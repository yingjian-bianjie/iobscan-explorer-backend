package contracts

//type MethodData struct {
//	Method    ABI.Method
//	Contracts string
//}

const (
	EvmDdcType        = "DDC"
	ContractDDC721    = "DDC721"
	ContractDDC1155   = "DDC1155"
	ContractAuthority = "Authority"
	ContractCharge    = "Charge"
)

const (
	_ = iota
	ContractDDC721Int
	ContractDDC1155Int
	ContractAuthorityInt
	ContractChargeInt
)

const (
	_ = iota
	MintDdc
	TransferDdc
	EditDdc
	SetURI
	//FreezeDdc
	//UnFreezeDdc
	BurnDdc
	MintBatchDdc
	SafeMintDdc
)

var (
	DdcMethod = map[string]int{
		"mint":                  MintDdc,
		"safeMint":              SafeMintDdc,
		"safeMintBatch":         MintBatchDdc,
		"safeTransferFrom":      TransferDdc,
		"transferFrom":          TransferDdc,
		"safeBatchTransferFrom": TransferDdc,
		"setNameAndSymbol":      EditDdc,
		"setURI":                SetURI,
		//"freeze":                FreezeDdc,
		//"unFreeze":              UnFreezeDdc,
		"burn":      BurnDdc,
		"burnBatch": BurnDdc,
	}
	DdcType = map[string]int{
		ContractDDC721:    ContractDDC721Int,
		ContractDDC1155:   ContractDDC1155Int,
		ContractAuthority: ContractAuthorityInt,
		ContractCharge:    ContractChargeInt,
	}
	DdcTypeName = map[int]string{
		ContractDDC721Int:    ContractDDC721,
		ContractDDC1155Int:   ContractDDC1155,
		ContractAuthorityInt: ContractAuthority,
		ContractChargeInt:    ContractCharge,
	}
)

type LegacyTx struct {
	// nonce corresponds to the account nonce (transaction sequence).
	Nonce uint64 `json:"nonce,omitempty"`
	// gas price defines the value for each gas unit
	GasPrice string `json:"gas_price,omitempty"`
	// gas defines the gas limit defined for the transaction.
	GasLimit uint64 `json:"gas,omitempty"`
	// hex formatted address of the recipient
	To string `json:"to,omitempty"`
	// value defines the unsigned integer value of the transaction amount.
	Amount string ` json:"value,omitempty"`
	// input defines the data payload bytes of the transaction.
	Data []byte `json:"data,omitempty"`
	// v defines the signature value
	V []byte `json:"v,omitempty"`
	// r defines the signature value
	R []byte `json:"r,omitempty"`
	// s define the signature value
	S []byte `json:"s,omitempty"`
}

// MsgEthereumTx encapsulates an Ethereum transaction as an SDK message.
type DocMsgEthereumTx struct {
	Data  string  `bson:"data" json:"data"`
	Size_ float64 `bson:"size" json:"size_"`
	Hash  string  `bson:"hash" json:"hash"`
	From  string  `bson:"from" json:"from"`

	EvmType      string
	ContractAddr string
	DdcType      string
	DdcId        int64
	Method       string
	Inputs       []string
	Outputs      []string
	TxTime       int64
	TxHeight     int64
	Signer       string
}
