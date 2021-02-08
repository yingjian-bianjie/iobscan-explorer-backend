import {Controller, Query, Get, Param} from '@nestjs/common';
import {ApiTags} from '@nestjs/swagger';
import {ListStruct, Result} from "../api/ApiResult";
import StakingService from "../service/staking.service";
import {
    // ValCommissionRewReqDto,
    CommissionInfoReqDto,
    // ValCommissionRewResDto,
    CommissionInfoResDto,
    ValidatorDelegationsReqDto,
    ValidatorDelegationsQueryReqDto,
    ValidatorDelegationsResDto,
    ValidatorUnBondingDelegationsReqDto,
    ValidatorUnBondingDelegationsQueryReqDto,
    ValidatorUnBondingDelegationsResDto,
    allValidatorReqDto,
    stakingValidatorResDto,
    ValidatorDetailAddrReqDto,
    ValidatorDetailResDto,
    AccountAddrReqDto,
    AccountAddrResDto,
    DelegatorsDelegationsReqDto,
    DelegatorsDelegationsResDto,
    DelegatorsDelegationsParamReqDto,
    DelegatorsUndelegationsReqDto,
    DelegatorsUndelegationsResDto,
    DelegatorsUndelegationsParamReqDto,
    ValidatorVotesResDto,
    ValidatorDepositsResDto
} from "../dto/staking.dto";

@ApiTags('Staking')
@Controller('staking')
export class StakingController {
    constructor(private readonly stakingService: StakingService) {
    }

    @Get('/commission_info')
    async getAllValCommissionInfo(@Query()q: CommissionInfoReqDto): Promise<Result<ListStruct<CommissionInfoResDto>>> {
        const allValCommissionData: ListStruct<CommissionInfoResDto> = await this.stakingService.getAllValCommission(q)
        return new Result<ListStruct<CommissionInfoResDto>>(allValCommissionData)
    }

    @Get('/validators/:address/delegations')
    async getValidatorDelegations(@Param()p: ValidatorDelegationsReqDto,@Query()q: ValidatorDelegationsQueryReqDto): Promise<Result<ListStruct<ValidatorDelegationsResDto>>> {
        const validatorDelegations = await this.stakingService.getValidatorDelegationList(p,q)
        return new Result<ListStruct<ValidatorDelegationsResDto>>(validatorDelegations)
    }

    @Get('/validators/:address/unbonding-delegations')
    async getValidatorUnBondingDelegations(@Param()p: ValidatorUnBondingDelegationsReqDto,@Query()q: ValidatorUnBondingDelegationsQueryReqDto): Promise<Result<ListStruct<ValidatorUnBondingDelegationsResDto>>> {
        const validatorUnBondingDelegations = await this.stakingService.getValidatorUnBondingDelegations(p,q)
        return new Result<ListStruct<ValidatorUnBondingDelegationsResDto>>(validatorUnBondingDelegations)
    }

    @Get('/validators')
    async getValidators(@Query()q: allValidatorReqDto): Promise<Result<ListStruct<stakingValidatorResDto>>> {
        const queryValidators = await this.stakingService.getValidatorsByStatus(q)
        return new Result<ListStruct<stakingValidatorResDto>>(queryValidators)
    }

    @Get('/validators/:address')
    async getValidatorDetail(@Param()q: ValidatorDetailAddrReqDto): Promise<Result<ValidatorDetailResDto>> {
        const queryValidatorDetail = await this.stakingService.getValidatorDetail(q)
        return new Result<ValidatorDetailResDto>(queryValidatorDetail)
    }

    @Get('/account/:address')
    async getAddressAccount(@Param()q: AccountAddrReqDto): Promise<Result<AccountAddrResDto>> {
        const addressAccount = await this.stakingService.getAddressAccount(q)
        return new Result<AccountAddrResDto>(addressAccount)
    }

    @Get('/delegators/:delegatorAddr/delegations')
    async getDelegatorsDelegations(@Param()p: DelegatorsDelegationsParamReqDto,@Query()q: DelegatorsDelegationsReqDto): Promise<Result<ListStruct<DelegatorsDelegationsResDto>>> {
        const delegatorsDelegations = await this.stakingService.getDelegatorsDelegations(p,q)
        return new Result<ListStruct<DelegatorsDelegationsResDto>>(delegatorsDelegations)
    }

    @Get('/delegators/:delegatorAddr/unbonding_delegations')
    async getDelegatorsUndelegations(@Param()p: DelegatorsUndelegationsParamReqDto,@Query()q: DelegatorsUndelegationsReqDto): Promise<Result<ListStruct<DelegatorsUndelegationsResDto>>> {
        const delegatorsUndelegations = await this.stakingService.getDelegatorsUndelegations(p,q)
        return new Result<ListStruct<DelegatorsUndelegationsResDto>>(delegatorsUndelegations)
    }

    @Get('/validators/:address/votes')
    async getValidatorVotes(@Param()p: ValidatorDelegationsReqDto,@Query()q: ValidatorDelegationsQueryReqDto): Promise<Result<ListStruct<ValidatorVotesResDto>>> {
        const validatorVotes = await this.stakingService.getValidatorVotesList(p,q)
        return new Result<ListStruct<ValidatorVotesResDto>>(validatorVotes)
    }

    @Get('/validators/:address/deposit')
    async getValidatorDeposits(@Param()p: ValidatorDelegationsReqDto,@Query()q: ValidatorDelegationsQueryReqDto): Promise<Result<ListStruct<ValidatorDepositsResDto>>> {
        const validatorDeposits = await this.stakingService.getValidatorDepositsList(p,q)
        return new Result<ListStruct<ValidatorDepositsResDto>>(validatorDeposits)
    }
    
}
