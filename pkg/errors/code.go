package errors

import (
	"fmt"
)

const (
	// 请使用 newCommonErr 方法构造该通用错误码
	eC10000 = "11111" // 通用错误码，errMsg 为中英文，RB 可直接弹出提示给用户

	EC40000 = "40000" // 错误的请求参数

	EC40007 = "40007" // 无效的参数

	EC50001 = "50001"
)

var (
	ErrInvalidParams = fmt.Errorf("invalid params")
)
