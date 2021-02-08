import { Module } from '@nestjs/common';
import { MongooseModule } from '@nestjs/mongoose';
import {ParametersSchema} from "../schema/parameters.schema";
import { TxSchema } from '../schema/tx.schema';
import { ProposalSchema } from '../schema/proposal.schema';
import { ProposalDetailSchema } from '../schema/proposal.detail.schema';
import { ProposalTaskService } from "../task/proposal.service";
import { GovHttp } from "../http/lcd/gov.http";
import { StakingValidatorSchema } from "../schema/staking.validator.schema";
import {StakingHttp} from "../http/lcd/staking.http";
@Module({
    imports:[
        MongooseModule.forFeature([
            {
                name:'ParametersTask',
                schema: ParametersSchema,
                collection:'ex_sync_parameters'
            },
            {
                name: 'Tx',
                schema: TxSchema,
                collection: 'sync_tx'
            },
            {
                name: 'Proposal',
                schema: ProposalSchema,
                collection: 'ex_sync_ proposal'
            },
            {
                name: 'StakingSyncValidators',
                schema: StakingValidatorSchema,
                collection: 'ex_staking_validator'
            },
            {
                name: 'ProposalDetail',
                schema: ProposalDetailSchema,
                collection: 'ex_sync_ proposal_details'
            }
        ])
    ],
    providers:[ProposalTaskService,GovHttp,StakingHttp],
    exports: [ProposalTaskService]
})
export class ProposalTaskModule{}
