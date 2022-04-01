package service

import (
	"github.com/kevinburke/twilio-go"
	"sms/conf"
	"sms/dao"
)

var (
	srv *service
)

type service struct {
	twilio *twilio.Client
	dao    *dao.Dao
}

type Service interface {
	Number(id int, zone int, no uint64) (number *dao.Number, err error)
	Numbers(zone int, county string, start, offset int) (list []*dao.Number, total int, err error)
	Message(numID int, offset, limit int) (list []*dao.Message, total int, err error)
	AddMessage(*dao.Message) error
}

func New() Service {
	srv = &service{
		twilio: twilio.NewClient(conf.C.Twilio.SID, conf.C.Twilio.Token, nil),
		dao:    dao.New(conf.C.Database),
	}
	return srv
}
