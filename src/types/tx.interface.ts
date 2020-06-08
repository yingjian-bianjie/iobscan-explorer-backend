import { IQueryBase } from '.';
import { Document } from 'mongoose';

export interface ITxsQueryParams extends IQueryBase {
	type?:string,
	status?:number,
	time?:{
		$gte?:Date,
		$lte?:Date,
	}
}

export interface IListStruct {
	data?: any[],
    count?: number
}
