package types

import (
	"github.com/bianjieai/iobscan-explorer-backend/internal/app/constant"
	"github.com/bianjieai/iobscan-explorer-backend/pkg/errors"
	"github.com/sirupsen/logrus"
)

type BaseResponse struct {
	Status  string      `json:"status"`
	ErrCode string      `json:"err_code"`
	ErrMsg  string      `json:"err_msg"`
	Data    interface{} `json:"data"`
}

type BasePageResponse struct {
	Status   string      `json:"status"`
	ErrCode  string      `json:"err_code"`
	ErrMsg   string      `json:"err_msg"`
	Data     interface{} `json:"data"`
	Total    int         `json:"total"`
	PageNum  int         `json:"page_num" `
	PageSize int         `json:"page_size"`
}

type Page struct {
	PageNum  int `json:"page_num" form:"page_num"`
	PageSize int `json:"page_size" form:"page_size"`
}

func BuildResponse(data interface{}) BaseResponse {
	return BaseResponse{
		Status: constant.Success,
		Data:   data,
	}
}
func BuildPageResponse(data interface{}, total int, pageNum int, pageSize int) BasePageResponse {
	return BasePageResponse{
		Status:   constant.Success,
		Data:     data,
		Total:    total,
		PageNum:  pageNum,
		PageSize: pageSize,
	}
}

func BuildErrResponse(err *errors.RsError) BaseResponse {
	logrus.Errorf("server response failed,err_code:%s,err_msg:%s", err.ErrCode, err.ErrMsg)
	return BaseResponse{
		Status:  constant.Fail,
		ErrCode: err.ErrCode,
		ErrMsg:  err.ErrMsg,
	}
}

func ParseParamPage(pageNum int, pageSize int) (skip int, limit int) {
	if pageNum == 0 && pageSize == 0 {
		pageSize = 10
	}
	return (pageNum - 1) * pageSize, pageSize
}
