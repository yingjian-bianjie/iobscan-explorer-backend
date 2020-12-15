import { Module } from '@nestjs/common';
import { MongooseModule } from '@nestjs/mongoose';
import { IritaModule } from './module/irita.module';
import { DenomModule } from './module/denom.module';
import { NftModule } from './module/nft.module';
import { BlockModule } from './module/block.module';
import { StatisticsModule } from './module/statistics.module';
import { TxModule } from './module/tx.module';
import { TxTaskModule } from './module/tx.task.module';
import { APP_FILTER, APP_PIPE } from '@nestjs/core';
import { HttpExceptionFilter } from './exception/HttpExceptionFilter';
import ValidationPipe from './pipe/validation.pipe';
import { ScheduleModule } from '@nestjs/schedule';
import { TasksService } from './task/task.service';

import { cfg } from './config/config';
import { TaskDispatchModule } from './module/task.dispatch.module';
import { DenomTaskModule } from './module/denom.task.module';
import { NftTaskModule } from './module/nft.task.module';
import { ValidatorTaskModule } from './module/validator.task.module';
import { ValidatorModule } from './module/validator.module';

import { IdentityTaskModule } from './module/identity.task.module';
import { IdentityModule } from './module/identity.module';
import { DistributionModule } from './module/distribution.module';
import {StakingValidatorTaskModule} from "./module/staking.validator.task.module";
import {ParametersTaskModule} from "./module/parameters.taskModule";
import {ProfilerModule} from "./module/profiler.module";
import {StakingModule} from "./module/staking.module";
import { TokensTaskModule } from "./module/tokens.task.module";
import { AssetModule } from './module/asset.module';

const url: string = `mongodb://${cfg.dbCfg.user}:${cfg.dbCfg.psd}@${cfg.dbCfg.dbAddr}/${cfg.dbCfg.dbName}`;
const params = {
    imports: [
        MongooseModule.forRoot(url),
        ScheduleModule.forRoot(),
        IritaModule,
        DenomModule,
        NftModule,
        BlockModule,
        StatisticsModule,
        TaskDispatchModule,
        DenomTaskModule,
        NftTaskModule,
        ValidatorTaskModule,
        ValidatorModule,
        TxModule,
        TxTaskModule,
        IdentityTaskModule,
        IdentityModule,
        DistributionModule,
        StakingValidatorTaskModule,
        ParametersTaskModule,
        ProfilerModule,
        StakingModule,
        TokensTaskModule,
        AssetModule
    ],
    providers:<any> [
        {
            provide: APP_FILTER,
            useClass: HttpExceptionFilter,
        },
        {
            provide: APP_PIPE,
            useClass: ValidationPipe,
        }
    ],
};
params.providers.push(TasksService);

// if (cfg.env !== 'development') {
//     params.providers.push(TasksService);
// }

@Module(params)
export class AppModule {
}
