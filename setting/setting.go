package setting

import (
	"log"
	"net/http"
	"time"

	"github.com/go-ini/ini"
)

type App struct {
	JwtSecret                 string
	RuntimeRootPath           string
	LogSavePath               string
	LogSaveName               string
	LogLevel                  string
	PluginMonitorIntervalSpec string
}

var AppSetting = &App{
	LogLevel: "error",
}

type Server struct {
	RunMode             string
	IsStartCors         bool
	HttpPort            int
	HttpsPort           int
	ReadTimeout         time.Duration
	WriteTimeout        time.Duration
	FileUploadMaxSize   int
	AllowPluginFileSize int
	AllowPluginFileType string
}

var ServerSetting = &Server{}

type Database struct {
	IsStartLog  bool
	Type        string
	TablePrefix string
	SQLiteDB    string
}

var DatabaseSetting = &Database{}

type Cookie struct {
	Name     string
	Path     string
	Secure   bool
	HttpOnly bool
	MaxAge   int
	SameSite string
}

func (c *Cookie) GetSameSite() http.SameSite {
	switch c.SameSite {
	case "strict":
		return http.SameSiteStrictMode
	case "lax":
		return http.SameSiteLaxMode
	case "none":
		return http.SameSiteNoneMode
	default:
		return http.SameSiteDefaultMode

	}
}

var CookieSetting = &Cookie{}

type IotPlatform struct {
	SimulationLoginSecret        string
	HeartbeatIntervalSpec        string
	DeviceSOSRequestIntervalSpec string
	IsStartIotPlatformCheck      bool
}

var IotPlatformSetting = &IotPlatform{}

type Frpc struct {
	IsStart        bool
	FrpcInstallDir string
}

var FrpcSetting = &Frpc{}

// 获取static目录
func GetStaticPath() string {
	return AppSetting.RuntimeRootPath + "static"
}

var cfg *ini.File

// Setup initialize the configuration instance
func Setup() {
	var err error
	cfg, err = ini.Load("conf/app.ini")
	if err != nil {
		log.Fatalf("setting.Setup, fail to parse 'conf/app.ini': %v", err)
	}

	mapTo("app", AppSetting)
	mapTo("server", ServerSetting)
	mapTo("database", DatabaseSetting)
	mapTo("cookie", CookieSetting)
	mapTo("iotPlatform", IotPlatformSetting)
	mapTo("frpc", FrpcSetting)

	ServerSetting.ReadTimeout = ServerSetting.ReadTimeout * time.Second
	ServerSetting.WriteTimeout = ServerSetting.WriteTimeout * time.Second
}

// 是否为debug模式
func IsDebugMode() bool {
	return ServerSetting.RunMode == "debug"
}

// mapTo map section
func mapTo(section string, v interface{}) {
	err := cfg.Section(section).MapTo(v)
	if err != nil {
		log.Fatalf("Cfg.MapTo %s err: %v", section, err)
	}
}
