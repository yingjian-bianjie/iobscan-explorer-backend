package repository

import (
	"github.com/bianjieai/iobscan-explorer-backend/internal/common/constant"
	"strings"
	"time"

	"github.com/bianjieai/iobscan-explorer-backend/internal/common/configs"
	"github.com/bianjieai/iobscan-explorer-backend/pkg/logger"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/mgo.v2/txn"
)

type (
	Docs interface {
		Name() string
		PkKvPair() map[string]interface{}
	}

	DocV2 interface {
		Docs
		EnsureIndexes()
	}
)

var (
	session *mgo.Session
)

const CollectionNameTxn = "sync_mgo_txn"

var (
	conf     configs.DataBaseConf
	_srvConf *configs.Server
)

func GetSrvConf() *configs.Server {
	if _srvConf != nil {
		return _srvConf
	}
	panic("Start() do not work")
}
func Start(dataCfg *configs.Config) {
	conf = dataCfg.DataBaseConf
	_srvConf = &dataCfg.Server
	initDbClient(dataCfg.DataBaseConf)
}
func initDbClient(dataBaseConf configs.DataBaseConf) {
	addrs := strings.Split(dataBaseConf.Addrs, ",")
	dialInfo := &mgo.DialInfo{
		Addrs:     addrs,
		Database:  dataBaseConf.Database,
		Username:  dataBaseConf.User,
		Password:  dataBaseConf.Passwd,
		Direct:    true,
		Timeout:   time.Second * 10,
		PoolLimit: 4096, // Session.SetPoolLimit
	}

	var err error
	session, err = mgo.DialWithInfo(dialInfo)
	if err != nil {
		logger.Fatal("connect db fail", logger.String("err", err.Error()))
	}
	session.SetMode(mgo.Strong, true)
	logger.Info("init db success")
}

func DbStatus() bool {
	if session == nil {
		return false
	}
	err := session.Ping()
	if err != nil {
		return false
	}
	return true
}

func Stop() {
	logger.Info("release resource :mongoDb")
	if session != nil {
		session.Close()
	}
}

func GetSession() *mgo.Session {
	// max session num is 4096
	return session.Clone()
}

func GetDataBase() string {
	return conf.Database
}

// get collection object
func ExecCollection(collectionName string, s func(*mgo.Collection) error) error {
	session := GetSession()
	defer session.Close()
	c := session.DB(conf.Database).C(collectionName)
	return s(c)
}

func Save(h Docs) error {
	save := func(c *mgo.Collection) error {
		pk := h.PkKvPair()
		n, _ := c.Find(pk).Count()
		if n >= 1 {
			return constant.ErrDbExist
		}
		return c.Insert(h)
	}
	return ExecCollection(h.Name(), save)
}

//mgo transaction method
//detail to see: https://godoc.org/gopkg.in/mgo.v2/txn
func Txn(ops []txn.Op) error {
	session := GetSession()
	defer session.Close()

	c := session.DB(conf.Database).C(CollectionNameTxn)
	runner := txn.NewRunner(c)

	txObjectId := bson.NewObjectId()
	err := runner.Run(ops, txObjectId, nil)
	if err != nil {
		if err == txn.ErrAborted {
			err = runner.Resume(txObjectId)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}

	return nil
}
