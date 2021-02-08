import { Test } from '@nestjs/testing';
import { AppModule } from '../app.module';
import { ParameterService } from "./parameter.service";
import { 
    ParametersListReqDto,
} from '../dto/parameter.dto';

describe('parameterController', () => {
    let parameterService : ParameterService
    beforeEach(async () => {
        const module = await Test.createTestingModule({
            imports:[
                AppModule
            ]
        }).compile();
        parameterService = module.get<ParameterService>(ParameterService);
    });

    describe('getProposals', () => {
        it('should return proposals array', async () => {
            let req:ParametersListReqDto = {}
            let data:any = await parameterService.queryParametersList(req);
            expect(data).toBeDefined();
        });
    });
})

