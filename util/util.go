package util

import (
	uuid "github.com/satori/go.uuid"
	"github.com/skszcool/iot-device/setting"
	"github.com/skszcool/iot-device/shortid"
	"reflect"
	"regexp"
	"strings"
)

var Crypto = new(crypto)
var Array = new(array)
var IoHelper = new(ioHelper)
var GinHelper = new(ginHelper)

// Setup Initialize the util
func Setup() {
	jwtSecret = []byte(setting.AppSetting.JwtSecret)
}

// 生成用户密码
func MkUserPassword(password string) string {
	return Crypto.Md5(password)
}

// 生成uuid
func MkUUID() string {
	return uuid.NewV5(uuid.NewV1(), "sksz").String()
}

// 生成简短的id
func MkShortId() string {
	if id, err := shortid.Generate(); err != nil {
		return MkShortId()
	} else {
		return id
	}
}

// 替换模板字符串
/*
tplStr := `aiStreaming. ! queue leaky=1 ! tee name=${ai,ai2} ${ai,ai2}. \
! queue !  videoscale ! video/x-raw,width=${width,480},height=${height,360} ! \
videoconvert ! ${encoded,vaapih264enc} !  flvmux name=${ssrc,} ! rtmpsink location=rtmp://127.0.0.1:1936/live/${ssrc,}`
	newTplStr := util.ReplaceTplStr(tplStr, map[string]string{
		"ssrc": "ssrc_xxxxx",
	})
	fmt.Printf("newTplStr = %s\n", newTplStr)
*/
func ReplaceTplStr(tplStr string, matchData map[string]string) string {
	r := regexp.MustCompile(`\$\{.*?\}`)
	for {
		matchResult := r.FindStringSubmatch(tplStr)
		if len(matchResult) > 0 {
			for _, oldStr := range matchResult {
				origin := strings.Split(strings.TrimSuffix(strings.TrimPrefix(oldStr, "${"), "}"), ",")
				key := strings.TrimSpace(origin[0])
				newStr := strings.TrimSpace(origin[1])

				if val, ok := matchData[key]; ok {
					newStr = val
				}

				tplStr = strings.Replace(tplStr, oldStr, newStr, 1)
			}
		} else {
			break
		}
	}

	return tplStr
}

// 将map数据结构转换为
func MapToStruct(data map[string]interface{}, inStructPtr interface{}) {
	rType := reflect.TypeOf(inStructPtr)
	rVal := reflect.ValueOf(inStructPtr)
	if rType.Kind() == reflect.Ptr {
		// 传入的inStructPtr是指针，需要.Elem()取得指针指向的value
		rType = rType.Elem()
		rVal = rVal.Elem()
	} else {
		panic("inStructPtr must be ptr to struct")
	}
	// 遍历结构体
	for i := 0; i < rType.NumField(); i++ {
		t := rType.Field(i)
		f := rVal.Field(i)
		// 得到tag中的字段名
		key := t.Tag.Get("key")
		if v, ok := data[key]; ok {
			// 检查是否需要类型转换
			dataType := reflect.TypeOf(v)
			structType := f.Type()
			if structType == dataType {
				f.Set(reflect.ValueOf(v))
			} else {
				if dataType.ConvertibleTo(structType) {
					// 转换类型
					f.Set(reflect.ValueOf(v).Convert(structType))
				} else {
					panic(t.Name + " type mismatch")
				}
			}
		} else {
			panic(t.Name + " not found")
		}
	}
}
