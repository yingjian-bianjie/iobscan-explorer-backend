import {
  ITxsQuery,
  ITxsWithAddressQuery,
  ITxsWithContextIdQuery,
} from '../types/schemaTypes/tx.interface';
import { TxStatus } from '../constant';
import { ITxsQueryParams } from '../types/tx.interface';
import {
  stakingTypes,
  declarationTypes,
  govTypes,
  coinswapTypes
} from '../helper/txTypes.helper';
import Cache from '../helper/cache';

export function txListParamsHelper(query: ITxsQuery){
  const queryParameters: ITxsQueryParams = {};
  if (query.type && query.type.length) {
      queryParameters['msgs.type'] = query.type;
  } else {
      queryParameters['msgs.type'] = {'$in':stakingTypes()};
  }
  if (query.status && query.status.length) {
      switch (query.status) {
          case '1':
              queryParameters.status = TxStatus.SUCCESS;
              break;
          case '2':
              queryParameters.status = TxStatus.FAILED;
              break;
      }
  }
  if (query.address && query.address.length) {
      queryParameters['addrs'] = { $elemMatch: { $eq: query.address } };
  }
  if ((query.beginTime && query.beginTime.length) || (query.endTime && query.endTime.length)) {
      queryParameters.time = {};
  }
  if (query.beginTime && query.beginTime.length) {
      queryParameters.time.$gte = Number(query.beginTime);
  }
  if (query.endTime && query.endTime.length) {
      queryParameters.time.$lte = Number(query.endTime);
  }
  return queryParameters
}

export function StakingTxListParamsHelper(query: ITxsQuery){
  const queryParameters: any = {};
  if (query.type && query.type.length) {
      queryParameters['msgs.type'] = query.type;
  } else {
      queryParameters['msgs.type'] = {'$in':stakingTypes()};
  }
  if (query.status && query.status.length) {
      switch (query.status) {
          case '1':
              queryParameters.status = TxStatus.SUCCESS;
              break;
          case '2':
              queryParameters.status = TxStatus.FAILED;
              break;
      }
  }
  if (query.address && query.address.length) {
      queryParameters['addrs'] = { $elemMatch: { $eq: query.address } };
  }
  if ((query.beginTime && query.beginTime.length) || (query.endTime && query.endTime.length)) {
      queryParameters.time = {};
  }
  if (query.beginTime && query.beginTime.length) {
      queryParameters.time.$gte = Number(query.beginTime);
  }
  if (query.endTime && query.endTime.length) {
      queryParameters.time.$lte = Number(query.endTime);
  }
  return queryParameters
}

export function CoinswapTxListParamsHelper(query: ITxsQuery){
  const queryParameters: any = {};
  const { type } = query;
  if (query.type && query.type.length) {
      const typeArr = type.split(",");
      queryParameters['msgs.type'] = {
          $in: typeArr
      }
  } else {
      queryParameters['msgs.type'] = {'$in':coinswapTypes()};
  }
  if (query.status && query.status.length) {
      switch (query.status) {
          case '1':
              queryParameters.status = TxStatus.SUCCESS;
              break;
          case '2':
              queryParameters.status = TxStatus.FAILED;
              break;
      }
  }
  if (query.address && query.address.length) {
      queryParameters['addrs'] = { $elemMatch: { $eq: query.address } };
  }
  if ((query.beginTime && query.beginTime.length) || (query.endTime && query.endTime.length)) {
      queryParameters.time = {};
  }
  if (query.beginTime && query.beginTime.length) {
      queryParameters.time.$gte = Number(query.beginTime);
  }
  if (query.endTime && query.endTime.length) {
      queryParameters.time.$lte = Number(query.endTime);
  }
  return queryParameters
}

export function DeclarationTxListParamsHelper(query: ITxsQuery){
  const queryParameters: any = {};
  if (query.type && query.type.length) {
      queryParameters['msgs.type'] = query.type;
  } else {
      queryParameters['msgs.type'] = { '$in': declarationTypes() };
  }
  if (query.status && query.status.length) {
      switch (query.status) {
          case '1':
              queryParameters.status = TxStatus.SUCCESS;
              break;
          case '2':
              queryParameters.status = TxStatus.FAILED;
              break;
      }
  }
  if (query.address && query.address.length) {
      queryParameters['addrs'] = { $elemMatch: { $eq: query.address } };
  }
  if ((query.beginTime && query.beginTime.length) || (query.endTime && query.endTime.length)) {
      queryParameters.time = {};
  }
  if (query.beginTime && query.beginTime.length) {
      queryParameters.time.$gte = Number(query.beginTime);
  }
  if (query.endTime && query.endTime.length) {
      queryParameters.time.$lte = Number(query.endTime);
  }
  return queryParameters
}

export function GovTxListParamsHelper(query: ITxsQuery){
  const queryParameters: any = {};
  if (query.type && query.type.length) {
      queryParameters['msgs.type'] = query.type;
  } else {
      queryParameters['msgs.type'] = { '$in': govTypes() };
  }
  if (query.status && query.status.length) {
      switch (query.status) {
          case '1':
              queryParameters.status = TxStatus.SUCCESS;
              break;
          case '2':
              queryParameters.status = TxStatus.FAILED;
              break;
      }
  }
  if (query.address && query.address.length) {
      queryParameters['addrs'] = { $elemMatch: { $eq: query.address } };
  }
  if ((query.beginTime && query.beginTime.length) || (query.endTime && query.endTime.length)) {
      queryParameters.time = {};
  }
  if (query.beginTime && query.beginTime.length) {
      queryParameters.time.$gte = Number(query.beginTime);
  }
  if (query.endTime && query.endTime.length) {
      queryParameters.time.$lte = Number(query.endTime);
  }
  return queryParameters
}

export function TxListEdgeParamsHelper(types, gt_height, status, address, include_event_addr){
  const queryParameters: any = {};
  if (types && types.length) {
      queryParameters['msgs.type'] = {'$in':types.split(',')};
  }
  if (gt_height) {
      queryParameters['height'] = {'$gt':gt_height};
  }
  if (status || status === 0) {
      queryParameters['status'] = status;
  }
  if (include_event_addr && include_event_addr == true && address && address.length) {
      queryParameters['$or'] = [
          { 'events.attributes.value': address },
          { 'addrs': { $elemMatch: {'$in': address.split(',')} }}
      ]
  } else if (address && address.length) {
      queryParameters['addrs'] = { $elemMatch: { '$in': address.split(',') } };
  }
  return queryParameters
}

export function TxWithAddressParamsHelper(query: ITxsWithAddressQuery){
  let queryParameters: any = {};
  if (query.address && query.address.length) {
      queryParameters = {
          // $or:[
          // 	{"from":query.address},
          // 	{"to":query.address},
          // 	{"signer":query.address},
          // ],
          addrs: { $elemMatch: { $eq: query.address } },
      };
  }
  if (query.type && query.type.length) {
      queryParameters['msgs.type'] = query.type;
  } else {
      // queryParameters.$or = [{ 'msgs.type': filterExTxTypeRegExp() }];
      queryParameters['msgs.type'] = {
          $in: Cache.supportTypes || []
      }
  }
  if (query.status && query.status.length) {
      switch (query.status) {
          case '1':
              queryParameters.status = TxStatus.SUCCESS;
              break;
          case '2':
              queryParameters.status = TxStatus.FAILED;
              break;
      }
  }
  return queryParameters
}

export function TxWithContextIdParamsHelper(query: ITxsWithContextIdQuery){
  let queryParameters: any = {};
  if (query.contextId && query.contextId.length) {
      queryParameters = {
          $or: [
              { 'events.attributes.value': query.contextId },
              { 'msgs.msg.ex.request_context_id': query.contextId },
              { 'msgs.msg.request_context_id': query.contextId },
          ],
      };
  }
  if (query.type && query.type.length) {
      queryParameters['msgs.type'] = query.type;
  } else {
      // queryParameters.$or = [{ 'msgs.type': filterExTxTypeRegExp() }];
      queryParameters['msgs.type'] = {
          $in: Cache.supportTypes || []
      }
  }

  if (query.status && query.status.length) {
      switch (query.status) {
          case '1':
              queryParameters.status = TxStatus.SUCCESS;
              break;
          case '2':
              queryParameters.status = TxStatus.FAILED;
              break;
      }
  }
  return queryParameters
}