package service

import (
	"sms/conf"
	"sms/dao"
)

func (s *service) Numbers(zone int, county string, start, offset int) (list []*dao.Number, total int, err error) {
	cond := &dao.Number{
		ZoneID: zone,
		Free:   true,
	}
	if county != "" {
		cond.County = county
	}
	list, total, err = s.dao.Numbers(cond, start, offset)
	for _, val := range list {
		for _, zoneInfo := range conf.C.Zone {
			if (zoneInfo.ID == val.ZoneID && val.County == "") || (val.County != "" && zoneInfo.Short == val.County) {
				val.Zone = zoneInfo
			}
		}
	}
	return
}

func (s *service) Number(id int, zone int, no uint64) (number *dao.Number, err error) {
	number, err = s.dao.Number(id, zone, no)
	for _, zoneInfo := range conf.C.Zone {
		if (zoneInfo.ID == number.ZoneID && number.County == "") || (number.County != "" && zoneInfo.Short == number.County) {
			number.Zone = zoneInfo
		}
	}
	return
}
