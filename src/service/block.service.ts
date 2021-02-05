import { Injectable } from '@nestjs/common';
import { Model } from 'mongoose';
import { InjectModel } from '@nestjs/mongoose';
import { ListStruct } from '../api/ApiResult';
import {
    BlockListReqDto,
    BlockDetailReqDto,
    ValidatorsetsReqDto } from '../dto/block.dto';
import {
    BlockListResDto,
    ValidatorsetsResDto,
    BlockStakingResDto} from '../dto/block.dto';
import { IBlock, IBlockStruct } from '../types/schemaTypes/block.interface';
import { BlockHttp } from '../http/lcd/block.http';
import { Logger } from '../logger';
import { addressPrefix } from '../constant';
import {getAddress, hexToBech32} from '../util/util';
@Injectable()
export class BlockService {

    constructor(
        @InjectModel('Block') private blockModel: Model<IBlock>,
        @InjectModel('StakingValidator') private stakingValidatorModel: any) {}

    async queryBlockList(query: BlockListReqDto): Promise<ListStruct<BlockListResDto[]>> {
        const { pageNum, pageSize, useCount } = query;
        let count: number;
        const b: IBlockStruct[] = await (this.blockModel as any).findList(pageNum, pageSize);
        if (useCount) {
            count = await (this.blockModel as any).findCount();
        }
        const res: BlockListResDto[] = b.map((b) => {
            return new BlockListResDto(b.height, b.hash, b.txn, b.time);
        });
        return new ListStruct(res, pageNum, pageSize, count);
    }

    async queryBlockDetail(p: BlockDetailReqDto): Promise<BlockListResDto | null> {
        let data: BlockListResDto | null = null;
        const { height } = p;
        const res: IBlockStruct | null = await (this.blockModel as any).findOneByHeight(height);
        if (res) {
            data = new BlockListResDto(res.height, res.hash, res.txn, res.time);
        }
        return data;
    }

    // blocks/staking/{height}
    async queryBlockStakingDetail(query: BlockDetailReqDto): Promise<BlockStakingResDto | null> {
        const { height } = query;
        let result: BlockStakingResDto | null = null;
        let data:any = {};

        let block_db = await (this.blockModel as any).findOneByHeight(height);
        block_db = JSON.parse(JSON.stringify(block_db));
        if (block_db) {
            let block_lcd =  await BlockHttp.queryBlockFromLcd(height);
            let latestBlock = await BlockHttp.queryLatestBlockFromLcd();
            let proposer = await this.stakingValidatorModel.findValidatorByPropopserAddr(block_db.proposer || '');
            let validatorsets = await BlockHttp.queryValidatorsets(height);
            data = {
                height: block_db.height,
                hash: block_db.hash,
                txn: block_db.txn,
                time: block_db.time,
                proposer: block_db.proposer
            };

            if (proposer && proposer.length) {
                data.proposer_moniker = (proposer[0].description || {}).moniker || '';
                data.proposer_addr = proposer[0].operator_address || '';
            }

            let signaturesMap:any = {};
            block_lcd.block.last_commit.signatures.forEach((item:any)=>{
                let address = hexToBech32(item.validator_address, addressPrefix.ica);
                signaturesMap[address] = item;
            }) 
            if (validatorsets) {
                data.total_validator_num = validatorsets ? validatorsets.length : 0;
                let icaAddr = hexToBech32(block_db.proposer, addressPrefix.ica);
                data.total_voting_power = 0;
                data.precommit_voting_power = 0;
                validatorsets.forEach((item)=>{
                    //TODO:hangtaishan 使用大数计算
                    data.total_voting_power += Number(item.voting_power || 0);
                    if (signaturesMap[item.address]) {
                        data.precommit_voting_power += Number(item.voting_power || 0);
                    }
                });
            }
            data.precommit_validator_num = 0;
            if (block_lcd) {
                try{
                    data.precommit_validator_num = block_lcd.block.last_commit.signatures.filter((item)=>{
                        return item.validator_address && item.validator_address.length;
                    }).length;
                }catch(e){
                    data.precommit_validator_num = 0;
                }
            }
            if (latestBlock) {
                data.latest_height = (latestBlock.block.header || {}).height;
            }
            result = new BlockStakingResDto(data);
        }
        return result;
    }

    // validatorset/{height}
    async queryValidatorset(query: ValidatorsetsReqDto): Promise<ListStruct<ValidatorsetsResDto[]>> {
        const { height, pageNum, pageSize, useCount } = query;
        let data_lcd = await BlockHttp.queryValidatorsets(height);
        let data = (data_lcd || []).slice((pageNum - 1) * pageSize, pageNum * pageSize);
        if (data && data.length) {
            let block = await (this.blockModel as any).findOneByHeight(Number(height));
            block = JSON.parse(JSON.stringify(block || '{}'));
            let validators = await this.stakingValidatorModel.queryAllValidators();
            if (validators.length) {
                let validatorMap = {};
                validators.forEach((v)=>{
                    validatorMap[v.proposer_addr] = v;
                });
                data.forEach((item)=>{
                    const proposer_addr = item.pub_key ? getAddress(item.pub_key).toLocaleUpperCase() : null
                    let validator = validatorMap[proposer_addr];
                    if (validator) {
                        (item as any).moniker = (validator.description || {}).moniker || '';
                        (item as any).operator_address = validator.operator_address || '';
                        (item as any).is_proposer = (validator.proposer_addr == block.proposer)
                    }
                })
            }
        }
        return new ListStruct(ValidatorsetsResDto.bundleData(data), pageNum, pageSize, data_lcd.length);
    }

    async queryLatestBlock(): Promise<IBlockStruct> {
        try {
            const blockStruct: IBlockStruct | null = await this.queryLatestBlockFromLcd();
            if(blockStruct){
                return blockStruct;
            }else {
                return await this.queryLatestBlockFromDB();
            }
        } catch (e) {
            Logger.warn('api-error:', e.message);
            return await this.queryLatestBlockFromDB();
        }

    }

    private async queryLatestBlockFromDB(): Promise<IBlockStruct> {
        return await (this.blockModel as any).findOneByHeightDesc();
    }

    private async queryLatestBlockFromLcd(): Promise<IBlockStruct | null> {
        const res = await BlockHttp.queryLatestBlockFromLcd();
        if(res && res.block_id && res.block && res.block.header && res.block.data){
            const blockStruct: IBlockStruct = {};
            blockStruct.height = res.block.header.height;
            blockStruct.time = res.block.header.time;
            blockStruct.txn = res.block.data.txs ? res.block.data.txs.length : 0;
            blockStruct.hash = res.block_id.hash;
            return blockStruct;
        }else {
            return null;
        }

    }
}
