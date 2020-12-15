import {Injectable} from '@nestjs/common';
import {InjectModel} from '@nestjs/mongoose';
import {Model} from 'mongoose';
import {StakingHttp} from "../http/lcd/staking.http";
import {addressTransform, pageNation} from "../util/util";
import {
    activeValidatorLabel,
    addressPrefix,
    jailedValidatorLabel,
    moduleSlashing,
    ValidatorNumberStatus,
    ValidatorStatus
} from "../constant";
import {
    AccountAddrReqDto,
    AccountAddrResDto,
    allValidatorReqDto,
    CommissionInfoReqDto,
    CommissionInfoResDto, stakingValidatorResDto,
    ValidatorDelegationsReqDto,ValidatorDelegationsQueryReqDto,
    ValidatorDelegationsResDto, ValidatorDetailResDto,
    ValidatorUnBondingDelegationsReqDto, ValidatorUnBondingDelegationsResDto,ValidatorUnBondingDelegationsQueryReqDto,
    DelegatorsDelegationsReqDto,DelegatorsDelegationsResDto,
    DelegatorsUndelegationsReqDto, DelegatorsUndelegationsResDto,
    DelegatorsDelegationsParamReqDto, DelegatorsUndelegationsParamReqDto,
} from "../dto/staking.dto";
import {ListStruct} from "../api/ApiResult";
import {BlockHttp} from "../http/lcd/block.http";
import { cfg } from '../config/config';
import {DistributionHttp} from "../http/lcd/distribution.http";

@Injectable()
export default class StakingService {
    constructor(@InjectModel('Profiler') private profilerModel: Model<any>,
                @InjectModel('StakingSyncValidators') private stakingValidatorsModel: Model<any>,
                @InjectModel('Parameters') private parametersModel: Model<any>,
                @InjectModel('Tx') private txModel: Model<any>,
                @InjectModel('Tokens') private tokensModel: any,
                private readonly stakingHttp: StakingHttp,
    ) {
    }

    async getAllValidatorMonikerMap() {
        const allValidators = await (this.stakingValidatorsModel as any).queryAllValidators()
        let allValidatorsMonikerMap = new Map()
        allValidators.forEach(item => {
            allValidatorsMonikerMap.set(item.operator_address, item)
        })
        return allValidatorsMonikerMap
    }

    async getTotalVotingPower() {
        const allValidators = await (this.stakingValidatorsModel as any).queryAllValidators()
        let totalVotingPower: number = 0;
        allValidators.forEach(item => {
            if (item.status === ValidatorStatus['bonded'] && item.jailed === false) {
                totalVotingPower += Number(item.voting_power);
            }
        })
        return totalVotingPower
    }


    async getAllValCommission(q: CommissionInfoReqDto): Promise<ListStruct<CommissionInfoResDto>> {
        const allValCommissionInfo: any = await (this.stakingValidatorsModel as any).queryAllValCommission(q)
        allValCommissionInfo.data = CommissionInfoResDto.bundleData(allValCommissionInfo.data)
        return allValCommissionInfo
    }

    async getValidatorDelegationList(p: ValidatorDelegationsReqDto,q: ValidatorDelegationsQueryReqDto): Promise<ListStruct<ValidatorDelegationsResDto>> {
        const validatorAddr = p.address
        const allValidatorsMap = await this.getAllValidatorMonikerMap()
        const validatorDelegationsFromLcd = await this.stakingHttp.queryValidatorDelegationsFromLcd(validatorAddr)
        const allShares: number[] = [];
        (validatorDelegationsFromLcd || []).forEach(item => {
            //TODO:zhangjinbiao 使用位移进行大数字的计算
            if(item){
                allShares.push(Number(item.delegation.shares))
            }
        })
        let totalShares;
        if (allShares && allShares.length) {
            totalShares = allShares.reduce((total: number, item: number) => {
                return item + total
            })
        }
        let resultData = (validatorDelegationsFromLcd || []).map(item => {
            let validator = allValidatorsMap.get(item.delegation.validator_address);
            return {
                moniker: validator ? validator.description.moniker : '',
                address: item.delegation.delegator_address || '',
                amount: item.balance || '',
                self_shares: item.delegation.shares || '',
                total_shares: totalShares || '',
            }
        })
        const count = resultData.length
        let result: any = {}
        if (q.useCount) {
            result.count = count
        }
        if (resultData.length <= q.pageSize) {
            result.data = ValidatorDelegationsResDto.bundleData(resultData)
        } else {
            let pageNationData = pageNation(resultData, q.pageSize)
            result.data = ValidatorDelegationsResDto.bundleData(pageNationData[q.pageNum - 1])
        }
        return new ListStruct(result.data, q.pageNum, q.pageSize, count)
    }

    async getValidatorUnBondingDelegations(p: ValidatorUnBondingDelegationsReqDto,q: ValidatorUnBondingDelegationsQueryReqDto): Promise<ListStruct<ValidatorUnBondingDelegationsResDto>> {
        const validatorAddr = p.address
        const allValidatorsMoniker = await this.getAllValidatorMonikerMap()
        const valUnBondingDelegationsFromLcd = await this.stakingHttp.queryValidatorUnBondingDelegations(validatorAddr)
        let resultData = (valUnBondingDelegationsFromLcd || []).map(item => {
            let validator = allValidatorsMoniker.get(item.validator_address);
            return {
                moniker: validator ? validator.description.moniker : '',
                address: item.delegator_address || '',
                amount: item.entries[0].balance || '',
                block: item.entries[0].creation_height || '',
                until: item.entries[0].completion_time || '',
            }

        })
        const count = resultData.length
        let result: any = {}
        if (q.useCount) {
            result.count = count
        }
        if (resultData.length <= q.pageSize) {
            result.data = ValidatorUnBondingDelegationsResDto.bundleData(resultData)
        } else {
            let pageNationData = pageNation(resultData, q.pageSize)
            result.data = ValidatorUnBondingDelegationsResDto.bundleData(pageNationData[q.pageNum - 1])
        }
        return new ListStruct(result.data, q.pageNum, q.pageSize, count)

    }

    async getValidatorsByStatus(q: allValidatorReqDto): Promise<ListStruct<stakingValidatorResDto>> {
        const validatorList = await (this.stakingValidatorsModel as any).queryValidatorsByStatus(q)
        const totalVotingPower = await this.getTotalVotingPower()
        validatorList.data.forEach(item => {
            item.voting_rate = item.voting_power / totalVotingPower
        })
        let result: any = {}
        result.data = stakingValidatorResDto.bundleData(validatorList.data)
        result.count = validatorList.count
        return new ListStruct(result.data, q.pageNum, q.pageSize, result.count)
    }

    async getValidatorDetail(q: ValidatorDelegationsReqDto): Promise<ValidatorDetailResDto> {

        const validatorAddress = q.address
        let result: any = null;
        let validatorDetail = await (this.stakingValidatorsModel as any).queryDetailByValidator(validatorAddress);
        if (validatorDetail) {
            const moduleName = moduleSlashing;
            const signedBlocksWindow = await (this.parametersModel as any).querySignedBlocksWindow(moduleName);
            const latestBlock = await BlockHttp.queryLatestBlockFromLcd();
            const signedBlocksWindowCurVal = signedBlocksWindow.cur_value || undefined;
            const latestBlockHeight = latestBlock ? latestBlock.block.header.height : 0;
            const startHeight = validatorDetail.start_height || 0;
            if (Number(signedBlocksWindowCurVal) < (Number(latestBlockHeight) - Number(startHeight))) {
                validatorDetail.stats_blocks_window = signedBlocksWindowCurVal;
            } else {
                validatorDetail.stats_blocks_window = (Number(latestBlockHeight) - Number(startHeight));
            }
            if(!validatorDetail.jailed){
                validatorDetail.valStatus = ValidatorNumberStatus[validatorDetail.status]
            }else {
                validatorDetail.valStatus = jailedValidatorLabel
            }
            validatorDetail.total_power = await this.getTotalVotingPower()
            validatorDetail.tokens = Number(validatorDetail.tokens)
            validatorDetail.bonded_stake = (Number(validatorDetail.self_bond.amount) * (Number(validatorDetail.tokens) / Number(validatorDetail.delegator_shares))).toString()
            validatorDetail.owner_addr = addressTransform(validatorDetail.operator_address, addressPrefix.iaa)
            result = new ValidatorDetailResDto(validatorDetail)
        }
        return result;
    }

    async getAddressAccount(q: AccountAddrReqDto): Promise<AccountAddrResDto> {
        const address = addressTransform(q.address, addressPrefix.iaa)
        const operatorAddress = addressTransform(q.address, addressPrefix.iva)
        const balancesArray = await this.stakingHttp.queryBalanceByAddress(address)
        const withdrawAddress =  await DistributionHttp.queryWithdrawAddressByDelegator(address)
        const allValidatorsMap = await this.getAllValidatorMonikerMap()
        const allProfilerAddress = await (this.profilerModel as any).queryProfileAddress()
        const validator = allValidatorsMap.get(operatorAddress)
        // const deposits = await (this.txModel as any).queryDepositsByAddress(address)

        let profilerAddressMap = new Map()
        if (allProfilerAddress && allProfilerAddress.length > 0) {
            allProfilerAddress.forEach(item => {
                profilerAddressMap.set(item.address, item)
            })
        }
        let result: any = {}
        result.amount = balancesArray || []
        result.withdrawAddress = withdrawAddress.address
        result.address = address
        result.moniker =  validator && validator.description ? validator.description.moniker :''
        result.operator_address = allValidatorsMap.has(operatorAddress) ? validator.operator_address : ''
        result.isProfiler = profilerAddressMap.size > 0 ? profilerAddressMap.has(address) : false
        if (allValidatorsMap.has(operatorAddress) && !validator.jailed) {
            result.status = ValidatorNumberStatus[validator.status]
        } else {
            result.status = jailedValidatorLabel
        }
        // if (deposits && deposits.data && deposits.data.length > 0) {
            //TODO:zhangjinbiao 处理查询出来与Gov相关的交易列表计算总的amount

        // } else {
            result.deposits = {}
        // }
        return new AccountAddrResDto(result)
    }

    async getDelegatorsDelegations(p:DelegatorsDelegationsParamReqDto,q: DelegatorsDelegationsReqDto): Promise<ListStruct<DelegatorsDelegationsResDto>> {
        const { pageNum, pageSize } = q
        const { delegatorAddr } = p
        const delegatorsDelegationsFromLcd = await this.stakingHttp.queryDelegatorsDelegationsFromLcd(delegatorAddr)
        const dataLcd = delegatorsDelegationsFromLcd ? delegatorsDelegationsFromLcd.result : []
        const count = dataLcd ? dataLcd.length : 0;
        const data = dataLcd ? dataLcd.slice((pageNum - 1) * pageSize, pageNum * pageSize) : []
        const allValidatorsMap = await this.getAllValidatorMonikerMap()
        const resultData = data.map(item => {
            let validator = allValidatorsMap.get(item.delegation.validator_address);
            return {
                address: item.delegation.validator_address || '',
                moniker: validator ? validator.description.moniker : '',
                amount: item.balance || '',
                shares: item.delegation.shares,
                //height: delegatorsDelegationsFromLcd.height || '',
            }
        })
        const result: any = {}
        result.count = count
        result.data = DelegatorsDelegationsResDto.bundleData(resultData)
        return new ListStruct(result.data, pageNum, pageSize, result.count)
    }

    async getDelegatorsUndelegations(p:DelegatorsUndelegationsParamReqDto,q: DelegatorsUndelegationsReqDto): Promise<ListStruct<DelegatorsUndelegationsResDto>> {
        const { pageNum, pageSize } = q
        const { delegatorAddr } = p
        const delegatorsDelegationsFromLcd = await this.stakingHttp.queryDelegatorsUndelegationsFromLcd(delegatorAddr)
        const dataLcd = delegatorsDelegationsFromLcd ? delegatorsDelegationsFromLcd.result : []
        const aminToken = await this.tokensModel.queryMainToken();
        const count =dataLcd ? dataLcd.length : 0
        const data = dataLcd ? dataLcd.slice((pageNum - 1) * pageSize, pageNum * pageSize) : []
        const allValidatorsMap = await this.getAllValidatorMonikerMap()
        const resultData = data.map(item => {
            const denom:string = (aminToken || {}).min_unit || '';
            let entries:any = item && item.entries || []
            const amount =  entries && entries.length > 0 ? entries[0].balance : ''
            let validator = allValidatorsMap.get(item.validator_address);
            return {
                address: item.validator_address || '',
                moniker: validator ? validator.description.moniker : '',
                amount: {
                    denom: denom || '',
                    amount: amount || ''
                },
                height:  entries && entries.length > 0 ? entries[0].creation_height : '',
                end_time: entries && entries.length > 0 ? entries[0].completion_time : ''
            }
        })
        const result: any = {}
        result.count = count
        result.data = DelegatorsUndelegationsResDto.bundleData(resultData)
        return new ListStruct(result.data, pageNum, pageSize, result.count)
    }
}
