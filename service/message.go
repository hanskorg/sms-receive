package service

import (
	"fmt"
	"sms/dao"
)

func (s *service) Message(numID int, offset, limit int) (list []*dao.Message, total int, err error) {
	list, total, err = s.dao.Messages(numID, offset, limit)
	for _, message := range list {
		message.FromEncoded = fmt.Sprintf("%s***%s", message.From[0:3],message.From[len(message.From)-3:])
	}
	return
}

func (s *service) AddMessage(message *dao.Message) (err error) {
	//TODO 验证来源
	err = s.dao.AddMessage(message)
	return
}
