import { Test } from '@nestjs/testing';
import { AppModule } from '../app.module';
import StakingService from "./staking.service";
import {
    CommissionInfoReqDto,
    ValidatorDelegationsReqDto,
    ValidatorUnBondingDelegationsReqDto,
    ValidatorUnBondingDelegationsQueryReqDto,
    allValidatorReqDto,
    ValidatorDetailAddrReqDto,
    AccountAddrReqDto,
    ValidatorDelegationsQueryReqDto
} from '../dto/staking.dto';

describe('stakingValidatorController', () => {
    let stakingService : StakingService
    beforeEach(async () => {
        const module = await Test.createTestingModule({
            imports:[
                AppModule
            ]
        }).compile();
        stakingService = module.get<StakingService>(StakingService);
    });

    describe('getAllValCommission', () => {
        it('should return an array', async () => {
            let req:CommissionInfoReqDto = {
                pageNum:1,
                pageSize:10,
                useCount:true
            };
            let data:any = await stakingService.getAllValCommission(req);
            expect(data).toBeDefined();
        });
    });

    describe('getValidatorDelegationList', () => {
        it('should return Validator Delegation array', async () => {
            let req:ValidatorDelegationsReqDto = {
                address:'iva1nfgkyn6ux5mvavhyk5aq8zgzushrecutlt5pr8',
            };
            let query:ValidatorUnBondingDelegationsQueryReqDto = {
                pageNum: 1,
                pageSize: 10,
                useCount: true
            };
            let data:any = await stakingService.getValidatorDelegationList(req,query);
            expect(data).toBeDefined();
        });
    });

    describe('getValidatorUnBondingDelegations', () => {
        it('should return Validator UnBonding Delegation array', async () => {
            let req:ValidatorUnBondingDelegationsReqDto = {
                address:'iva1nfgkyn6ux5mvavhyk5aq8zgzushrecutlt5pr8',
            };
            let query:ValidatorUnBondingDelegationsQueryReqDto = {
                pageNum: 1,
                pageSize: 10,
                useCount: true
            };
            let data:any = await stakingService.getValidatorUnBondingDelegations(req,query);
            expect(data).toBeDefined();
        });
    });

    describe('getValidatorsByStatus', () => {
        it('should return Validator array', async () => {
            let req:allValidatorReqDto = {
                status:'active',
                pageNum:1,
                pageSize:10,
                useCount:true
            };
            let data:any = await stakingService.getValidatorsByStatus(req);
            expect(data).toBeDefined();
        });
    });

    describe('getValidatorDetail', () => {
        it('should return Validator array', async () => {
            let req:ValidatorDetailAddrReqDto = {
                address:'iva1nfgkyn6ux5mvavhyk5aq8zgzushrecutlt5pr8',
            };
            let data:any = await stakingService.getValidatorDetail(req);
            if (data) {
                expect(data.operator_addr).toBe(req.address);
            }else{
                expect(data).toBeDefined();
            }

        });
    });

    describe('getAddressAccount', () => {
        it('should return Validator array', async () => {
            let req:AccountAddrReqDto = {
                address:'iaa1nfgkyn6ux5mvavhyk5aq8zgzushrecut267w7q',
            };
            let data:any = await stakingService.getAddressAccount(req);
            if (data) {
                expect(data.address).toBe(req.address);
            }else{
                expect(data).toBeDefined();
            }

        });
    });

    describe('getValidatorVotesList', () => {
        it('should return VotesList array', async () => {
            let param:ValidatorDelegationsReqDto = {
                address:'iaa1nfgkyn6ux5mvavhyk5aq8zgzushrecut267w7q',
            }
            let query:ValidatorDelegationsQueryReqDto = {
                pageNum:1,
                pageSize:10,
                useCount:true
            };
            let data:any = await stakingService.getValidatorVotesList(param,query);
            expect(data).toBeDefined();
        });
    });

    describe('getValidatorDepositsList', () => {
        it('should return DepositsList array', async () => {
            let param:ValidatorDelegationsReqDto = {
                address:'iva1nfgkyn6ux5mvavhyk5aq8zgzushrecutlt5pr8',
            }
            let query:ValidatorDelegationsQueryReqDto = {
                pageNum:1,
                pageSize:10,
                useCount:true
            };
            let data:any = await stakingService.getValidatorDepositsList(param,query);
            expect(data).toBeDefined();
        });
    });
})

