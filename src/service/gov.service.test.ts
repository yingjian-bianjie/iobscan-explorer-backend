import { Test } from '@nestjs/testing';
import { AppModule } from '../app.module';
import { GovService } from "./gov.service";
import {
    proposalsReqDto,
    ProposalDetailReqDto,
    proposalsVoterReqDto
} from '../dto/gov.dto';

describe('govController', () => {
    let govService : GovService
    beforeEach(async () => {
        const module = await Test.createTestingModule({
            imports:[
                AppModule
            ]
        }).compile();
        govService = module.get<GovService>(GovService);
    });

    describe('getProposals', () => {
        it('should return proposals array', async () => {
            let req:proposalsReqDto = {
                pageNum:1,
                pageSize:10,
                useCount:true
            };
            let data:any = await govService.getProposals(req);
            expect(data).toBeDefined();
        });
    });

    describe('getProposalDetail', () => {
        it('should return proposalsDetail', async () => {
            let req:ProposalDetailReqDto = {
                id: 1
            };
            let data:any = await govService.getProposalDetail(req);
            expect(data).toBeDefined();
        });
    });

    describe('getProposalsVoter', () => {
        it('should return ProposalsVoter', async () => {
            let query:proposalsVoterReqDto = {
                pageNum:1,
                pageSize:2,
                useCount: true,
                voterType:'all'
            };
            let param: ProposalDetailReqDto = {
                id: 1
            }
            let data: any = await govService.getProposalsVoter(param,query);
            expect(data).toBeDefined();
        });
    });

    describe('getProposalsDepositor', () => {
        it('should return proposalsDepositor array', async () => {
            let param: ProposalDetailReqDto = {
                id: 1
            };
            let req:proposalsReqDto = {
                pageNum:1,
                pageSize:10,
                useCount:true
            };
            let data:any = await govService.getProposalsDepositor(param,req);
            expect(data).toBeDefined();
        });
    });
})
