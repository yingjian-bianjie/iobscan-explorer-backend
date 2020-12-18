import { Controller, Get,Query } from '@nestjs/common';
import { Result } from '../api/ApiResult';
import { StatisticsResDto,NetworkStatisticsResDto } from '../dto/statistics.dto';
import { StatisticsService } from '../service/statistics.service';
import { ApiTags,ApiQuery } from '@nestjs/swagger';

@ApiTags('Statistics')
@Controller('statistics')
export class StatisticsController {
    constructor(private readonly statisticsService: StatisticsService) {
    }

    @Get('/db')
    @ApiQuery({
        name: 'params'
    })
    async queryStatistics(@Query() query:string): Promise<Result<StatisticsResDto>> {
        const data: StatisticsResDto = await this.statisticsService.queryStatistics(query);
        return new Result<StatisticsResDto>(data);
    }

    @Get('/network')
    @ApiQuery({
        name: 'params'
    })
    async queryNetworkStatistics(@Query() query:string): Promise<Result<NetworkStatisticsResDto>> {
        const data: NetworkStatisticsResDto = await this.statisticsService.queryNetworkStatistics(query);
        return new Result<NetworkStatisticsResDto>(data);
    }
}