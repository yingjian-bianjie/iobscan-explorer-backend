/* 一个通用的module 一些不方便归类模块的api可以放在这里 */

import {Module} from '@nestjs/common';
import {IritaController} from '../controller/irita.controller';
import {IritaService} from '../service/irita.service';
import {MongooseModule} from '@nestjs/mongoose';
import {NetworkSchema} from '../schema/network.schema';
import {TokensSchema} from "../schema/tokens.schema";
import {TokensTaskService} from "../task/tokens.service";
import {TokensHttp} from "../http/lcd/tokens.http";
import {ParametersSchema} from "../schema/parameters.schema";
import {ParametersTaskService} from "../task/parameters.task.service";
import {StakingHttp} from "../http/lcd/staking.http";
import { TxSchema } from '../schema/tx.schema';
import { SyncTaskSchema } from '../schema/sync.task.schema';
import {GovHttp} from "../http/lcd/gov.http";
@Module({
    imports: [
        MongooseModule.forFeature([
            {
                name: 'Network',
                schema: NetworkSchema,
                collection: 'ex_network'
            },
            {
                name: 'Tokens',
                schema: TokensSchema,
                collection: 'ex_tokens'
            },
            {
                name:'ParametersTask',
                schema: ParametersSchema,
                collection:'ex_sync_parameters'
            },{
                name: 'Tx',
                schema: TxSchema,
                collection: 'sync_tx'
            },{
                name: 'SyncTask',
                schema: SyncTaskSchema,
                collection: 'sync_task'
            }
        ])
    ],
    providers: [IritaService, TokensTaskService, TokensHttp,ParametersTaskService,StakingHttp,GovHttp],
    controllers: [IritaController],
})
export class IritaModule {
}
