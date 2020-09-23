import {Injectable} from '@nestjs/common';
import {InjectModel} from '@nestjs/mongoose';
import {Model} from 'mongoose';
import {StakingHttp} from "../http/lcd/staking.http";
import {addressTransform, pageNation} from "../util/util";
import {activeValidatorLabel, addressPrefix, moduleSlashing, ValidatorNumberStatus, ValidatorStatus} from "../constant";
import {
    AccountAddrReqDto, AccountAddrResDto,
    allValidatorReqDto,
    CommissionInfoReqDto,
    CommissionInfoResDto, stakingValidatorResDto,
    ValidatorDelegationsReqDto,
    ValidatorDelegationsResDto, ValidatorDetailResDtO,
    ValidatorUnBondingDelegationsReqDto, ValidatorUnBondingDelegationsResDto,
} from "../dto/staking.dto";
import {ListStruct} from "../api/ApiResult";
import {BlockHttp} from "../http/lcd/block.http";

@Injectable()
export default class StakingService {
    constructor(@InjectModel('Profiler') private profilerModel: Model<any>,
                @InjectModel('StakingSyncValidators') private stakingValidatorsModel: Model<any>,
                @InjectModel('Parameters') private parametersModel: Model<any>,
                @InjectModel('Tx') private txModel: Model<any>,
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
        let allValidatorVotingPower: number[] = []
        await allValidators.forEach(item => {
            if (item.status === ValidatorStatus[activeValidatorLabel] && item.jailed === false) {
                allValidatorVotingPower.push(Number(item.voting_power))
            }
        })
        let totalVotingPower: number = 0
        if (allValidatorVotingPower.length > 0) {
            totalVotingPower = await allValidatorVotingPower.reduce((total: number, item: number) => {
                return item + total
            })
        }
        return totalVotingPower
    }


    async getAllValCommission(q: CommissionInfoReqDto): Promise<ListStruct<CommissionInfoResDto>> {
        const allValCommissionInfo: any = await (this.stakingValidatorsModel as any).queryAllValCommission(q)
        allValCommissionInfo.data = CommissionInfoResDto.bundleData(allValCommissionInfo.data)
        return allValCommissionInfo
    }

    async getValidatorDelegationList(q: ValidatorDelegationsReqDto): Promise<ListStruct<ValidatorDelegationsResDto>> {
        const validatorAddr = q.address
        const allValidatorsMap = await this.getAllValidatorMonikerMap()
        const validatorDelegationsFromLcd = await this.stakingHttp.queryValidatorDelegationsFromLcd(validatorAddr)
        const allShares: number[] = []
        validatorDelegationsFromLcd.forEach(item => {
            //TODO:zhangjinbiao 使用位移进行大数字的计算
            allShares.push(Number(item.delegation.shares))
        })
        let totalShares = allShares.reduce((total: number, item: number) => {
            return item + total
        })
        let resultData = validatorDelegationsFromLcd.map(item => {
            return {
                moniker: allValidatorsMap.get(item.delegation.validator_address).description.moniker || '',
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
        if (resultData.length < q.pageSize) {
            result.data = ValidatorDelegationsResDto.bundleData(resultData)
        } else {
            let pageNationData = pageNation(resultData, q.pageSize)
            result.data = ValidatorDelegationsResDto.bundleData(pageNationData[q.pageNum - 1])
        }
        return new ListStruct(result.data, q.pageNum, q.pageSize, count)
    }

    async getValidatorUnBondingDelegations(q: ValidatorUnBondingDelegationsReqDto): Promise<ListStruct<ValidatorUnBondingDelegationsResDto>> {
        const validatorAddr = q.address
        const allValidatorsMoniker = await this.getAllValidatorMonikerMap()
        const valUnBondingDelegationsFromLcd = await this.stakingHttp.queryValidatorUnBondingDelegations(validatorAddr)
        let resultData = valUnBondingDelegationsFromLcd.map(item => {
            return {
                moniker: allValidatorsMoniker.get(item.validator_address).description.moniker || '',
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
        if (resultData.length < q.pageSize) {
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

    async getValidatorDetail(q: ValidatorDelegationsReqDto): Promise<ValidatorDetailResDtO> {

        const validatorAddress = q.address
        let result: any = null;
        let validatorDetail = await (this.stakingValidatorsModel as any).queryDetailByValidator(validatorAddress);
        if (validatorDetail) {
            const moduleName = moduleSlashing
            const signedBlocksWindow = await (this.parametersModel as any).querySignedBlocksWindow(moduleName)
            const latestBlock = await BlockHttp.queryLatestBlockFromLcd()
            const signedBlocksWindowCurVal = signedBlocksWindow.cur_value || undefined
            const latestBlockHeight = latestBlock.block.header.height || 0
            const startHeight = validatorDetail.start_height || 0
            if (Number(signedBlocksWindowCurVal) < (Number(latestBlockHeight) - Number(startHeight))) {
                validatorDetail.stats_blocks_window = signedBlocksWindowCurVal
            } else {
                validatorDetail.stats_blocks_window = (Number(latestBlockHeight) - Number(startHeight))
            }
            validatorDetail.valStatus = ValidatorNumberStatus[validatorDetail.status]
            validatorDetail.total_power = await this.getTotalVotingPower()
            validatorDetail.tokens = Number(validatorDetail.tokens)
            validatorDetail.bonded_stake = (Number(validatorDetail.self_bond.amount) * (Number(validatorDetail.tokens) / Number(validatorDetail.delegator_shares))).toString()
            validatorDetail.owner_addr = addressTransform(validatorDetail.operator_address, addressPrefix.iaa)
            result = new ValidatorDetailResDtO(validatorDetail)
        }
        return result
    }

    async getAddressAccount(q: AccountAddrReqDto): Promise<AccountAddrResDto> {
        const address = addressTransform(q.address, addressPrefix.iaa)
        const operatorAddress = addressTransform(q.address, addressPrefix.iva)
        const balancesArray = await this.stakingHttp.queryBalanceByAddress(address)
        //TODO:zhangjinbiao 奖励地址需要加上
        const allValidatorsMap = await this.getAllValidatorMonikerMap()
        const allProfilerAddress = await (this.profilerModel as any).queryProfileAddress()
        const validator = allValidatorsMap.get(operatorAddress)
        const deposits = await (this.txModel as any).queryDepositsByAddress(address)


        let profilerAddressMap = new Map()
        if (allProfilerAddress && allProfilerAddress.length > 0) {
            allProfilerAddress.forEach(item => {
                profilerAddressMap.set(item.address, item)
            })
        }
        let result: any = {}
        result.amount = balancesArray
        //TODO:zhangjinbiao 奖励地址需要加上
        result.withdrawAddress = ''
        result.address = address
        result.moniker = validator.description.moniker
        result.operator_address = allValidatorsMap.has(address) ? validator.operator_address : '--'
        result.isProfiler = profilerAddressMap.size > 0 ? profilerAddressMap.has(address) : false
        if (allValidatorsMap.has(operatorAddress) && !validator.jailed) {
            result.status = ValidatorNumberStatus[validator.status]
        } else {
            result.status = ''
        }
        if (deposits && deposits.data && deposits.data.length > 0) {
            //TODO:zhangjinbiao 处理查询出来与Gov相关的交易列表计算总的amount

        } else {
            result.deposits = {}
        }
        return new AccountAddrResDto(result)
    }
}
