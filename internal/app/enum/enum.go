package enum

type (
	TxType   string
	TxStatus int
)

const (
	DefineService TxType   = "define_service"
	Success       TxStatus = 1
	Failed        TxStatus = 0
)
