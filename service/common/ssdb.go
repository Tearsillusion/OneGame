package common

import (
	"errors"
	"log"

	"github.com/seefan/gossdb"
)

var (
	ssdb_pool *gossdb.Connectors
)

//=============================================================
//ssdb-map基类
type SsdbHash struct {
	value   int64 // 属于谁的数据
	mainkey string
	ssdb    *gossdb.Client // ssdb驱动
}

func (self *SsdbHash) InitHash(value int64, mainkey string, client *gossdb.Client) {
	self.value = value
	self.mainkey = mainkey
	self.ssdb = client
}

func (self *SsdbHash) SetValue(key string, value interface{}) error {
	if self.ssdb == nil {
		return errors.New("ssdb == nil")
	}
	return self.ssdb.Hset(self.mainkey, key, value)
}

func (self *SsdbHash) GetValue(key string) (value gossdb.Value, err error) {
	if self.ssdb == nil {
		return "", errors.New("ssdb == nil")
	}
	return self.ssdb.Hget(self.mainkey, key)
}

func (self *SsdbHash) IncrValue(key string, value int64) (val int64, err error) {
	if self.ssdb == nil {
		return 0, errors.New("ssdb == nil")
	}
	return self.ssdb.Hincr(self.mainkey, key, value)
}

func (self *SsdbHash) Size() (value int64, err error) {
	if self.ssdb == nil {
		return 0, errors.New("ssdb == nil")
	}
	return self.ssdb.Hsize(self.mainkey)
}

func (self *SsdbHash) Clear() error {
	if self.ssdb == nil {
		return errors.New("ssdb == nil")
	}
	return self.ssdb.Hclear(self.mainkey)
}

func (self *SsdbHash) Del(key string) error {
	if self.ssdb == nil {
		return errors.New("ssdb == nil")
	}
	return self.ssdb.Hdel(self.mainkey, key)
}

func (self *SsdbHash) Exist(key string) (b bool, err error) {
	if self.ssdb == nil {
		return false, errors.New("ssdb == nil")
	}
	return self.ssdb.Hexists(self.mainkey, key)
}

//=============================================================
//ssdb-list基类
type SsdbQueue struct {
	userid  int64 // 属于谁的数据
	mainkey string
	ssdb    *gossdb.Client // ssdb驱动
}

func (self *SsdbQueue) InitQueue(userid int64, mainkey string, client *gossdb.Client) {
	self.userid = userid
	self.mainkey = mainkey
	self.ssdb = client
}

func (self *SsdbQueue) QSize() (size int64, err error) {
	if self.ssdb == nil {
		return 0, errors.New("ssdb == nil")
	}
	return self.ssdb.Qsize(self.mainkey)
}

func (self *SsdbQueue) QClear() error {
	if self.ssdb == nil {
		return errors.New("ssdb == nil")
	}
	return self.ssdb.Qclear(self.mainkey)
}

func (self *SsdbQueue) QPushBack(value ...interface{}) (size int64, err error) {
	if self.ssdb == nil {
		return 0, errors.New("ssdb == nil")
	}
	return self.ssdb.Qpush_back(self.mainkey, value...)
}

func (self *SsdbQueue) QPopBack() (v gossdb.Value, err error) {
	if self.ssdb == nil {
		return "", errors.New("ssdb == nil")
	}
	return self.ssdb.Qpop_back(self.mainkey)
}

func (self *SsdbQueue) QPushFront(value ...interface{}) (size int64, err error) {
	if self.ssdb == nil {
		return 0, errors.New("ssdb == nil")
	}
	return self.ssdb.Qpush_front(self.mainkey, value)
}

func (self *SsdbQueue) QpopFront() (v gossdb.Value, err error) {
	if self.ssdb == nil {
		return "", errors.New("ssdb == nil")
	}
	return self.ssdb.Qpop_front(self.mainkey)
}

func (self *SsdbQueue) QRange(offset int, limit int) (v []gossdb.Value, err error) {
	if self.ssdb == nil {
		return nil, errors.New("ssdb == nil")
	}
	return self.ssdb.Qrange(self.mainkey, offset, limit)
}

//=============================================================

type SsdbConf struct {
	Host             string
	Port             int
	MinPoolSize      int
	MaxPoolSize      int
	AcquireIncrement int
}

func InitSsdb(conf *SsdbConf) {
	if ssdb_pool != nil {
		return
	}
	if conf.Host == "" {
		log.Fatal(errors.New("没有配置ssdb服务器地址"))
		return
	}
	pool, err := gossdb.NewPool(&gossdb.Config{
		Host:             conf.Host,
		Port:             conf.Port,
		MinPoolSize:      conf.MinPoolSize,
		MaxPoolSize:      conf.MaxPoolSize,
		AcquireIncrement: conf.AcquireIncrement,
	})
	if err != nil {
		log.Fatal(err)
		return
	}
	ssdb_pool = pool
}

func NewSsdbClient() (client *gossdb.Client, err error) {
	return ssdb_pool.NewClient()
}
