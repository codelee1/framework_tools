package conf

import (
	"flag"
	"fmt"
	"github.com/BurntSushi/toml"
	"learn_tools/gin_demo/library/database/orm"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var (
	confPath string
	Conf     = &Config{}
	Addr     string
)

type Config struct {
	App        *App
	HttpServer *HttpServer
	Log        *Log
	DB         *orm.Config
}

type App struct {
	CacheDir string
	Version  int
	Pid      string
}

type HttpServer struct {
	Addr string
	Port int
}

type Log struct {
	Dir    string
	Stdout bool
}

func Init() (err error) {
	if confPath != "" {
		path := GetAppPath()
		_, err = toml.DecodeFile(path+"/"+confPath, &Conf)
		if err != nil {
			return
		}

	}
	if e := os.MkdirAll(GetLogDir(), os.ModePerm); e != nil {
		err = fmt.Errorf("mdkir log error! ", err)
	}
	if e := os.MkdirAll(GetCacheDir(), os.ModePerm); e != nil {
		err = fmt.Errorf("mdkir cache error! ", err)
	}
	return
}

func GetAppPath() string {
	fmt.Println(os.Args)
	file, _ := exec.LookPath(os.Args[0])
	path, _ := filepath.Abs(file)
	index := strings.LastIndex(path, string(os.PathSeparator))

	return path[:index]
}

func GetLogDir() string {
	if Conf.Log.Dir == "" || Conf.Log.Dir == "./" {
		return GetAppPath() + "/log/"
	}
	return Conf.Log.Dir
}

func GetCacheDir() string {
	if Conf.App.CacheDir == "" || Conf.App.CacheDir == "./" {
		return GetAppPath() + "/cache/"
	}
	return Conf.App.CacheDir
}

func init() {
	flag.StringVar(&confPath, "conf", "../conf/conf.toml", "config path")
	flag.StringVar(&Addr, "addr", ":8080", "ip:port")
}
