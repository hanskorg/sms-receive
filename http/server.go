package http

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/gin-gonic/gin"
	"github.com/hanskorg/logkit"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"net/http"
	"sms/conf"
	"sms/service"
	"sync"
)

//Server http server 定义
type Server struct {
	Listen    string
	debug     bool
	Router    *gin.Engine
	Server    *http.Server
	svc       service.Service
	bundle    *i18n.Bundle
	localizes map[string]*i18n.Localizer
	locker   sync.RWMutex
}
var (
	matcher = language.NewMatcher([]language.Tag{
		language.English,
		language.SimplifiedChinese,
		language.Chinese,
	})
)
//New 开启http server
func New(svc service.Service) *Server {
	s := &Server{
		Listen: conf.C.Bind,
		svc:    svc,
		localizes: map[string]*i18n.Localizer{},
	}
	if conf.Debug {
		gin.SetMode(gin.DebugMode)
		s.Router = gin.New()
		s.Router.Use(gin.Logger())

	} else {
		gin.SetMode(gin.ReleaseMode)
		s.Router = gin.New()
	}
	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	bundle.MustLoadMessageFile("i18n/active.en.toml")
	bundle.MustLoadMessageFile("i18n/active.zh-CN.toml")
	s.bundle = bundle

	s.Router.LoadHTMLGlob(fmt.Sprintf("%shtml/template/*", conf.C.TemplateDir))
	s.Router.Static("/assets", fmt.Sprintf("%shtml/assets", conf.C.TemplateDir))

	gin.DefaultErrorWriter = logkit.NewLogWriter(logkit.LevelError)
	gin.DefaultWriter = logkit.NewLogWriter(logkit.LevelInfo)

	t := s.Router.Group("/hook/twilio")
	{
		t.POST("/status", s.status)
		t.POST("/income", s.twilo).Use(twiloValid)
	}
	s.Router.POST("/delivery/nexmo", func(context *gin.Context) {
		context.Status(http.StatusNoContent)
	})
	web := s.Router.Group("/")
	{
		web.Use(s.langMatch)
		web.Use(s.globalVal)
		web.GET("/", s.number)
		web.GET("/info/:nid", s.message)
		web.GET("/zone/:zone", s.number)
		web.GET("/zone/:zone/:c", s.number)
	}
	go func() {
		s.Server = &http.Server{
			Addr:    s.Listen,
			Handler: s.Router,
		}
		s.Router.Use(gin.Recovery())
		if err := s.Server.ListenAndServe(); err != nil {
			logkit.Infof("http server fail, %s\n", err.Error())
		}
	}()
	logkit.Infof("http server started, %s\n", s.Listen)
	return s
}
