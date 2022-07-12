package errors

import (
	"fmt"
)

const (
	// 请使用 newCommonErr 方法构造该通用错误码
	eC10000 = "11111" // 通用错误码，errMsg 为中英文，RB 可直接弹出提示给用户

	EC40000 = "40000" // 错误的请求参数

	EC40001 = "40001" // 未认证

	EC40002 = "40002" // 时效性验证失败

	EC40003 = "40003" // 记录已存在

	EC40004 = "40004" // 记录未找到

	EC40005 = "40005" // 操作被拒绝

	EC40006 = "40006" // 无效的token

	EC40007 = "40007" // 无效的参数

	EC40010 = "40010" // 手机号未注册

	EC40020 = "40020" // 验证码错误

	EC40030 = "40030" // nft 不存在
	EC40031 = "40031" // nft 不属于当前用户
	EC40032 = "40032" //无效地址
	EC40033 = "40033" //invalid chain
	EC40034 = "40034" //请10分钟之后重试

	EC50001 = "50001"
)

var (
	ErrRecordExist      = fmt.Errorf("record exists while save record")
	ErrUnAuthorization  = fmt.Errorf("unAuthorization")
	ErrNotFound         = fmt.Errorf("record not found")
	ErrNotRegistered    = fmt.Errorf("phone number is not registered")
	ErrVerificationCode = fmt.Errorf("verification code error")
	ErrInvalidToken     = fmt.Errorf("invalid token")
	ErrInvalidParams    = fmt.Errorf("invalid params")
	ErrNftNotExist      = fmt.Errorf("NFT does not exist")
	ErrNftOwner         = fmt.Errorf("NFT does not belong to the current user")
	ErrInvalidSignature = fmt.Errorf("invalid signature")
	ErrTimeliness       = fmt.Errorf("timeliness error")
	ErrInvalidAddress   = fmt.Errorf("invalid address")
	ErrInvalidChain     = fmt.Errorf("invalid chain")
	ErrMetadataRefresh  = fmt.Errorf("We've queued this item for an update! Check back in a minute")
)
