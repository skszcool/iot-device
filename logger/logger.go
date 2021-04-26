package logger

import (
	"fmt"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
	"github.com/skszcool/iot-device/setting"
	"io"
	"path"
	"time"
)

var logger *logrus.Logger
var levelMapping = map[string]logrus.Level{
	"trace": logrus.TraceLevel,
	"debug": logrus.DebugLevel,
	"info":  logrus.InfoLevel,
	"warn":  logrus.WarnLevel,
	"error": logrus.ErrorLevel,
	"fatal": logrus.FatalLevel,
	"panic": logrus.PanicLevel,
}

func Setup() {
	logFilePath := getLogFilePath()
	logFileName := getLogFileName()

	// 日志文件
	fileName := path.Join(logFilePath, logFileName)

	// 实例化
	logger = logrus.New()

	// 设置 rotatelogs
	logWriter, err := rotatelogs.New(
		// 分割后的文件名称
		fileName+".%Y%m%d.log",

		// 生成软链，指向最新日志文件
		rotatelogs.WithLinkName(fileName),

		// 设置最大保存时间(7天)
		rotatelogs.WithMaxAge(2*24*time.Hour),

		// 设置日志切割时间间隔(1天)
		rotatelogs.WithRotationTime(24*time.Hour),
	)

	if err != nil {
		fmt.Println("err", err)
	}

	// 设置日志级别
	if level, ok := levelMapping[setting.AppSetting.LogLevel]; ok {
		logger.SetLevel(level)
	} else {
		logger.SetLevel(logrus.ErrorLevel)
	}
	logger.SetOutput(io.MultiWriter(logWriter))

	logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})
}

func GetInstance() *logrus.Logger {
	return logger
}

func Trace(args ...interface{}) {
	logger.Trace(args)
}

func Debug(args ...interface{}) {
	logger.Debug(args)
}

func Info(args ...interface{}) {
	logger.Info(args)
}

func Print(args ...interface{}) {
	logger.Print(args)
}

func Warn(args ...interface{}) {
	logger.Warn(args)
}

func Error(args ...interface{}) {
	logger.Error(args)
}

func Fatal(args ...interface{}) {
	logger.Fatal(args)
}

func Panic(args ...interface{}) {
	logger.Panic(args)
}

// getLogFilePath get the log file save path
func getLogFilePath() string {
	return fmt.Sprintf("%s%s", setting.AppSetting.RuntimeRootPath, setting.AppSetting.LogSavePath)
}

// getLogFileName get the save name of the log file
func getLogFileName() string {
	return setting.AppSetting.LogSaveName
}
