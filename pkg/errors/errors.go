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

func NewInvalidParams() *RsError {
	return &RsError{
		ErrCode: EC40007,
		ErrMsg:  ErrInvalidParams.Error(),
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
