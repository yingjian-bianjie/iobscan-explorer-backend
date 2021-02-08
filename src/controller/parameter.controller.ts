import { Controller, Get, Param, Query } from '@nestjs/common';
import { ParameterService } from '../service/parameter.service';
import { Result } from '../api/ApiResult';
import { ListStruct } from '../api/ApiResult';
import { 
    ParametersListReqDto,
} from '../dto/parameter.dto';
import { 
    ParametersListResDto
} from '../dto/parameter.dto';

import { ApiTags } from '@nestjs/swagger';

@ApiTags('Parameter')
@Controller('parameter')
export class ParameterController {
    constructor(private readonly paramterService: ParameterService) {
    }

    @Get('/')
    async queryParametersList(@Query() query: ParametersListReqDto): Promise<Result<ParametersListResDto[]>> {
        const data: ParametersListResDto[] = await this.paramterService.queryParametersList(query);
        return new Result<ParametersListResDto[]>(data);
    }
}
