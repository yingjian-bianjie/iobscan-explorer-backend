import { Injectable } from '@nestjs/common';
import { InjectModel } from '@nestjs/mongoose';
import { Model } from "mongoose";
import { GovHttp } from "../http/lcd/gov.http";
import { proposalStatus, govParams, voteOptions } from '../constant'
import { StakingHttp } from "../http/lcd/staking.http";
import { addressTransform } from "../util/util";
import { addressPrefix,proposal } from "../constant";
import { getTimestamp, formatDateStringToNumber, splitString } from '../util/util';
import { cfg } from "../config/config"
@Injectable()
export class ProposalTaskService {
    constructor(
        @InjectModel('ParametersTask') private parametersTaskModel: Model<any>,
        @InjectModel('Tx') private txModel: any,
        @InjectModel('Proposal') private proposalModel: any,
        @InjectModel('ProposalDetail') private proposalDetailModel: any,
        @InjectModel('StakingSyncValidators') private stakingValidatorsModel: any,
        private readonly govHttp: GovHttp,
        private readonly stakingHttp: StakingHttp
    ) {
        this.doTask = this.doTask.bind(this);
    }
    async doTask(): Promise<void> {
        const proposalFromLcd = await this.govHttp.getProposals(cfg.taskCfg.proposalsLimit);
        const proposalFromDb = await this.proposalModel.queryAllProposals();
        const proposalDetailFromDb = await this.proposalDetailModel.queryAllProposalsDetail();
        const govParamsFromDb = await (this.parametersTaskModel as any).queryGovParams();
        let govParamsFromDbMap = new Map();
        let proposalFromLcdMap = new Map();
        let proposalFromDbMap = new Map();
        if (govParamsFromDb && govParamsFromDb.length > 0) {
            govParamsFromDb.forEach(item => {
                govParamsFromDbMap.set(item.key, item)
            });
        }
        if (proposalFromLcd && proposalFromLcd.length > 0) {
            proposalFromLcd.forEach(item => {
                proposalFromLcdMap.set(item.id, item)
            })
        }
        let deleteIdData = []
        if (proposalFromDb && proposalFromDb.length > 0) {
            proposalFromDb.forEach(item => {
                proposalFromDbMap.set(item.id, item)
                if (!proposalFromLcdMap.has(item.id)) {
                    deleteIdData.push(item.id)
                }
            })
        }
        if (deleteIdData.length > 0) {
            this.proposalModel.updateProposals(deleteIdData)
        }
        // 以上删除提案已完成
        let insertData = [];
        let insertIds = [];
        if (proposalFromLcdMap.size > 0) {
            proposalFromLcdMap.forEach(item => {
                if (!proposalFromDbMap.has(item.id)) {
                    insertData.push(item);
                    insertIds.push(item.id);
                }
            })
        }
        for (let item of insertData) {
            let txData = await this.txModel.querySubmitProposalById(String(item.id));
            item.hash = txData && txData.tx_hash;
            item.proposer = txData && txData.msgs && txData.msgs[0] && txData.msgs[0].msg && txData.msgs[0].msg.proposer || '';
            item.initial_deposit = txData && txData.msgs && txData.msgs[0] && txData.msgs[0].msg && txData.msgs[0].msg.initial_deposit || {};
            item.status = proposalStatus[item.status];
            if (item.content && item.content['@type']) {
                item.content.type = splitString(item.content['@type'], '.').replace(proposal,'');
                delete item.content['@type']
            }
            item.submit_time = formatDateStringToNumber(item.submit_time);
            item.deposit_end_time = formatDateStringToNumber(item.deposit_end_time);
            item.voting_start_time = formatDateStringToNumber(item.voting_start_time);
            item.voting_end_time = formatDateStringToNumber(item.voting_end_time);
            item.min_deposit = govParamsFromDbMap.get(govParams.min_deposit).cur_value;
            item.quorum = govParamsFromDbMap.get(govParams.quorum).cur_value;
            item.threshold = govParamsFromDbMap.get(govParams.threshold).cur_value;
            item.veto_threshold = govParamsFromDbMap.get(govParams.veto_threshold).cur_value;
        }
        // 以上添加 proposer and initial_deposit 修改字段完成
        let updateDbData = proposalFromDb.filter(item => {
            return item.status == proposalStatus['PROPOSAL_STATUS_DEPOSIT_PERIOD'] || item.status == proposalStatus['PROPOSAL_STATUS_VOTING_PERIOD']
        })
        let updateData = updateDbData.concat(insertData);
        let updateProposalDetails = [];
        for (const proposal of updateData) {
            if (proposal.status == proposalStatus['PROPOSAL_STATUS_VOTING_PERIOD'] || proposal.status == proposalStatus['PROPOSAL_STATUS_DEPOSIT_PERIOD']) {
                let { current_tally_result, tally_details } = await this.tallyDetails(proposal.id);
                proposal.current_tally_result = current_tally_result;
                updateProposalDetails.push({
                    id: proposal.id,
                    tally_details
                })
            } 
            if (!insertIds.includes(proposal.id)) {
                let txData = await this.txModel.querySubmitProposalById(String(proposal.id));
                proposal.hash = txData && txData.tx_hash;
                proposal.proposer = txData && txData.msgs && txData.msgs[0] && txData.msgs[0].msg && txData.msgs[0].msg.proposer || '';
                proposal.initial_deposit = txData && txData.msgs && txData.msgs[0] && txData.msgs[0].msg && txData.msgs[0].msg.initial_deposit || {};
                let proposalFromLcd = proposalFromLcdMap.get(proposal.id);
                proposal.status = proposalStatus[proposalFromLcd.status];
                proposal.final_tally_result = proposalFromLcd.final_tally_result
                proposal.total_deposit = proposalFromLcd.total_deposit
                proposal.voting_start_time = formatDateStringToNumber(proposalFromLcd.voting_start_time)
                proposal.voting_end_time = formatDateStringToNumber(proposalFromLcd.voting_end_time)
            }
        }
        await this.insertAndUpdateProposalDetail(updateProposalDetails, proposalDetailFromDb)
        await this.insertAndUpdateProposal(updateData,proposalFromDb)
    }
    async tallyDetails(proposal_id) {
        const validators = await this.stakingValidatorsModel.queryActiveVal();
        let systemVotingPower = 0;
        let totalVotingPower = 0;
        let delegatorsGovInfo = {};
        let validatorGovInfo = {};
        if (validators && validators.length > 0) {
            validators.forEach(validator => {
                validatorGovInfo[validator.operator_address] = {
                    address: validator.operator_address,
                    moniker: validator.description && validator.description.moniker || '',
                    bondedtokens: Number(validator.tokens),
                    delShares: Number(validator.delegator_shares),
                    vote: '',
                    delDeductionShares: 0,
                    selfDelDeductionShares: 0
                }
                // todo: duanjie 使用大数计算
                systemVotingPower += Number(validator.tokens)
            });
        }
        const votes = await this.txModel.queryVoteByProposalId(Number(proposal_id))
        if (votes && votes.length > 0) {
            for (let vote of votes) {
                vote = vote && vote.msg && vote.msg[0]
                delegatorsGovInfo[vote.voter] = {
                    address: vote.voter,
                    vote: voteOptions[vote.option],
                    moniker: '',
                    selfDelVotingPower: 0,
                    delVotingPower: 0,
                    notVoteVotingPower: 0,
                    isValidator: false
                }
                let delegatorsDelegationsFromLcd = await this.stakingHttp.queryDelegatorsDelegationsFromLcd(vote.voter)
                let delegators = []
                if (delegatorsDelegationsFromLcd && delegatorsDelegationsFromLcd.result && delegatorsDelegationsFromLcd.result.length > 0) {
                    delegatorsDelegationsFromLcd.result.forEach(item => {
                        if (validatorGovInfo[item.delegation.validator_address]) {
                            delegators.push(item.delegation)
                        }
                    })
                }
                if (delegators.length > 0) {
                    delegators.forEach(delegator => {
                        let voteIva = addressTransform(vote.voter, addressPrefix.iva)
                        if (delegator.validator_address == voteIva) {
                            delegatorsGovInfo[vote.voter].isValidator = true;
                            let votingPower = Number(delegator.shares) * (validatorGovInfo[voteIva].bondedtokens / validatorGovInfo[voteIva].delShares);
                            totalVotingPower += votingPower;
                            delegatorsGovInfo[vote.voter].selfDelVotingPower += votingPower;
                            delegatorsGovInfo[vote.voter].moniker = validatorGovInfo[voteIva].moniker;
                            validatorGovInfo[voteIva].selfDelDeductionShares += Number(delegator.shares);
                            validatorGovInfo[voteIva].vote = delegatorsGovInfo[vote.voter].vote;
                        } else {
                            let votingPower = Number(delegator.shares) * (validatorGovInfo[delegator.validator_address].bondedtokens / validatorGovInfo[delegator.validator_address].delShares);
                            totalVotingPower += votingPower;
                            delegatorsGovInfo[vote.voter].delVotingPower += votingPower;
                            validatorGovInfo[delegator.validator_address].delDeductionShares += Number(delegator.shares);
                        }
                    })
                }
            }
        }
        let delAndValGovInfo = {};
        for (const address in delegatorsGovInfo) {
            if (delegatorsGovInfo[address].isValidator) {
                delAndValGovInfo[address] = delegatorsGovInfo[address];
            }
        }
        for (const address in delAndValGovInfo) {
            let ivaAddress = addressTransform(address, addressPrefix.iva)
            let validator = validatorGovInfo[ivaAddress];
            let votingPower = (validator.delShares - validator.delDeductionShares - validator.selfDelDeductionShares) * (validator.bondedtokens / validator.delShares);
            delegatorsGovInfo[address].notVoteVotingPower += votingPower;
            totalVotingPower += votingPower
        }
        let yes = 0, abstain = 0, no = 0, no_with_veto = 0;
        for (const address in delegatorsGovInfo) {
            let info = delegatorsGovInfo[address]
            switch (info.vote) {
                case 'yes':
                    yes += info.selfDelVotingPower + info.delVotingPower + info.notVoteVotingPower
                    break;
                case 'abstain':
                    abstain += info.selfDelVotingPower + info.delVotingPower + info.notVoteVotingPower
                    break;
                case 'no':
                    no += info.selfDelVotingPower + info.delVotingPower + info.notVoteVotingPower
                    break;
                case 'no_with_veto':
                    no_with_veto += info.selfDelVotingPower + info.delVotingPower + info.notVoteVotingPower
                    break;
                default:
                    break;
            }
        }
        let tally_details = Object.values(delegatorsGovInfo)
        let result = {
            current_tally_result: {
                system_voting_power:systemVotingPower,
                total_voting_power: totalVotingPower,
                yes,
                abstain,
                no,
                no_with_veto
            },
            tally_details
        }
        return result
    }

    private async insertAndUpdateProposal(updateData,proposalFromDb) {
        if (updateData && updateData.length > 0) {
            for (let proposal of updateData) {
                if (proposalFromDb && proposalFromDb.length > 0) {
                    let created = proposalFromDb.filter(item => {
                        return item.id == proposal.id
                    })
                    if (created && created.length > 0) {
                        proposal.create_time = created[0].create_time
                    }
                }else {
                    proposal.create_time = getTimestamp()
                }
                proposal.update_time = getTimestamp()
                await this.proposalModel.insertProposal(proposal)
            }
        }
    }

    private async insertAndUpdateProposalDetail(updateProposalDetails,proposalDetailFromDb) {
        if (updateProposalDetails && updateProposalDetails.length > 0) {
            for (const proposal of updateProposalDetails) {
                if (proposalDetailFromDb && proposalDetailFromDb.length > 0) {
                    let created = proposalDetailFromDb.filter(item => {
                        return item.id == proposal.id
                    })
                    if (created && created.length > 0) {
                        proposal.create_time = created[0].create_time
                    }
                } else {
                    proposal.create_time = getTimestamp()
                }
                proposal.update_time = getTimestamp()
                await this.proposalDetailModel.insertProposalDetail(proposal)
            }
        }
    }
}
