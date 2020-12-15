import dotenv from 'dotenv';
dotenv.config();

const {
    LCD_ADDR,
    DB_USER,
    DB_PASSWD,
    DB_ADDR,
    ICONURI,
    DB_DATABASE,
    NODE_ENV,
    DENOM_EXECUTE_TIME,
    NFT_EXECUTE_TIME,
    TX_SERVICE_NAME_EXECUTE_TIME,
    FAULT_TOLERANCE_EXECUTE_TIME,
    IDENTITY_EXECUTE_TIME,
    SYNC_TX_SERVICE_NAME_SIZE,
    HEARTBEAT_RATE,
    VALIDATORS_EXECTUTE_TIME,
    DisableLog,
    STAKING_VALIDATORS_TIME,
    STAKING_PARAMETERS,
    TOKENS,
    INCREASE_HEIGHT,
    MAX_OPERATE_TX_COUNT,
    CURRENT_CHAIN,
    MAIN_TOKEN,
    DELETE_CRON_JOBS
} = process.env;
export const cfg = {
    env: NODE_ENV,
    disableLog:Boolean(DisableLog=='true'),
    dbCfg: {
        user: DB_USER,
        psd: DB_PASSWD,
        dbAddr: DB_ADDR,
        dbName: DB_DATABASE,
    },
    serverCfg:{
        lcdAddr:LCD_ADDR,
        iconUri:ICONURI || 'https://keybase.io/_/api/1.0/user/lookup.json'
    },
    taskCfg:{
        interval:{
            heartbeatRate:Number(HEARTBEAT_RATE || 10000),
        },
        executeTime:{
            denom:DENOM_EXECUTE_TIME || '1 * * * * *',
            nft:NFT_EXECUTE_TIME || '21 * * * * *',
            txServiceName:TX_SERVICE_NAME_EXECUTE_TIME || '30 * * * * *',
            faultTolerance:FAULT_TOLERANCE_EXECUTE_TIME || '41 * * * * *',
            validators:VALIDATORS_EXECTUTE_TIME || '1 * * * * *',
            identity: IDENTITY_EXECUTE_TIME || '1 * * * * *',
            stakingValidators: STAKING_VALIDATORS_TIME || '15 * * * * *',
            stakingParameters: STAKING_PARAMETERS || '10 * * * * *',
            Tokens: TOKENS || '5 * * * * *',
        },
        syncTxServiceNameSize: Number(SYNC_TX_SERVICE_NAME_SIZE) || 200,
        increaseHeight: INCREASE_HEIGHT || 1000,
        maxOperateTxCount: MAX_OPERATE_TX_COUNT || 100,
        DELETE_CRON_JOBS: DELETE_CRON_JOBS ? JSON.parse(DELETE_CRON_JOBS) : []
    },
    currentChain: CURRENT_CHAIN || 'IRIS',
    MAIN_TOKEN: MAIN_TOKEN ? JSON.parse(MAIN_TOKEN) : {"min_unit":"umuon","scale":"6","symbol":"muon"}
};

