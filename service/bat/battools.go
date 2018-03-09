package bat

import (
	"os"
	"github.com/gogap/logrus"
	"time"
	"fmt"
	"path/filepath"
	"strings"
	"archive/zip"
	"io/ioutil"
	"bytes"
	"github.com/gin-gonic/gin"
)

func HasArg(key string) bool {
	for _, v := range os.Args {
		if v == key {
			return true
		}
	}
	return false
}

type EmptyWriter struct {
}

func (self *EmptyWriter) Write(p []byte) (n int, err error) {
	// nothing to do
	return len(p), nil
}

type LogrusWriter struct {
}

func (self *LogrusWriter) Write(p []byte) (n int, err error) {
	logrus.Infof(string(p))
	return len(p), nil
}

//打包zip
func CompressLogToZip() (filelist []string, err error) {
	tm := time.Unix(time.Now().Unix()-60*60*24, 0)
	folder := fmt.Sprintf("%d-%.2d-%.2d", tm.Year(), tm.Month(), tm.Day())
	if err := os.MkdirAll("log/"+folder, 0777); err != nil {
		return nil, err
	}
	filelist = make([]string, 0, 32)
	if err := filepath.Walk("log", func(path string, fileinfo os.FileInfo, err error) error {
		if fileinfo == nil {
			return err
		}
		if fileinfo.IsDir() {
			return nil
		}
		//fmt.Printf("搜索到文件: %s\n", path)
		if strings.Contains(path, folder) {
			filelist = append(filelist, path)
			array := strings.Split(path, "/")
			fileName := array[len(array)-1]
			//fmt.Printf("需要压缩的文件: %s\n", fileName)

			buf := new(bytes.Buffer)
			w := zip.NewWriter(buf)
			if f, err := w.Create(fileName); err != nil {
				return err
			} else {
				if bin, err := ioutil.ReadFile(path); err != nil {
					return err
				} else {
					if _, err = f.Write(bin); err != nil {
						return err
					}
				}
			}
			if err := w.Close(); err != nil {
				return err
			}

			if f, err := os.OpenFile(fmt.Sprintf("log/%s/%.2dh%.2dm.zip", folder, fileinfo.ModTime().Hour(), fileinfo.ModTime().Minute()), os.O_CREATE|os.O_WRONLY, 0666); err != nil {
				return err
			} else {
				buf.WriteTo(f)
			}
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return filelist, nil
}
//删除昨天的日志
func ClearLogYesterday(filelist []string) error {
	for _, v := range filelist {
		//fmt.Printf("删除的文件: %s\n", v)
		os.Remove(v)
	}
	tm := time.Unix(time.Now().Unix()-60*60*24*7, 0)
	folder := fmt.Sprintf("%d-%.2d-%.2d", tm.Year(), tm.Month(), tm.Day())
	os.RemoveAll(folder)
	return nil
}


func MakeAuthBasic()(authbasic gin.Accounts){
	authbasic["userid"]="misasky"
	authbasic["password"]="misasky-2018-03-02"
	return authbasic
}