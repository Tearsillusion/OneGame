package main

import (
	"fmt"
	"runtime"
	"onegame-master/service/bat"
	"onegame-master/service/third"
	"onegame-master/service/common"
	"github.com/gogap/logrus"
	"github.com/gogap/logrus/hooks/file"
	"onegame-master/service/keywords"
	"onegame-master/service"
	"github.com/astaxie/beego/orm"
	"github.com/robfig/cron"
	"github.com/gin-gonic/gin"
	"onegame-master/service/test"
	"os"
	"log"
	"errors"
	"net/http"
)

type Rounters struct {
	Path    string
	Handler gin.HandlerFunc
}
func options(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": 200,
	})
}
func init(){
	fmt.Println("Go SDK Version:",runtime.Version())

	// ssdb 配置
	third.LoadConfig("config/web.conf")
	ssdb_host := third.Config.GetString("ssdb.host")
	ssdb_port := third.Config.GetInt64("ssdb.port")
	ssdb_MinPoolSize := third.Config.GetInt64("ssdb.MinPoolSize")
	ssdb_MaxPoolSize := third.Config.GetInt64("ssdb.MaxPoolSize")
	ssdb_AcquireIncrement := third.Config.GetInt64("ssdb.AcquireIncrement")
	common.InitSsdb(&common.SsdbConf{
		Host:             ssdb_host,
		Port:             int(ssdb_port),
		MinPoolSize:      int(ssdb_MinPoolSize),
		MaxPoolSize:      int(ssdb_MaxPoolSize),
		AcquireIncrement: int(ssdb_AcquireIncrement),
	})

	// 日志
	logrus.SetLevel(logrus.InfoLevel)
	logrus.AddHook(file.NewHook("log/log.txt"))

	// 加载敏感词库
	sensitives.LoadTokens()

	// mysql数据库连接
	db_host := third.Config.GetString("db_host")
	db_port := third.Config.GetString("db_port")
	db_user := third.Config.GetString("db_user")
	db_pass := third.Config.GetString("db_pass")
	db_name := third.Config.GetString("db_name")

	debug := bat.HasArg("-debug")
	//注册model
	service.RegisterAllModel()
	common.InitOrm("default", debug)
	common.ConnectMySQL(db_user, db_pass, db_name, db_host+":"+db_port, 500, 500)

	// 命令行参数处理
	bat.HandleToolAPI()
	// 处理数据库命令行
	bat.HandleDatabaseAPI()

	if !debug {
		gin.SetMode(gin.ReleaseMode)
		// 不把日志打印到控制台
		logrus.SetOutput(&bat.EmptyWriter{})
		// orm数据库日志输出到logrus日志系统
		orm.DebugLog = orm.NewLog(&bat.LogrusWriter{})
	}else {
		gin.SetMode(gin.DebugMode)
	}
	iszip:=false
	c := cron.New()
	if iszip {
		c.AddFunc("* * 1 * * ", func() {
			if filelist, err := bat.CompressLogToZip(); err != nil {
				fmt.Printf("日志压缩失败： %s\n", err.Error())
			} else {
				if err := bat.ClearLogYesterday(filelist); err != nil {
					fmt.Printf("日志清除失败: %s\n", err.Error())
				}
			}
		})
		c.Start()
	}

}
func main(){
	//查看权限
	//gin.BasicAuthForRealm(bat.MakeAuthBasic(),"misasky's zone")
	rounter:=gin.Default()
	//管理员接口
	admin:=rounter.Group("/admin",gin.BasicAuth(gin.Accounts{
		third.Config.GetString("basic_auth"):third.Config.GetString("auth_pwd"),
	}))
	aapis:=[]Rounters{
		{"/secrets",test.Secrets},
	}
	for _,v:=range aapis{
		admin.GET(v.Path,v.Handler)
	}
	//普通用户接口
	apis:=[]Rounters{
		{"/login",test.Secrets},
	}
	api:=rounter.Group("/api")
	for _,v:=range apis{
		api.OPTIONS(v.Path,options)
		api.POST(v.Path,v.Handler)
	}
	// 启动服务器
	use_http := true
	port := ""
	if len(os.Args) < 3 {
		log.Fatal(errors.New("缺少服务器启动参数[-http 80 / -https 80]"))
		return
	}
	for i, v := range os.Args {
		if i == 1 {
			switch v {
			case "-http":
				use_http = true
				port = os.Args[2]
			case "-https":
				use_http = false
				port = os.Args[2]
			default:
				log.Fatal(errors.New("无效的命令行参数"))
			}
		}
	}
	if use_http {
		log.Fatal(http.ListenAndServe(":"+port, nil))
	} else {
		//https证书 还没有咯放着吧
		//log.Fatal(http.ListenAndServeTLS(":"+port, "config/ssl/onegame.crt", "config/ssl/onegame.key", nil))
	}
	rounter.Run(":"+port)
}

