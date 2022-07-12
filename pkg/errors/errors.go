package errors

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

type RsError struct {
	ErrCode string
	ErrMsg  string
}

func (e *RsError) Error() string {
	return e.ErrMsg
}

func (e *RsError) IsNotNull() bool {
	if e.ErrCode != "" {
		return true
	}
	return false
}

func New(errCode string, errMsg string) *RsError {
	return &RsError{
		ErrCode: errCode,
		ErrMsg:  errMsg,
	}
}

func NewBadRequestErr(err error) *RsError {
	return &RsError{
		ErrCode: EC40000,
		ErrMsg:  err.Error(),
	}
}

func NewInvalidSignatureErr() *RsError {
	return &RsError{
		ErrCode: EC40001,
		ErrMsg:  ErrInvalidSignature.Error(),
	}
}

func NewTimelinessErr() *RsError {
	return &RsError{
		ErrCode: EC40002,
		ErrMsg:  ErrTimeliness.Error(),
	}
}

func NewNotFoundErr() *RsError {
	return &RsError{
		ErrCode: EC40004,
		ErrMsg:  ErrNotFound.Error(),
	}
}

func NewRecordExistErr() *RsError {
	return &RsError{
		ErrCode: EC40003,
		ErrMsg:  ErrRecordExist.Error(),
	}
}

func NewForbiddenErr(err error) *RsError {
	return &RsError{
		ErrCode: EC40005,
		ErrMsg:  err.Error(),
	}
}

func NewUnAuthorizationErr() *RsError {
	return &RsError{
		ErrCode: EC40001,
		ErrMsg:  ErrUnAuthorization.Error(),
	}
}

func NewNotRegisteredErr() *RsError {
	return &RsError{
		ErrCode: EC40010,
		ErrMsg:  ErrNotRegistered.Error(),
	}
}

func NewVerificationCodeErr() *RsError {
	return &RsError{
		ErrCode: EC40020,
		ErrMsg:  ErrVerificationCode.Error(),
	}
}

func NewInvalidTokenErr() *RsError {
	return &RsError{
		ErrCode: EC40006,
		ErrMsg:  ErrInvalidToken.Error(),
	}
}

func NewInvalidParams() *RsError {
	return &RsError{
		ErrCode: EC40007,
		ErrMsg:  ErrInvalidParams.Error(),
	}
}

func NewNftNotExistErr() *RsError {
	return &RsError{
		ErrCode: EC40030,
		ErrMsg:  ErrNftNotExist.Error(),
	}
}

func NewNftOwnerErr() *RsError {
	return &RsError{
		ErrCode: EC40031,
		ErrMsg:  ErrNftOwner.Error(),
	}
}

func NewInternalServerErr(err error) *RsError {
	errMsg := "Internal server err"
	logrus.Errorf(fmt.Sprintf("%s, %s", errMsg, err.Error()))
	return &RsError{
		ErrCode: EC50001,
		ErrMsg:  errMsg,
	}
}

func NewInvalidAddress() *RsError {
	return &RsError{
		ErrCode: EC40032,
		ErrMsg:  ErrInvalidAddress.Error(),
	}
}

func NewErrInvalidChain() *RsError {
	return &RsError{
		ErrCode: EC40033,
		ErrMsg:  ErrInvalidChain.Error(),
	}
}

func NewErrMetadataRefresh() *RsError {
	return &RsError{
		ErrCode: EC40034,
		ErrMsg:  ErrMetadataRefresh.Error(),
	}
}
