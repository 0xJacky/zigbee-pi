package settings

import (
	"gopkg.in/ini.v1"
	"log"
)

var Conf *ini.File

type App struct {
	PageSize     int
	TrustedToken string
}

var AppSettings = &App{}

type Server struct {
	HttpPort string
	RunMode  string
}

var ServerSettings = &Server{}

func init() {
	var err error
	Conf, err = ini.Load("app.ini")
	if err != nil {
		log.Fatalf("setting.init, fail to parse 'app.ini': %v", err)
	}

	mapTo("app", AppSettings)
	mapTo("server", ServerSettings)
}

func mapTo(section string, v interface{}) {
	err := Conf.Section(section).MapTo(v)
	if err != nil {
		log.Fatalf("setting.mapTo %s err: %v", section, err)
	}
}
