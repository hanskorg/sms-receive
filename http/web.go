package http

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hanskorg/logkit"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"net/http"
	"sms/conf"
	"sms/dao"
	"time"
)

func (s *Server) message(ctx *gin.Context) {
	var (
		messages []*dao.Message
		number   *dao.Number
		total    int
		err      error
		params   = struct {
			Page     int `uri:"page,default=0"`
			Size     int `uri:"size,default=20"`
			NumberID int `uri:"nid"`
		}{}
	)
	err = ctx.ShouldBindUri(&params)
	if err != nil {
		logkit.Errorf("params err, %s", err.Error())
		ctx.HTML(http.StatusBadRequest, "info.html", gin.H{
			"err": "Bad Error",
		})
		return
	}
	messages, total, err = s.svc.Message(params.NumberID, params.Page*params.Size, params.Size)
	number, err = s.svc.Number(params.NumberID, 0, 0)
	if err != nil {
		ctx.HTML(http.StatusOK, "info.html", gin.H{
			"err": "Number Error",
		})
		return
	}
	if err != nil {
		logkit.Errorf("message query fail, %s", err.Error())
		ctx.AbortWithError(http.StatusBadRequest, fmt.Errorf("errors"))
		return
	}
	logkit.Infof("query message total: %d", total)
	for _, message := range messages {
		dur := time.Now().Sub(message.CreateAt).Seconds()
		if dur < 60 {
			message.Duration = fmt.Sprintf("%d secs ago", int(dur))
		} else if dur < 3600 {
			message.Duration = fmt.Sprintf("%d mins %d secs ago", int(dur)/60, int(dur)%60)
		} else if dur < 3600*24 {
			message.Duration = fmt.Sprintf("%d hours %d mins ago", int(dur)/3600, int(dur)%3600/60)
		} else {
			message.Duration = fmt.Sprintf("%d days %d hours ago", int(dur)/86400, int(dur) % 86400 / 3600 )

		}
	}
	global,_ := ctx.Get("global")
	values := global.(gin.H)
	values["number"] = number
	values["zones"] = conf.C.Zone
	values["zoneID"] = number.Zone.ID
	values["cname"] = number.Zone.County
	values["err"] = nil
	ctx.HTML(http.StatusOK, "info.html", values)
}

func (s *Server) number(ctx *gin.Context) {
	var (
		zone struct {
			ID     int    `uri:"zone,default=0"`
			County string `uri:"c"`
			Page   int    `uri:"page,default=0"`
		}
		numbers []*dao.Number
		err     error
		total   int
	)
	err = ctx.ShouldBindUri(&zone)
	if err != nil {
		zone.ID = 1
	}
	numbers, total, err = s.svc.Numbers(zone.ID, zone.County, 0, 20)
	if err != nil {
		logkit.Errorf("number query fail, %s", err.Error())
		ctx.AbortWithError(http.StatusBadRequest, fmt.Errorf("errors"))
		return
	}
	logkit.Infof("query message total: %d", total)
	global,_ := ctx.Get("global")
	params := global.(gin.H)
	params["numbers"] = numbers
	params["zones"] = conf.C.Zone
	params["zoneID"] = zone.ID
	params["cname"] = zone.County
	ctx.HTML(http.StatusOK, "index.html" , params)
}
func (s *Server) globalVal(ctx *gin.Context)  {
	localizer,_ := ctx.Get("localizer")
	ctx.Set("global", gin.H{
		"title": localizer.(*i18n.Localizer).MustLocalize(&i18n.LocalizeConfig{
			MessageID:      "Title",
			DefaultMessage: &i18n.Message{ID: "Title"},
		}),
		"keywords": localizer.(*i18n.Localizer).MustLocalize(&i18n.LocalizeConfig{
			MessageID:      "Keywords",
		}),
		"description": localizer.(*i18n.Localizer).MustLocalize(&i18n.LocalizeConfig{
			MessageID:      "Description",
		}),
		"country": localizer.(*i18n.Localizer).MustLocalize(&i18n.LocalizeConfig{
			MessageID:      "Country",
		}),
	})
	ctx.Next()
}

func  (s *Server) langMatch(ctx *gin.Context) {

	accept := ctx.Request.Header.Get("Accept-Language")
	tags,_, _ := language.ParseAcceptLanguage(accept)
	tag, _, _ :=  matcher.Match(tags...)
	switch tag.Parent() {
	case language.Chinese, language.SimplifiedChinese, language.TraditionalChinese:
		ctx.Set("region", "cn")
	case language.Und, language.English, language.Armenian:
		ctx.Set("region", "en")
	default:
		ctx.Set("region", "en")
	}
	if localizer, ok := s.localizes[tag.Parent().String()]; ok {
		ctx.Set("localizer", localizer)
	}else{
		localizer = i18n.NewLocalizer(s.bundle, tag.Parent().String() , accept)
		s.locker.Lock()
		defer s.locker.Unlock()
		s.localizes[tag.Parent().String()] = localizer
		ctx.Set("localizer", localizer )
	}
	ctx.Next()
}
