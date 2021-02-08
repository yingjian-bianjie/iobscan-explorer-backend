import { Injectable } from '@nestjs/common';
import { Model } from 'mongoose';
import { InjectModel } from '@nestjs/mongoose';
import { ListStruct } from '../api/ApiResult';
import { 
    ParametersListReqDto,
} from '../dto/parameter.dto';
import { 
    ParametersListResDto
} from '../dto/parameter.dto';

@Injectable()
export class ParameterService {

    constructor(
        @InjectModel('ParametersTask') private parametersTaskModel: Model<any>) { }
    
    async queryParametersList(query: ParametersListReqDto): Promise<ParametersListResDto[]> {
            return await (this.parametersTaskModel as any).queryParams(query)
        }

}
