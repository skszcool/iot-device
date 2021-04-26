package util

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/skszcool/iot-device/setting"
	"mime/multipart"
	"os"
	"path"
	"strconv"
	"strings"
)

type ginHelper struct {
}

type ginUploadOption struct {
	AllowMaxFileSize        int64
	AllowFileType           []string
	FormatMaxFileSizeErrMsg func(allowSize int64) string
	FormatFileTypeErrMsg    func(fileType string, allowFileType []string) string
}

type GinUploadOption interface {
	apply(option *ginUploadOption)
}

type ginUploadOptionFunc struct {
	f func(*ginUploadOption)
}

func (fdo *ginUploadOptionFunc) apply(do *ginUploadOption) {
	fdo.f(do)
}

func newGinUploadOptionFunc(f func(*ginUploadOption)) *ginUploadOptionFunc {
	return &ginUploadOptionFunc{
		f: f,
	}
}

type GinUploadSingleFileOpt struct{}

func (opt *GinUploadSingleFileOpt) WithAllowMaxFileSize(size int) GinUploadOption {
	return newGinUploadOptionFunc(func(option *ginUploadOption) {
		option.AllowMaxFileSize = int64(size) << 20
	})
}

func (opt *GinUploadSingleFileOpt) WithAllowMaxFileType(fileType []string) GinUploadOption {
	return newGinUploadOptionFunc(func(option *ginUploadOption) {
		option.AllowFileType = fileType
	})
}

func (opt *GinUploadSingleFileOpt) WithFormatMaxFileSizeErrMsg(f func(allowSize int64) string) GinUploadOption {
	return newGinUploadOptionFunc(func(option *ginUploadOption) {
		option.FormatMaxFileSizeErrMsg = f
	})
}

func (opt *GinUploadSingleFileOpt) WithFormatFileTypeErrMsg(f func(string, []string) string) GinUploadOption {
	return newGinUploadOptionFunc(func(option *ginUploadOption) {
		option.FormatFileTypeErrMsg = f
	})
}

// 校验文件的合法性
func checkUploadFile(fileHeader *multipart.FileHeader, allOpts ...GinUploadOption) error {
	opts := &ginUploadOption{
		AllowMaxFileSize: 1000000 << 20, // 1000000M
		AllowFileType:    []string{},
		FormatMaxFileSizeErrMsg: func(allowSize int64) string {
			return fmt.Sprintf("上传大小不能超过%dM", allowSize/1024/1024)
		},
		FormatFileTypeErrMsg: func(fileType string, allowFileType []string) string {
			return fmt.Sprintf("允许的文件类型为%s", strings.Join(allowFileType, ","))
		},
	}

	for _, opt := range allOpts {
		opt.apply(opts)
	}

	allowSize := setting.ServerSetting.FileUploadMaxSize << 20
	if fileHeader.Size > int64(allowSize) {
		return errors.New("文件上传大小不能超过" + strconv.Itoa(setting.ServerSetting.FileUploadMaxSize) + "M")
	}

	if fileHeader.Size > opts.AllowMaxFileSize {
		return errors.New(opts.FormatMaxFileSizeErrMsg(opts.AllowMaxFileSize))
	}

	if len(opts.AllowFileType) > 0 {
		currentType := strings.TrimLeft(path.Ext(fileHeader.Filename), ".")
		isAllow := false
		for i := 0; i < len(opts.AllowFileType); i++ {
			if currentType == opts.AllowFileType[i] {
				isAllow = true
				break
			}
		}

		if !isAllow {
			return errors.New(opts.FormatFileTypeErrMsg(currentType, opts.AllowFileType))
		}
	}

	return nil
}

// 单文件上传
func (ginHelper *ginHelper) SingleFile(c *gin.Context, filename string, allOpts ...GinUploadOption) (*multipart.FileHeader, error) {
	fileHeader, err := c.FormFile(filename)
	if err != nil {
		return fileHeader, err
	}

	return fileHeader, checkUploadFile(fileHeader, allOpts...)
}

// 多文件上传
func (ginHelper *ginHelper) MultiFile(c *gin.Context, filename string, allOpts ...GinUploadOption) ([]*multipart.FileHeader, error) {
	form, err := c.MultipartForm()
	if err != nil {
		return []*multipart.FileHeader{}, err
	}

	files := form.File[filename]

	totalSize := int64(0)
	for _, fileHeader := range files {
		if err := checkUploadFile(fileHeader, allOpts...); err != nil {
			return files, err
		}

		totalSize += fileHeader.Size
	}

	allowSize := setting.ServerSetting.FileUploadMaxSize << 20
	if totalSize > int64(allowSize) {
		return files, errors.New("文件上传大小不能超过" + strconv.Itoa(setting.ServerSetting.FileUploadMaxSize) + "M")
	}

	return files, err
}

func (ginHelper *ginHelper) DownloadFile(c *gin.Context, downloadFilePath string) error {
	//downloadFile := path.Join(setting.GetDeviceEventPath(), "restart_supervisor.zip")
	downloadFileOpen, err := os.Open(downloadFilePath)
	if err != nil {
		return err
	}

	fileInfo, _ := downloadFileOpen.Stat()
	extraHeaders := map[string]string{
		"Content-Disposition": fmt.Sprintf(`attachment; filename="%s"`, path.Base(downloadFilePath)),
	}

	fmt.Println(downloadFilePath)
	fmt.Println(extraHeaders)

	c.DataFromReader(200, fileInfo.Size(), "", downloadFileOpen, extraHeaders)

	return nil
}
