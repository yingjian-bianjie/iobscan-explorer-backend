package enum

type (
	TxType   string
	TxStatus int
	TaskEnum string
)

const (
	DefineService TxType   = "define_service"
	Success       TxStatus = 1
	Failed        TxStatus = 0
)

const (
	Denom                     TaskEnum = "ex_sync_denom"
	Nft                       TaskEnum = "ex_sync_nft"
	Identity                  TaskEnum = "sync_identity"
	TxServiceName             TaskEnum = "sync_tx_service_name"
	Validaters                TaskEnum = "sync_validators"
	StakingSyncValidatorsInfo TaskEnum = "staking_sync_validators_info"
	Tokens                    TaskEnum = "tokens"
)
