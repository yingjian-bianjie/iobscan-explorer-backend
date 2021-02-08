import {Module } from '@nestjs/common';
import {ParameterController} from '../controller/parameter.controller';
import {ParameterService} from '../service/parameter.service';
import { MongooseModule } from '@nestjs/mongoose';
import {ParametersSchema} from "../schema/parameters.schema";
@Module({
    imports:[
        MongooseModule.forFeature([
            {
                name:'ParametersTask',
                schema: ParametersSchema,
                collection:'ex_sync_parameters'
            }
        ]),
    ],
    providers:[ParameterService],
    controllers:[ParameterController],
})
export class ParameterModule{}
