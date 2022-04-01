package dao

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/hanskorg/logkit"
	"github.com/jinzhu/gorm"
	"sms/conf"
	"time"
)

type Number struct {
	ID        int        `gorm:"Column:id; primary_key"`
	ZoneID    int        `gorm:"Column:zone"`
	County    string     `gorm:"Column:county"`
	Number    uint64     `gorm:"Column:number"`
	Valid     bool       `gorm:"Column:valid"`
	Free      bool       `gorm:"Column:free"`
	Carrier   int        `gorm:"Column:carrier"`
	CreateAt  time.Time  `gorm:"Column:create_at"`
	ReleaseAt time.Time  `gorm:"Column:release_at"`
	Zone      *conf.Zone `gorm:"-"`
}

func (*Number) TableName() string {
	return "number"
}

type Message struct {
	ID          int       `gorm:"Column:id; primary_key"`
	Number      int       `gorm:"Column:num_id"`
	From        string    `gorm:"type:varchar; Column:from"`
	FromEncoded string    `gorm:"-"`
	IsDel       bool      `gorm:"type:tinyint; Column:is_del"`
	Message     string    `gorm:"Column:message"`
	CreateAt    time.Time `gorm:"Column:created_at"`
	Duration    string    `gorm:"-"`
}

func (*Message) TableName() string {
	return "message"
}

type Dao struct {
	db *gorm.DB
}

func New(c *conf.Database) *Dao {
	var (
		db  *gorm.DB
		err error
	)
	if db, err = gorm.Open("mysql", c.DSN); err == nil {
		if err = db.DB().Ping(); err != nil {
			panic("[bootstrap], db connection error, err: " + err.Error())
		} else {
			logkit.Infof("database connect success %s", c.DSN)
		}
		db.DB().SetMaxOpenConns(c.MaxConnection)
		db.DB().SetMaxIdleConns(c.MaxIdle)
		db.DB().SetConnMaxLifetime(time.Duration(time.Second * 300))
		db.LogMode(conf.C.Debug)
	} else {
		logkit.Error(err.Error())
		panic("[bootstrap], db connection error, err: " + err.Error())
	}
	return &Dao{
		db: db,
	}
}

func (d *Dao) Numbers(cond interface{}, offset, limit int) (numbers []*Number, total int, err error) {
	err = d.db.Where(cond).Offset(offset).Limit(limit).
		Order("id desc").Find(&numbers).Error
	if err != nil {
		return
	}
	err = d.db.Table("number").Where(cond).Or("county = ?", "").Offset(offset).Limit(limit).
		Order("id desc").Count(&total).Error
	return
}

func (d *Dao) Number(id int, zone int, no uint64) (number *Number, err error) {
	number = new(Number)
	if id != 0 {
		err = d.db.Where("id = ?", id).Take(number).Error
	} else {
		err = d.db.Where("zone = ? and number = ? ", zone, no).Take(number).Error
	}
	if err != nil {
		return
	}
	return
}

func (d *Dao) Messages(numId int, offset, limit int) (messages []*Message, total int, err error) {
	err = d.db.Where("is_del = ? AND num_id = ?", false, numId).Offset(offset).Limit(limit).
		Order("id desc").Find(&messages).Error
	if err != nil {
		return
	}
	err = d.db.Table("message").Where("is_del = ? AND num_id = ?", false, numId).Offset(offset).Limit(limit).
		Order("id desc").Count(&total).Error
	return
}

func (d *Dao) AddMessage(message *Message) (err error) {
	err = d.db.Save(message).Error
	if err != nil {
		return
	}
	return
}
