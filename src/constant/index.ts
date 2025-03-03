import { cfg } from '../config/config';

export enum TaskEnum {
    denom = 'ex_sync_denom',
    nft = 'ex_sync_nft',
    txServiceName = "sync_tx_service_name",
    validators = 'sync_validators',
    identity = 'sync_identity',
    stakingSyncValidatorsInfo = 'staking_sync_validators_info',
    stakingSyncValidatorsMoreInfo = 'staking_sync_validators_more_info',
    stakingSyncParameters = 'staking_sync_parameters',
    tokens = 'tokens',
    proposal = 'ex_sync_proposal',
    statistics = 'ex_statistics',
    account = 'ex_sync_account',
    accountInfo =  'ex_sync_account_info'
}

export const StatisticsNames = [
    'tx_all',
    'service_all',
    'validator_all',
    'validator_active',
    'identity_all',
    'nft_all',
    'denom_all',
    'bonded_tokens',
    'total_supply',
    'community_pool',
    'accounts_all',
];

export const DefaultPaging = {
    pageNum: 1,
    pageSize: 10,
};

export enum ENV {
    development = 'development',
    production = 'production',
};

export enum DDCType {
    contractTag= 'ddc_',
    dataDdc = 'DDC',
    ddc721 = 'DDC721',
    ddc1155 = 'DDC1155',
}
export const ContractType = {
    'ddc_other':-1,
    'ddc_721':1,
    'ddc_1155':2,
}

//txs/addresses/statistic (datangchain-explorer)
export enum TxCntQueryCond {
    nftQueryCnt = 1,
    sendQueryCnt,
    multisendQueryCnt
}
export enum TxType {
    // service
    define_service = 'define_service',
    bind_service = 'bind_service',
    call_service = 'call_service',
    respond_service = 'respond_service',
    update_service_binding = 'update_service_binding',
    disable_service_binding = 'disable_service_binding',
    enable_service_binding = 'enable_service_binding',
    refund_service_deposit = 'refund_service_deposit',
    pause_request_context = 'pause_request_context',
    start_request_context = 'start_request_context',
    kill_request_context = 'kill_request_context',
    update_request_context = 'update_request_context',
    service_set_withdraw_address = 'service/set_withdraw_address',
    withdraw_earned_fees = 'withdraw_earned_fees',
    // nft
    burn_nft = 'burn_nft',
    transfer_nft = 'transfer_nft',
    edit_nft = 'edit_nft',
    issue_denom = 'issue_denom',
    mint_nft = 'mint_nft',
    transfer_denom = 'transfer_denom',

    // tibc
    tibc_nft_transfer = 'tibc_nft_transfer',
    tibc_recv_packet = 'tibc_recv_packet',
    tibc_acknowledge_packet = 'tibc_acknowledge_packet',

    // Asset
    issue_token = 'issue_token',
    edit_token = 'edit_token',
    mint_token = 'mint_token',
    transfer_token_owner = 'transfer_token_owner',
    burn_token= 'burn_token',
    //Transfer
    send = 'send',
    multisend = 'multisend',
    //Crisis
    verify_invariant = 'verify_invariant',
    //Evidence
    submit_evidence = 'submit_evidence',
    //Staking
    begin_unbonding = 'begin_unbonding',
    edit_validator = 'edit_validator',
    create_validator = 'create_validator',
    delegate = 'delegate',
    begin_redelegate = 'begin_redelegate',
    // Slashing
    unjail = 'unjail',
    // Distribution
    set_withdraw_address = 'set_withdraw_address',
    withdraw_delegator_reward = 'withdraw_delegator_reward',
    withdraw_validator_commission = 'withdraw_validator_commission',
    fund_community_pool = 'fund_community_pool',
    // Gov
    deposit = 'deposit',
    vote = 'vote',
    submit_proposal = 'submit_proposal',
    // Coinswap
    add_liquidity = 'add_liquidity',
    remove_liquidity = 'remove_liquidity',
    swap_order = 'swap_order',
    // Htlc
    create_htlc = 'create_htlc',
    claim_htlc = 'claim_htlc',
    refund_htlc = 'refund_htlc',
    // Evm
    ethereum_tx = 'ethereum_tx',
    // Guardian
    add_profiler = 'add_profiler',
    delete_profiler = 'delete_profiler',
    add_trustee = 'add_trustee',
    delete_trustee = 'delete_trustee',
    add_super = 'add_super',
    // Oracle
    create_feed = 'create_feed',
    start_feed = 'start_feed',
    pause_feed = 'pause_feed',
    edit_feed = 'edit_feed',
    // IBC
    recv_packet = 'recv_packet',
    create_client = 'create_client',
    update_client = 'update_client',
    // Identity
    create_identity = 'create_identity',
    update_identity = 'update_identity',
    // Record
    create_record = 'create_record',
    // Random
    request_rand = 'request_rand',
    channel_open_init = 'channel_open_init',
    channel_open_confirm = 'channel_open_confirm',
    channel_open_try = 'channel_open_try',
    connection_open_init = 'connection_open_init',
    connection_open_confirm = 'connection_open_confirm',
    connection_open_try = 'connection_open_try',
    connection_open_ack = 'connection_open_ack',
}

export enum TxStatus {
    SUCCESS = 1,
    FAILED = 0,
}

export const IdentityLimitSize = 1000

export enum LoggerLevel {
    ALL = 'ALL',
    TRACE = 'TRACE',
    DEBUG = 'DEBUG',
    INFO = 'INFO',
    WARN = 'WARN',
    ERROR = 'ERROR',
    FATAL = 'FATAL',
    MARK = 'MARK',
    OFF = 'OFF',
}

export const PubKeyAlgorithm = {
    0: 'UnknownPubKeyAlgorithm',
    1: 'RSA',
    2: 'DSA',
    3: 'ECDSA',
    4: 'ED25519',
    5: 'SM2',
}

export enum currentChain  {
    cosmos ='cosmos',
    iris ='iris',
    binance = 'binance'
}
export const deFaultGasPirce = 1
export const signedBlocksWindow = 'signed_blocks_window'
export const hubDefaultEmptyValue = '[do-not-modify]'
export const moduleSlashing = 'slashing'
export const moduleStaking = 'staking'
export const moduleStakingBondDenom = 'bond_denom'
export const moduleGov = 'gov'
export const moduleGovDeposit = 'min_deposit'
export const defaultEvmTxReceiptErrlog = 'failed to execute message'



export const ValidatorStatus = {
    'Unbonded': 1,
    'Unbonding': 2,
    'bonded': 3,
}

let addressPrefix,validatorStatusStr;
switch (cfg.currentChain) {
    case currentChain.iris:
        // validatorStatusStr = {
        //     'unbonded': 'unbonded',
        //     'unbonding': 'unbonding',
        //     'bonded': 'bonded'
        // };
        validatorStatusStr = {
            'unbonded': 'BOND_STATUS_UNBONDED', // 1 关押状态
            'unbonding': 'BOND_STATUS_UNBONDING', // 2  候选人状态
            'bonded': 'BOND_STATUS_BONDED' // 3 活跃的验证人状态
        };
        addressPrefix = {
            iaa: 'iaa',
            iva: 'iva',
            ica: 'ica',
            icp: 'icp'
        }
        break;
    case currentChain.cosmos:
        validatorStatusStr = {
            'unbonded': 'BOND_STATUS_UNBONDED',
            'unbonding': 'BOND_STATUS_UNBONDING',
            'bonded': 'BOND_STATUS_BONDED'
        };
        addressPrefix = {
            iaa: 'cosmos',
            iva: 'cosmosvaloper',
            ica: 'cosmosvalcons',
            icp: 'cosmosvalconspub'
        }
        break;
    default:
        break;
}
export {
    validatorStatusStr,addressPrefix
}

export const validatorStatusFromLcd = {
    'BOND_STATUS_UNBONDED': 1,
    'BOND_STATUS_UNBONDING': 2,
    'BOND_STATUS_BONDED': 3
}

export const ValidatorNumberStatus = {
    1: 'candidate',
    2: 'candidate',
    3: 'active',
}
export const activeValidatorLabel = 'active'
export const candidateValidatorLabel = 'candidate'
export const jailedValidatorLabel = 'jailed'

export const INCREASE_HEIGHT = Number(cfg.taskCfg.increaseHeight);
export const MAX_OPERATE_TX_COUNT = Number(cfg.taskCfg.maxOperateTxCount);
export const MAX_DENOM_TX_COUNT = Number(1000);

export const NFT_INFO_DO_NOT_MODIFY = '[do-not-modify]';

export const correlationStr = {
    '200': 'block',
    '201': 'txCount',
    '202': 'validatorCount',
    '203': 'avgBlockTime',
    '204': 'assetCount',
    '205': 'denomCount',
    '206': 'serviceCount',
    '207': 'identityCount',
    '208': 'validatorNumCount',
    '209': 'bondedTokensInformation',
    '210': 'communityPoolInformation',
    '211': 'accountsCount',
}

export const proposalStatus = {
    PROPOSAL_STATUS_DEPOSIT_PERIOD: 'DepositPeriod',
    PROPOSAL_STATUS_VOTING_PERIOD: 'VotingPeriod',
    PROPOSAL_STATUS_PASSED: 'Passed',
    PROPOSAL_STATUS_REJECTED: 'Rejected'
}

export const govParams = {
    min_deposit: 'min_deposit',
    quorum: 'quorum',
    threshold: 'threshold',
    veto_threshold:'veto_threshold'
}

export const voteOptions = {
    1: 'yes',
    2: 'abstain',
    3: 'no',
    4: 'no_with_veto'
}
export const proposal = 'Proposal'

export const queryVoteOptionCount = {
    yes: 1,
    abstain: 2,
    no: 3,
    no_with_veto: 4,
}

export const addressAccount = 'xxx'

export const SRC_PROTOCOL = {
    NATIVE:'native',
    HTLT:'htlt',
    IBC:'ibc',
    SWAP:'swap',
    PEG:'peg',
}
