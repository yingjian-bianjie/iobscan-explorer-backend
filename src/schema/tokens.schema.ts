import * as mongoose from 'mongoose';
import {ITokens} from "../types/schemaTypes/tokens.interface";
import { IAssetStruct } from '../types/schemaTypes/asset.interface';
export const TokensSchema = new mongoose.Schema({
    symbol: String,
    min_unit: String,
    scale: String,
    is_main_token: Boolean,
    initial_supply: String,
    max_supply: String,
    mintable: Boolean,
    owner: String,
    name: String,
    total_supply: String,
    update_block_height: Number
})
TokensSchema.index({symbol: 1}, {unique: true})
TokensSchema.index({min_unit: 1})

TokensSchema.statics = {
    async insertTokens(Tokens: ITokens) {
        //设置 options 查询不到就插入操作
        let {min_unit} = Tokens
        const options = {upsert: true, new: false, setDefaultsOnInsert: true}
        await this.findOneAndUpdate({ min_unit }, Tokens, options)
    },
    async queryAllTokens() {
        return await this.find({})
    },
    async queryMainToken() {
        return await this.findOne({is_main_token:true});
    },
    async findList(pageNum: number, pageSize: number): Promise<IAssetStruct[]> {
        return await this.find({'is_main_token':false})
            .select({
                symbol: 1,
                owner: 1,
                total_supply: 1,
                initial_supply: 1,
                max_supply: 1,
                mintable: 1,
            })
            .skip((pageNum - 1) * pageSize)
            .limit(pageSize).exec();
    },
    async findCount(): Promise<number> {
        return await this.find({'is_main_token':false}).countDocuments().exec();
    },
    async findOneBySymbol(symbol: string): Promise<IAssetStruct | null> {
        return await this.findOne({ symbol }).select({
            name: 1,
            owner: 1,
            total_supply: 1,
            initial_supply: 1,
            max_supply: 1,
            mintable: 1,
            scale: 1,
            min_unit:1
        });
    },
}
