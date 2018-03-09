package bat

import (
	_"errors"
	_"fmt"
	_"gosugar/qiniu"
	_"io/ioutil"
	_"log"
	"onegame-master/service/common"
	_"onegame-master/service/model"
	_"onegame-master/service/third"
	"os"
	"path/filepath"
	"strings"
	_  "github.com/Pallinder/go-randomdata"
	_"github.com/astaxie/beego/httplib"
	_"github.com/donnie4w/json4g"
	_ "github.com/go-sql-driver/mysql"
	_"github.com/seefan/to"
)

func createdb() {
	common.Syncdb()
}

//获取指定目录及所有子目录下的所有文件，可以匹配后缀过滤。
func WalkDir(dirPth, suffix string) (files []string, err error) {
	files = make([]string, 0, 30)
	suffix = strings.ToUpper(suffix)                                                     //忽略后缀匹配的大小写
	err = filepath.Walk(dirPth, func(filename string, fi os.FileInfo, err error) error { //遍历目录
		//if err != nil { //忽略错误
		// return err
		//}
		if fi.IsDir() { // 忽略目录
			return nil
		}
		if strings.HasSuffix(strings.ToUpper(fi.Name()), suffix) {
			files = append(files, filename)
		}
		return nil
	})
	return files, err
}

func HandleToolAPI() {
	args := os.Args
	for i, v := range args {
		if i == 1 {
			switch v {
			case "-robot":
				//方法
				os.Exit(0)
			}
		}
	}
}

func HandleDatabaseAPI() {
	args := os.Args
	for i, v := range args {
		if i == 1 {
			switch v {
			case "-syncdb":
				createdb()
				os.Exit(0)
			}
		}
	}
}
