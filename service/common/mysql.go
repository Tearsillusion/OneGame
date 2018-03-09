package common

import (
	"gosugar/db"

	"github.com/astaxie/beego/orm"
)

var (
	sql *db.Orm
)

func InitOrm(aliasName string, print_log bool) {
	if sql != nil {
		return
	}
	sql = db.NewOrm(aliasName)
	orm.Debug = print_log
}

func NewOrm() orm.Ormer {
	return sql.NewOrm()
}

func ConnectMySQL(username string, password string, dbname string, mysql_ip string, param ...int) {
	sql.ConnectMySQL(username, password, dbname, mysql_ip, param...)
}

func RegisterModel(model ...interface{}) {
	orm.RegisterModel(model...)
}

func Syncdb() {
	sql.BuildDB(false, true)
}
