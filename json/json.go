package json

import jsoniter "github.com/json-iterator/go"

// 100%兼用encoding/json包，但性能提升6倍多，占用更少的内存
var json = jsoniter.ConfigCompatibleWithStandardLibrary

var (
	Marshal             = json.Marshal
	Unmarshal           = json.Unmarshal
	UnmarshalFromString = json.UnmarshalFromString
	MarshalIndent       = json.MarshalIndent
	NewDecoder          = json.NewDecoder
	NewEncoder          = json.NewEncoder
	MarshalToString     = json.MarshalToString
	Get                 = json.Get
	Valid               = json.Valid
	RegisterExtension   = json.RegisterExtension
	DecoderOf           = json.DecoderOf
	EncoderOf           = json.EncoderOf
)
