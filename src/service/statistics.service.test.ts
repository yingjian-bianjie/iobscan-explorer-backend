import { Test } from '@nestjs/testing';
import { AppModule } from '../app.module';
import { StatisticsService } from './statistics.service';


describe('statistics module', () => {
    let statisticsService: StatisticsService;

    beforeEach(async () => {
        const module = await Test.createTestingModule({
            imports: [
                AppModule,
            ],
        }).compile();
        statisticsService = module.get<StatisticsService>(StatisticsService);
    });

    describe('queryStatistics', () => {
        it('should return statistics detail', async () => {
            const data: any = await statisticsService.queryStatistics('');
            expect(data).toBeDefined();
        });
    });

    describe('queryLatestHeight', () => {
        it('should return latest block height', async () => {
            const data: {height:number,latestBlockTime:number} | null = await statisticsService.queryLatestHeightAndTimeAndValidator();
            expect(data).toBeDefined();
        });
    });

    describe('queryAvgBlockTime', () => {
        it('should return block average time', async () => {
            const data: number | null = await statisticsService.queryAvgBlockTime();
            expect(data).toBeDefined();
        });
    });

    describe('queryAssetCount', () => {
        it('should return asset count', async () => {
            const data: number | null = await statisticsService.queryAssetCount();
            expect(data).toBeDefined();
        });
    });


    describe('queryConsensusValidatorCount', () => {
        it('should return validator count', async () => {
            const data: number | null = await statisticsService.queryConsensusValidatorCount();
            expect(data).toBeDefined();
        });
    });

    describe('queryTxCount', () => {
        it('should return tx count', async () => {
            const data: any = await statisticsService.queryTxCount();
            expect(data).toBeDefined();
        });
    });

    describe('queryNetworkStatistics', () => {
        it('should return NetworkStatistics', async () => {
            const data: any = await statisticsService.queryNetworkStatistics('200');
            expect(data).toBeDefined();
        });
    });
});

