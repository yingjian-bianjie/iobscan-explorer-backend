import {HttpService, Injectable} from '@nestjs/common';
import {cfg} from '../../config/config';
import {Logger} from '../../logger';
import {
    AddressBalancesLcdDto,
    IconUriLcdDto,
    StakingValidatorDelegationLcdDto,
    StakingValidatorLcdDto, StakingValidatorParametersLcdDto,
    StakingValidatorSlashLcdDto, StakingValUnBondingDelLcdDto,
    DelegatorsDelegationLcdDto,
    DelegatorsUndelegationLcdDto,
    BondedTokensLcdDto,
    TokensStakingLcdToken
} from "../../dto/http.dto";
import { currentChain } from '../../constant/index'

@Injectable()

export class StakingHttp {
    async queryValidatorListFromLcd(status: string, pageNum: number, pageSize: number) {
        let validatorLcdUri;
        switch (cfg.currentChain) {
            case currentChain.iris:
                // iris
                // validatorLcdUri = `${cfg.serverCfg.lcdAddr}/staking/validators?status=${status}&pageNum=${pageNum}&pageSize=${pageSize}`
                validatorLcdUri = `${cfg.serverCfg.lcdAddr}/staking/validators?status=${status}&page=${pageNum}&limit=${pageSize}`
                break;
            case currentChain.cosmos:
                // cosmos
                validatorLcdUri = `${cfg.serverCfg.lcdAddr}/staking/validators?status=${status}&page=${pageNum}&limit=${pageSize}`
                break;
            default:
                break;
        }
        try {
            let stakingValidatorData: any = await new HttpService().get(validatorLcdUri).toPromise().then(result => result.data)
            if (stakingValidatorData && stakingValidatorData.result) {
                return StakingValidatorLcdDto.bundleData(stakingValidatorData.result);
            } else {
                Logger.warn('api-error:', 'there is no result of validators from lcd');
            }
        } catch (e) {
            Logger.warn(`api-error from ${validatorLcdUri}`, e)
        }
    }

    async queryValidatorFormSlashing(address_ica: string) {
        // const slashValidatorUri = `${cfg.serverCfg.lcdAddr}/slashing/validators/${validatorPubkey}/signing_info`
        //todo:hangtaishan 测试网api
        const slashValidatorUri = `${cfg.serverCfg.lcdAddr}/cosmos/slashing/v1beta1/signing_infos/${address_ica}`
        try {
            const stakingSlashValidatorData: any = await new HttpService().get(slashValidatorUri).toPromise().then(result => result.data)
            if (stakingSlashValidatorData && stakingSlashValidatorData.val_signing_info) {
                return new StakingValidatorSlashLcdDto(stakingSlashValidatorData.val_signing_info);
            } else {
                Logger.warn('api-error:', 'there is no result of validators from lcd');
            }
        } catch (e) {
            Logger.warn(`api-error from ${slashValidatorUri}`, e)
        }
    }

    async querySelfBondFromLcd(valOperatorAddr) {
        const selfBondUri = `${cfg.serverCfg.lcdAddr}/staking/validators/${valOperatorAddr}/delegations`
        try {
            const selfBondData: any = await new HttpService().get(selfBondUri).toPromise().then(result => result.data)
            if (selfBondData && selfBondData.result) {
                return StakingValidatorDelegationLcdDto.bundleData(selfBondData.result);
            } else {
                Logger.warn('api-error:', 'there is no result of validators from lcd');
            }
        } catch (e) {
            Logger.warn(`api-error from ${selfBondUri}`, e)
        }
    }

    async queryValidatorIcon(valIdentity) {
        const getIconUri = `${cfg.serverCfg.iconUri}?fields=pictures&key_suffix=${valIdentity || ''}`
        try {
            const valIconData: any = await new HttpService().get(getIconUri).toPromise().then(result => result.data)
            if (valIconData) {
                return new IconUriLcdDto(valIconData)
            } else {
                Logger.warn('api-error:', 'there is no result of validators from getIconUri');
            }

        } catch (e) {
            Logger.warn(`api-error from ${getIconUri}`, e)
        }
    }

    async queryParametersFromSlashing() {
        const parameterUri = `${cfg.serverCfg.lcdAddr}/slashing/parameters`
        try {
            const parameterData: any = await new HttpService().get(parameterUri).toPromise().then(result => result.data)
            if (parameterData && parameterData.result) {
                return new StakingValidatorParametersLcdDto(parameterData.result)
            } else {
                Logger.warn('api-error:', 'there is no result of validators from lcd');
            }

        } catch (e) {
            Logger.warn(`api-error from ${parameterUri}`, e)
        }
    }

    async queryValidatorDelegationsFromLcd(address) {
        // const getValidatorDelegationsUri = `${cfg.serverCfg.lcdAddr}/staking/validators/${address}/delegations`
        const getValidatorDelegationsUri = `${cfg.serverCfg.lcdAddr}/cosmos/staking/v1beta1/validators/${address}/delegations`
        try {
            let validatorDelegationsData: any = await new HttpService().get(getValidatorDelegationsUri).toPromise().then(result => result.data)
            if (validatorDelegationsData && validatorDelegationsData.delegation_responses) {
                return StakingValidatorDelegationLcdDto.bundleData(validatorDelegationsData.delegation_responses);
            } else {
                Logger.warn('api-error:', 'there is no result of validator delegations from lcd');
            }
        } catch (e) {
            Logger.warn(`api-error from ${getValidatorDelegationsUri}`, e)
        }
    }

    async queryValidatorUnBondingDelegations(address) {
        const getValidatorUnBondingDelUri = `${cfg.serverCfg.lcdAddr}/staking/validators/${address}/unbonding_delegations`
        try {
            let validatorUnBondingDelegationsData: any = await new HttpService().get(getValidatorUnBondingDelUri).toPromise().then(result => result.data)
            if (validatorUnBondingDelegationsData && validatorUnBondingDelegationsData.result) {
                return StakingValUnBondingDelLcdDto.bundleData(validatorUnBondingDelegationsData.result);
            } else {
                Logger.warn('api-error:', 'there is no result of validator unBonding delegations from lcd');
            }
        } catch (e) {
            Logger.warn(`api-error from ${getValidatorUnBondingDelUri}`, e)
        }
    }

    async queryBalanceByAddress(address) {
        const getBalancesUri = `${cfg.serverCfg.lcdAddr}/bank/balances/${address}`
        try {
            let addressBalancesData: any = await new HttpService().get(getBalancesUri).toPromise().then(result => result.data)
            if (addressBalancesData && addressBalancesData.result) {
                return AddressBalancesLcdDto.bundleData(addressBalancesData.result);
            } else {
                Logger.warn('api-error:', 'there is no result of validator unBonding delegations from lcd');
            }

        } catch (e) {
            Logger.warn(`api-error from ${getBalancesUri}`, e)
        }
    }

    async queryDelegatorsDelegationsFromLcd(address) {
        const getDelegatorsDelegationsUri = `${cfg.serverCfg.lcdAddr}/staking/delegators/${address}/delegations`
        try {
            const delegatorsDelegationsData: any = await new HttpService().get(getDelegatorsDelegationsUri).toPromise().then(result => result.data)
            if (delegatorsDelegationsData && delegatorsDelegationsData.result) {
                return new DelegatorsDelegationLcdDto(delegatorsDelegationsData);
            } else {
                Logger.warn('api-error:', 'there is no result of delegators delegations from lcd');
            }
        } catch (e) {
            Logger.warn(`api-error from ${getDelegatorsDelegationsUri}`, e)
        }
    }

    async queryDelegatorsUndelegationsFromLcd(address) {
        const getDelegatorsUndelegationsUri = `${cfg.serverCfg.lcdAddr}/staking/delegators/${address}/unbonding_delegations`
        try {
            const delegatorsUndelegationsData: any = await new HttpService().get(getDelegatorsUndelegationsUri).toPromise().then(result => result.data)
            if (delegatorsUndelegationsData && delegatorsUndelegationsData.result) {
                return new DelegatorsUndelegationLcdDto(delegatorsUndelegationsData);
            } else {
                Logger.warn('api-error:', 'there is no result of delegators delegations from lcd');
            }
        } catch (e) {
            Logger.warn(`api-error from ${getDelegatorsUndelegationsUri}`, e)
        }
    }

    static async getBondedTokens () {
        const BondedTokensUrl = `${cfg.serverCfg.lcdAddr}/staking/pool`
        try {
            let BondedTokens: any = await new HttpService().get(BondedTokensUrl).toPromise().then(result => result.data)
            if (BondedTokens && BondedTokens.result) {
                return new BondedTokensLcdDto(BondedTokens.result);
            } else {
                Logger.warn('api-error:', 'there is no result of bonded_tokens from lcd');
            }
        } catch (e) {
            Logger.warn(`api-error from ${BondedTokensUrl}`, e)
        }
    }

    async getStakingTokens() {
        const stakingTokensUrl = `${cfg.serverCfg.lcdAddr}/cosmos/staking/v1beta1/params`
        try {
            let stakingTokensData: any = await new HttpService().get(stakingTokensUrl).toPromise().then(result => result.data)
            if (stakingTokensData && stakingTokensData.params) {
                return new TokensStakingLcdToken(stakingTokensData.params);
            } else {
                Logger.warn('api-error:', 'there is no result of validators from lcd');
            }
        } catch (e) {
            Logger.warn(`api-error from ${stakingTokensUrl}`, e)
        }
    }
}
