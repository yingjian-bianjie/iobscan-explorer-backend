import { Controller, Get, Param, Query } from '@nestjs/common';
import { GovService } from '../service/gov.service';
import { Result } from '../api/ApiResult';
import { ListStruct } from '../api/ApiResult';
import { 
    proposalsReqDto,
    ProposalDetailReqDto,
    proposalsVoterReqDto
} from '../dto/gov.dto';
import { 
    govProposalResDto,
    govProposalDetailResDto,
    govProposalVoterResDto,
    govProposalDepositorResDto
} from '../dto/gov.dto';

import { ApiTags } from '@nestjs/swagger';

@ApiTags('Gov')
@Controller('gov')
export class GovController {
    constructor(private readonly govService: GovService) {
    }

    @Get('proposals')
    async getProposals(@Query() query: proposalsReqDto): Promise<Result<ListStruct<govProposalResDto>>> {
        const queryProposals = await this.govService.getProposals(query)
        return new Result<ListStruct<govProposalResDto>>(queryProposals)
    }

    @Get('proposals/:id')
    async getProposalDetail(@Param() param: ProposalDetailReqDto): Promise<Result<govProposalDetailResDto>> {
        const queryProposalDetail = await this.govService.getProposalDetail(param)
        return new Result<govProposalDetailResDto>(queryProposalDetail)
    }

    @Get('proposals/:id/voter')
    async getProposalVoter(@Param() param: ProposalDetailReqDto,@Query() query: proposalsVoterReqDto): Promise<Result<ListStruct<govProposalVoterResDto>>> {
        const queryProposalsVoter = await this.govService.getProposalsVoter(param,query)
        return new Result<ListStruct<govProposalVoterResDto>>(queryProposalsVoter)
    }

    @Get('proposals/:id/depositor')
    async getProposalDepositor(@Param() param: ProposalDetailReqDto,@Query() query: proposalsReqDto): Promise<Result<ListStruct<govProposalDepositorResDto>>> {
        const queryProposalsDepositor = await this.govService.getProposalsDepositor(param,query)
        return new Result<ListStruct<govProposalDepositorResDto>>(queryProposalsDepositor)
    }

}
