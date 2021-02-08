import {Injectable} from '@nestjs/common';
import {InjectModel} from '@nestjs/mongoose';
import {Model} from 'mongoose'
import {StakingHttp} from "../http/lcd/staking.http";
import {GovHttp} from "../http/lcd/gov.http";
import {getTimestamp} from "../util/util";
import {moduleSlashing, moduleStaking, moduleStakingBondDenom,moduleGov,moduleGovDeposit} from "../constant";

@Injectable()
export class ParametersTaskService {
    constructor(
        @InjectModel('ParametersTask') private parametersTaskModel: Model<any>, 
        private readonly stakingHttp: StakingHttp,
        private readonly govHttp: GovHttp
    ) {
        this.doTask = this.doTask.bind(this);
    }

    async doTask(): Promise<any> {
        let parametersData = await this.stakingHttp.queryParametersFromSlashing(),
            needInsertData: any[] = [];
        let dbParametersData = await (this.parametersTaskModel as any).queryAllParameters()
        const stakingTokensData = await this.stakingHttp.getStakingTokens()
        const govTallyParams = await this.govHttp.getTallying()
        const govDepositParams = await this.govHttp.getDeposit()
        for (const parameterKey in parametersData) {
            const dbData = {
                module: moduleSlashing,
                key: parameterKey,
                cur_value: parametersData[parameterKey],
                create_time: '',
                update_time: ''
            }
            await needInsertData.push(dbData)
        }
        for (const stakingToken in stakingTokensData) {
            const stakingData = {
                module: moduleStaking,
                key: moduleStakingBondDenom,
                cur_value: stakingTokensData[moduleStakingBondDenom],
                create_time: '',
                update_time: ''
            }
            await needInsertData.push(stakingData)
        }
        for (const tallyParam in govTallyParams) {
            const tallyData = {
                module: moduleGov,
                key: tallyParam,
                cur_value: govTallyParams[tallyParam],
                create_time: '',
                update_time: ''
            }
            await needInsertData.push(tallyData)
        }
        const depositParams = govDepositParams && govDepositParams.min_deposit && govDepositParams.min_deposit[0] && govDepositParams.min_deposit[0].amount
        if (depositParams) {
            const depositData = {
                module: moduleGov,
                key: moduleGovDeposit,
                cur_value: depositParams,
                create_time: '',
                update_time: ''
            }
            await needInsertData.push(depositData)
        }
        if (dbParametersData.length === 0) {
            needInsertData.forEach(item => {
                item.create_time = getTimestamp()
                item.update_time = getTimestamp()
            })
            await (this.parametersTaskModel as any).insertParameters(needInsertData)
        } else {
            await this.handleData(dbParametersData, needInsertData)
        }
    }

    handleData(dbParameters, needInsertData) {
        let dbParametersMap = new Map()
        let needInsertDataMap = new Map()
        dbParameters.forEach(item => {
            dbParametersMap.set(item.key, item)
        })
        needInsertData.forEach(item => {
            needInsertDataMap.set(item.key, item)
        })
        if (needInsertDataMap.size > 0) {
            needInsertData.forEach(item => {
                if (dbParametersMap.has(item.key)) {
                    needInsertDataMap.get(item.key).update_time = getTimestamp()
                    needInsertDataMap.get(item.key).create_time = dbParametersMap.get(item.key).create_time
                } else {
                    needInsertDataMap.get(item.key).create_time = getTimestamp()
                }
            })
        }
        needInsertDataMap.forEach(item => {
            (this.parametersTaskModel as any).updateParameters(item)
        })
    }
}
