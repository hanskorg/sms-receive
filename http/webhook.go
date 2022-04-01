package http

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hanskorg/logkit"
	"github.com/kevinburke/twilio-go"
	"github.com/ttacon/libphonenumber"
	"net/http"
	"sms/conf"
	"sms/dao"
	"time"
)

const (
	contentType = "application/xml"
)

type TwiloRequestBody struct {
	To          string `form:"To" binding:"required"`
	From        string `form:"From" binding:"required"`
	Body        string `form:"Body" binding:"required"`
	MessageSid  string `form:"MessageSid" binding:"required"`
	Carrier     string
	ToCountry   string `form:"ToCountry" binding:"required"`
	FromCountry string `form:"FromCountry"`
}

//{
//"api-key": "abcd1234",
//"msisdn": "447700900001",
//"to": "447700900000",
//"messageId": "0A0000000123ABCD1",
//"text": "Hello world",
//"type": "text",
//"keyword": "HELLO",
//"message-timestamp": "2020-01-01 12:00:00 +0000",
//"timestamp": "1578787200",
//"nonce": "aaaaaaaa-bbbb-cccc-dddd-0123456789ab",
//"concat": "true",
//"concat-ref": "1",
//"concat-total": "3",
//"concat-part": "2"
//}
type NexmoSMS struct {
	ApiKey      string `json:"api-key" uri:"api-key"`
	From        string `json:"msisdn" uri:"msisdn"`
	To          string `json:"to" uri:"to"`
	MessageID   string `json:"messageId" uri:"messageId"`
	Body        string `json:"text" uri:"text"`
	Type        string `json:"type" uri:"type"` //unicode text binary
	Keyword     string `json:"keyword" uri:"keyword"`
	Time        string `json:"message-timestamp" uri:"message-timestamp"`
	Timestamp   string `json:"timestamp" uri:"timestamp"`
	Nonce       string `json:"nonce" uri:"nonce"`
	Concat      string `json:"concat" uri:"concat"`
	ConcatTotal string `json:"concat-total" uri:"concat-total"`
	ConcatRef   string `json:"concat-ref" uri:"concat-ref"`
	ConcatPart  string `json:"concat-part" uri:"concat-part"`
}

func (s *Server) status(ctx *gin.Context) {
	logkit.Infof("status hook: %+v", ctx.Request.Header)
	ctx.JSON(http.StatusOK, "")
}
func (s *Server) twilo(ctx *gin.Context) {
	var (
		params   = &TwiloRequestBody{}
		toNumber *libphonenumber.PhoneNumber
		number   *dao.Number
		err      error
	)
	resp := `<?xml version="1.0" encoding="UTF-8"?>
<Response>
    <Message action="http://%s/hook/twilio/status"><Body>%s</Body></Message>
</Response>`
	err = ctx.ShouldBind(params)
	if err != nil {
		logkit.Infof("incoming message parse fail, to number error, %s", err.Error())

		ctx.Data(http.StatusBadRequest, contentType, []byte(fmt.Sprintf(resp, conf.C.Twilio.Hook, "INVAILD")))
		return
	}
	//验证号码
	toNumber, err = libphonenumber.Parse(params.To, params.ToCountry)
	if err != nil {
		ctx.Data(http.StatusBadRequest, contentType, []byte(fmt.Sprintf(resp, conf.C.Twilio.Hook, "INVAILD")))
		logkit.Infof("incoming message parse fail, to number error, %s", err.Error())
		return
	}
	number, err = s.svc.Number(0, int(*toNumber.CountryCode), *toNumber.NationalNumber)
	if err != nil {
		logkit.Infof("incoming message query fail, %s", err.Error())
		ctx.Data(http.StatusInternalServerError, contentType, []byte(fmt.Sprintf(resp, conf.C.Twilio.Hook, "FAIL")))
		return
	}
	err = s.svc.AddMessage(&dao.Message{
		Number:   number.ID,
		From:     params.From,
		IsDel:    false,
		Message:  params.Body,
		CreateAt: time.Now(),
	})
	ctx.Data(http.StatusOK, contentType, []byte(fmt.Sprintf(resp, conf.C.Twilio.Hook, "SUCCESS")))
}

func (s *Server) nexmo(ctx *gin.Context) {
	var (
		params   = &NexmoSMS{}
		toNumber *libphonenumber.PhoneNumber
		number   *dao.Number
		err      error
	)
	err = ctx.ShouldBindJSON(params)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, fmt.Errorf("params error"))
		return
	}
	toNumber, err = libphonenumber.Parse(params.To, "US")
	//验证号码
	toNumber, err = libphonenumber.Parse(params.To, "US")
	if err != nil {
		logkit.Infof("incoming message parse fail, to number error, %s", err.Error())
		return
	}
	number, err = s.svc.Number(0, int(*toNumber.CountryCode), *toNumber.NationalNumber)
	if err != nil {
		logkit.Infof("incoming message query fail, %s", err.Error())
		return
	}
	err = s.svc.AddMessage(&dao.Message{
		Number:   number.ID,
		From:     params.From,
		IsDel:    false,
		Message:  ctx.GetString("Body"),
		CreateAt: time.Now(),
	})
	if err != nil {
		ctx.Status(http.StatusNoContent)
	}
}
func twiloValid(ctx *gin.Context) {
	if err := twilio.ValidateIncomingRequest(conf.C.Twilio.Hook, conf.C.Twilio.Token, ctx.Request); err != nil {
		ctx.AbortWithStatus(http.StatusUnauthorized)
	}
	ctx.Next()
}

func nexmoValid(ctx *gin.Context) {
	ctx.Next()
}
