import {Module } from '@nestjs/common';
import {DenomTaskService} from '../task/denom.task.service';
import { MongooseModule } from '@nestjs/mongoose';
import {DenomSchema} from '../schema/denom.schema';
import { DenomHttp } from '../http/lcd/denom.http';
import { TxSchema } from '../schema/tx.schema';
import { SyncTaskSchema } from '../schema/sync.task.schema';
@Module({
    imports:[
        MongooseModule.forFeature([{
            name: 'Denom',
            schema: DenomSchema,
            collection: 'ex_sync_denom'
        },{
            name: 'Tx',
            schema: TxSchema,
            collection: 'sync_tx',
        },{
            name: 'SyncTask',
            schema: SyncTaskSchema,
            collection: 'sync_task'
        }]),
    ],
    providers:[DenomTaskService, DenomHttp],
    exports:[DenomTaskService]
})
export class DenomTaskModule{}