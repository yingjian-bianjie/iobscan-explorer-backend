import { Injectable } from '@nestjs/common';
import { Model } from 'mongoose';
import { InjectModel } from '@nestjs/mongoose';
import { ListStruct } from '../api/ApiResult';
import {
    AssetListReqDto,
    AssetDetailReqDto
} from '../dto/asset.dto';
import {
    AssetListResDto,
    AssetDetailResDto
} from '../dto/asset.dto';
import { IAssetStruct } from '../types/schemaTypes/asset.interface';
@Injectable()
export class AssetService {

    constructor(
        @InjectModel('Tokens') private tokensModel: Model<any>,
        @InjectModel('Tx') private txModel: any) {}

    async queryTokensList(query: AssetListReqDto): Promise<ListStruct<AssetListResDto[]>> {
        const { pageNum, pageSize, useCount } = query;
        let count: number;
        const b: IAssetStruct[] = await (this.tokensModel as any).findList(pageNum, pageSize);
        if (useCount) {
            count = await (this.tokensModel as any).findCount();
        }
        const res: AssetListResDto[] = b.map((b) => {
            return new AssetListResDto(b);
        });
        return new ListStruct(res, pageNum, pageSize, count);
    }

    // tokens/:symbol
    async queryTokenDetail(query: AssetDetailReqDto): Promise<AssetDetailResDto | null> {
        const { symbol } = query;
        let result: AssetDetailResDto | null = null;
        let data:any = {};
        let token_db = await (this.tokensModel as any).findOneBySymbol(symbol);
        token_db = JSON.parse(JSON.stringify(token_db));
        if (token_db) {
            data = {
                name: token_db.name,
                owner: token_db.owner,
                total_supply: token_db.total_supply,
                initial_supply: token_db.initial_supply,
                max_supply: token_db.max_supply,
                mintable: token_db.mintable,
                scale: token_db.scale,
                min_unit: token_db.min_unit
            };
            result = new AssetDetailResDto(data);
        }
        return result;
    }
}
