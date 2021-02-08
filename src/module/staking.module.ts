import {Module} from '@nestjs/common';
import {MongooseModule} from '@nestjs/mongoose';
import {StakingValidatorSchema} from "../schema/staking.validator.schema";
import StakingService from "../service/staking.service";
import {StakingHttp} from "../http/lcd/staking.http";
import {ProfilerSchema} from "../schema/profiler.schema";
import {StakingController} from "../controller/staking.controller";
import {ParametersSchema} from "../schema/parameters.schema";
import {TxSchema} from "../schema/tx.schema";
import {DistributionHttp} from "../http/lcd/distribution.http";
import { TokensSchema } from "../schema/tokens.schema";
import { ProposalSchema } from '../schema/proposal.schema';
@Module({
    imports: [
        MongooseModule.forFeature([
            {
                name: 'StakingSyncValidators',
                schema: StakingValidatorSchema,
                collection: 'ex_staking_validator'
            },
            {
                name: 'Profiler',
                schema: ProfilerSchema,
                collection: 'ex_profiler'
            },
            {
                name: 'Parameters',
                schema: ParametersSchema,
                collection: 'ex_sync_parameters'
            },
            {
                name: 'Tx',
                schema: TxSchema,
                collection: 'sync_tx'
            },
            {
                name: 'Tokens',
                schema: TokensSchema,
                collection: 'ex_tokens'
            },
            {
                name: 'Proposal',
                schema: ProposalSchema,
                collection: 'ex_sync_ proposal'
            }
        ])
    ],
    providers: [StakingService, StakingHttp,DistributionHttp],
    controllers: [StakingController],
})
export class StakingModule {
}
