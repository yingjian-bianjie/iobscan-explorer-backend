import { Module } from '@nestjs/common';
import { TxController } from '../controller/tx.controller';
import { TxService } from '../service/tx.service';
import { MongooseModule } from '@nestjs/mongoose';
import { TxSchema } from '../schema/tx.schema';
import { TxTypeSchema } from '../schema/txType.schema';
import { DenomSchema } from '../schema/denom.schema';
import { NftSchema } from '../schema/nft.schema';
import { TxEvmSchema } from '../schema/txEvmSchema';
import { EvmContractConfigSchema } from '../schema/evmContractConfig.schema';
import { IdentitySchema } from '../schema/identity.schema';
import {StakingValidatorSchema} from "../schema/staking.validator.schema";
import { ProposalSchema } from '../schema/proposal.schema';
import { GovHttp } from "../http/lcd/gov.http";
import { StatisticsSchema } from '../schema/statistics.schema';
@Module({
    imports:[
        MongooseModule.forFeature([{
            name: 'Tx',
            schema: TxSchema,
            collection: 'sync_tx'
        },
        {
            name: 'TxType',
            schema: TxTypeSchema,
            collection: 'ex_tx_type'
        },
        {
            name: 'Denom',
            schema: DenomSchema,
            collection: 'ex_sync_denom'
        },
        {
            name: 'Identity',
            schema: IdentitySchema,
            collection: 'sync_identity'
        },
        {
            name: 'Nft',
            schema: NftSchema,
            collection: 'ex_sync_nft'
        },
        {
            name: 'TxEvm',
            schema: TxEvmSchema,
            collection: 'ex_sync_tx_evm'
        },{
            name: 'EvmContractConfig',
            schema: EvmContractConfigSchema,
            collection: 'ex_evm_contracts_config'
        },
        {
            name: 'StakingValidator',
            schema: StakingValidatorSchema,
            collection: 'ex_staking_validator'
        },
        {
            name: 'Proposal',
            schema: ProposalSchema,
            collection: 'ex_sync_proposal'
        },{
                name: 'Statistics',
                schema: StatisticsSchema,
                collection: 'ex_statistics'
            }])
    ],
    providers:[TxService,GovHttp],
    controllers:[TxController],
})
export class TxModule{}
