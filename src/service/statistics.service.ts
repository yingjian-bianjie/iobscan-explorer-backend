import { IconUriLcdDto } from './../dto/http.dto';
import { Injectable, Query } from '@nestjs/common';
import { Model } from 'mongoose';
import { InjectModel } from '@nestjs/mongoose';
import { StatisticsResDto,NetworkStatisticsResDto } from '../dto/statistics.dto';
import { IBlock, IBlockStruct } from '../types/schemaTypes/block.interface';
import { INft } from '../types/schemaTypes/nft.interface';
import { BlockHttp } from '../http/lcd/block.http';
import { StakingHttp } from '../http/lcd/staking.http';
import { BankHttp } from '../http/lcd/bank.http';
import { StatisticsStruct } from '../types/statistics.interface'
import { correlationStr } from '../constant'
@Injectable()
export class StatisticsService {

    constructor(
        @InjectModel('Block') private blockModel: Model<IBlock>,
        @InjectModel('Nft') private nftModel: Model<INft>,
        @InjectModel('Tx') private txModel: any,
        @InjectModel('Validators') private validatorModel: any,
        @InjectModel('Identity') private identityModel: any,
        @InjectModel('Denom') private denomModel: any,
        @InjectModel('StakingSyncValidators') private stakingValidatorsModel: any,
        @InjectModel('Tokens') private tokensModel: Model<any>,
    ) {
    }

    async queryStatistics(query: string): Promise<StatisticsResDto> {
        const params: Array<string> = query['params'].split(',');
        const that = this
        let indicators = []
        params.forEach( code => {
            switch (code) {
                case '202':
                    indicators.push(this.queryConsensusValidatorCount.bind(that))
                break;
                case '203':
                    indicators.push(this.queryAvgBlockTime.bind(that))
                break;
                case '204':
                    indicators.push(this.queryAssetCount.bind(that))
                break;
                case '205':
                    indicators.push(this.queryDenomCount.bind(that))
                break;
                case '206':
                    indicators.push(this.queryServiceCount.bind(that))
                break;
                case '207':
                    indicators.push(this.queryIdentityCount.bind(that))
                break;
                case '208':
                    indicators.push(this.queryValidatorNumCount.bind(that))
                    break;
                default:
                    break;
            }
        })
        const dateArray = await Promise.all(indicators.map(p=>p()))
        let statisticsDate: StatisticsStruct = {}
        params.forEach((func, index) => {
            statisticsDate[correlationStr[func]] = dateArray[index]
        })
        return new StatisticsResDto(statisticsDate)
    }

    async queryNetworkStatistics(query: string): Promise<NetworkStatisticsResDto> {
        const params: Array<string> = query['params'].split(',');
        const latestBlock = await BlockHttp.queryLatestBlockFromLcd();
        const that = this
        let indicators = []
        if ((!params.includes('200')) && params.includes('201')) {
            params.unshift('200')
        }
        params.forEach( code => {
            switch (code) {
                case '200':
                    indicators.push(this.queryLatestHeightAndTimeAndValidator.bind(that,latestBlock))
                    break;
                case '201':
                    indicators.push(this.queryTxCount.bind(that))
                break;
                case '209':
                    indicators.push(this.queryBondedTokensInformation.bind(that))
                    break;
                default:
                    break;
            }
        })
        const dateArray = await Promise.all(indicators.map(p=>p()))
        let statisticsDate: StatisticsStruct = {}
        params.forEach((func, index) => {
            statisticsDate[correlationStr[func]] = dateArray[index]
        })
        return new NetworkStatisticsResDto(statisticsDate);
    }

    async queryLatestHeightAndTimeAndValidator(latestBlock?:any): Promise<{height:number,latestBlockTime:number,moniker:string,validator_icon:string,operator_addr:string} | null> {
        let result: any = { height: 0, latestBlockTime: 0, moniker: '', validator_icon: '', operator_addr: '' };
        let proposer_address: string;
        if (latestBlock && latestBlock.block && latestBlock.block.header) {
            result.height = latestBlock.block.header.height;
            result.latestBlockTime = Math.floor(new Date(latestBlock.block.header.time || '').getTime() / 1000);
            proposer_address = latestBlock.block.header.proposer_address;
        }else {
            const res: IBlockStruct | null = await (this.blockModel as any).findOneByHeightDesc();
            if (res) {
                result.height = res.height;
                result.latestBlockTime = Number(res.time);
                proposer_address = res.proposer;
            }
        }
        let validator = await this.stakingValidatorsModel.findValidatorByPropopserAddr(proposer_address);
        result.moniker = validator && validator[0] && validator[0].description && validator[0].description.moniker;
        result.validator_icon = validator && validator[0] && validator[0].icon;
        result.operator_addr = validator && validator[0] && validator[0].operator_address;
        return result;
    }

    async queryAvgBlockTime(): Promise<number | null> {
        const latestBlock: IBlockStruct | null = await (this.blockModel as any).findOneByHeightDesc();
        const num100Block: IBlockStruct | null = await (this.blockModel as any).findNum100Height();
        if (latestBlock && num100Block) {
            const latestTime = Number(new Date(latestBlock.time).getTime());
            const num100Time = Number(new Date(num100Block.time).getTime());
            let avgTime: number;
            if (latestBlock.height - num100Block.height >= 100) {
                avgTime = (latestTime - num100Time) / 100;
            } else {
                //可能当前区块高度还不到100
                const diff: number = latestBlock.height - num100Block.height;
                if(diff){
                    avgTime = (latestTime - num100Time) / diff;
                }else{
                    return null;
                }
            }
            return !isNaN(avgTime) ? Number(avgTime.toFixed(2)) : null;
        } else {
            return null;
        }
    }

    async queryAssetCount(): Promise<number | null>{
        return await (this.nftModel as any).findCount();
    }

    async queryConsensusValidatorCount(): Promise<number | null>{
        return await (this.validatorModel as any).findCount(false);
    }

    async queryTxCount():Promise<any>{
        return await (this.txModel as any).queryTxCountStatistics();
    }

    async queryServiceCount():Promise<any>{
        return await (this.txModel as any).queryServiceCountStatistics();
    }

    async queryIdentityCount(query:any):Promise<any>{
        return await this.identityModel.queryIdentityCount(query);
    }

    async queryDenomCount():Promise<any>{
        return await this.denomModel.queryAllCount();
    }

    async queryValidatorNumCount():Promise<any>{
        return await this.stakingValidatorsModel.queryActiveValCount();
    }

    async queryBondedTokensInformation(): Promise<any>{
        const [bondedTokensLcd,totalSupplyLcd] = await Promise.all([StakingHttp.getBondedTokens(),BankHttp.getTotalSupply()])
        const bonded_tokens = bondedTokensLcd && bondedTokensLcd.bonded_tokens || '0'
        const mainToken = await (this.tokensModel as any).queryMainToken()
        let total_supply: string = '0';
        if (mainToken) {
            if (totalSupplyLcd && totalSupplyLcd.supply && totalSupplyLcd.supply.length > 0) {
                totalSupplyLcd.supply.map(item => {
                    if (item.denom === mainToken.min_unit) {
                        total_supply = item.amount
                    }
                })
            }
        }
        return { bonded_tokens, total_supply };
    }
}
