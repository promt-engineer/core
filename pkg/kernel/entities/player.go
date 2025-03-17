package entities

import "github.com/gin-gonic/gin"

type PlayerMetaData struct {
	IP        string `validate:"required,ip_addr"`
	UserAgent string `validate:"required"`
	Host      string `validate:"required,url"`
	Request   []byte
}

func (pmd *PlayerMetaData) CopyAndSetRequest(req []byte) *PlayerMetaData {
	return &PlayerMetaData{
		IP:        pmd.IP,
		UserAgent: pmd.UserAgent,
		Host:      pmd.Host,
		Request:   req,
	}
}

func NewPlayerMetaDataFromCtx(c *gin.Context, request []byte) *PlayerMetaData {
	return &PlayerMetaData{
		IP:        c.ClientIP(),
		UserAgent: c.GetHeader("User-Agent"),
		Host:      c.Request.Header.Get("Origin"),
		Request:   request,
	}
}
