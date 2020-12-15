import {IsString, IsInt, Length, Min, Max, IsOptional, Equals, MinLength, ArrayNotEmpty} from 'class-validator';
import {ApiProperty, ApiPropertyOptional} from '@nestjs/swagger';
import {BaseReqDto, BaseResDto, PagingReqDto} from './base.dto';
import {ApiError} from '../api/ApiResult';
import {ErrorCodes} from '../api/ResultCodes';
import {IBindTx} from '../types/tx.interface';

/************************   request dto   ***************************/


/************************   response dto   ***************************/
//txs response dto
export class NetworkResDto extends BaseResDto{
    network_id: string;
    network_name: string;
    uri: string;
    is_main: boolean;

    constructor(value) {
        super();
        const { network_id, network_name, uri, is_main } = value;
        this.network_id = network_id;
        this.network_name = network_name;
        this.uri = uri;
        this.is_main = is_main || false;
    }

    static bundleData(value: any): NetworkResDto[] {
        let data: NetworkResDto[] = [];
        data = value.map((v: any) => {
            return new NetworkResDto(v);
        });
        return data;
    }
}

export class TokensResDto extends BaseResDto {
    symbol:string;
    min_unit:string;
    scale:number;
    is_main_token:boolean;
    constructor(value) {
        super();
        this.symbol = value.symbol;
        this.min_unit = value.min_unit;
        this.scale = value.scale;
        this.is_main_token = value.is_main_token;
    }
    static bundleData(value: any): TokensResDto[] {
        let data: TokensResDto[] = [];
        data = value.map((v: any) => {
            return new TokensResDto(v);
        });
        return data;
    }

}

export class StatusResDto {
    is_follow: boolean;
    constructor(value) {
        this.is_follow = value;
    }
}