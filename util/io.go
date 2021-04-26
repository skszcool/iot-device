package util

import (
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/mholt/archiver/v3"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

type ioHelper struct {
}

// 判断文件或目录是否存在
func (ioHelper *ioHelper) IsExistsFileOrDir(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

// zip压缩
func (ioHelper *ioHelper) Zip(sources []string, destination string) error {
	z := archiver.Zip{ImplicitTopLevelFolder: true}
	return z.Archive(sources, destination)
}

// unzip解压
func (ioHelper *ioHelper) Unzip(sources string, destination string) error {
	z := archiver.Zip{ImplicitTopLevelFolder: false, MkdirAll: true}
	return z.Unarchive(sources, destination)
}

// 获取目录下所有的子目录
func (ioHelper *ioHelper) GetDirAllDirNames(dirPath string) ([]string, error) {
	//dir, err := ioutil.ReadDir(dirPath)
	dir, err := os.ReadDir(dirPath)

	if err != nil {
		return nil, err
	}
	var dirs []string

	// 获取path分隔符
	//PthSep := string(os.PathSeparator)

	for _, fi := range dir {
		if fi.IsDir() {
			dirs = append(dirs, fi.Name())
		}
	}

	return dirs, nil
}

// 获取目录下所有的子目录和文件
func (ioHelper *ioHelper) GetDirAllPath(dirPath string) ([]string, error) {
	dir, err := os.ReadDir(dirPath)
	var result []string

	if err != nil {
		return result, err
	}

	for _, fi := range dir {
		result = append(result, path.Join(dirPath, fi.Name()))
	}

	return result, nil
}

type TypeDirDeep struct {
	Title    string        `json:"title"`
	Key      string        `json:"key"`
	IsFile   bool          `json:"isFile"`
	IsDir    bool          `json:"isDir"`
	Children []TypeDirDeep `json:"children"`
}

// 遍历目录结构(包括文件)
func (ioHelper *ioHelper) ReadDirDeep(dirPath string, excludeName []string, delimiter string, key string) (TypeDirDeep, error) {
	pathArr := strings.Split(dirPath, delimiter)
	dirName := pathArr[len(pathArr)-1]

	if key == "" {
		key = dirName
	}
	data := TypeDirDeep{
		Title:    dirName,
		Key:      key,
		IsFile:   false,
		IsDir:    false,
		Children: []TypeDirDeep{},
	}

	if len(excludeName) > 0 {
		for _, v := range excludeName {
			if matched, _ := regexp.Match(dirPath, []byte(v+"$")); matched {
				return data, errors.New("跳出过滤目录")
			}
		}
	}

	if ioHelper.IsExistsFileOrDir(dirPath) {
		fileInfo, _ := os.Stat(dirPath)
		if fileInfo.IsDir() {
			data.IsDir = true
			data.IsFile = false
			fileInfos, _ := os.ReadDir(dirPath)
			for _, file := range fileInfos {
				childDir := file.Name()
				newDirPath := path.Join(dirPath, childDir)
				if key == "" {
					key = dirName + delimiter + childDir
				} else {
					key = key + delimiter + childDir
				}

				child, _ := ioHelper.ReadDirDeep(newDirPath, excludeName, delimiter, key)

				data.Children = append(data.Children, child)
			}
		} else {
			data.IsDir = false
			data.IsFile = true
			data.Title = dirName
			if key == "" {
				data.Key = dirName
			} else {
				data.Key = key
			}
		}
	}

	return data, nil
}

//拷贝文件夹,同时拷贝文件夹中的文件
func (ioHelper *ioHelper) CopyDir(srcPath string, destPath string) error {
	//检测目录正确性
	if srcInfo, err := os.Stat(srcPath); err != nil {
		return err
	} else {
		if !srcInfo.IsDir() {
			e := errors.New("srcPath不是一个正确的目录！")
			return e
		}
	}
	if destInfo, err := os.Stat(destPath); err != nil {
		return err
	} else {
		if !destInfo.IsDir() {
			e := errors.New("destInfo不是一个正确的目录！")
			return e
		}
	}
	err := filepath.Walk(srcPath, func(path string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}
		if !f.IsDir() {
			path := strings.Replace(path, "\\", "/", -1)
			destNewPath := strings.Replace(path, srcPath, destPath, -1)
			ioHelper.CopyFile(path, destNewPath)
		}
		return nil
	})

	return err
}

//生成目录并拷贝文件
func (ioHelper *ioHelper) CopyFile(src, dest string) (w int64, err error) {
	srcFile, err := os.Open(src)
	if err != nil {
		return
	}
	defer srcFile.Close()

	//分割path目录
	destSplitPathDirs := strings.Split(dest, "/")

	//检测时候存在目录
	destSplitPath := ""
	for index, dir := range destSplitPathDirs {
		if index < len(destSplitPathDirs)-1 {
			destSplitPath = destSplitPath + dir + "/"
			b := ioHelper.IsExistsFileOrDir(destSplitPath)
			if b == false {
				//创建目录
				err := os.Mkdir(destSplitPath, os.ModePerm)
				if err != nil {
					fmt.Println(err)
				}
			}
		}
	}
	dstFile, err := os.Create(dest)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer dstFile.Close()

	return io.Copy(dstFile, srcFile)
}

func (ioHelper *ioHelper) FileToBase64(filenamePath string) (string, error) {
	data, err := ioutil.ReadFile(filenamePath)
	if err != nil {
		return "", err
	}
	base64Str := base64.StdEncoding.EncodeToString(data)

	return base64Str, nil
}

func (ioHelper *ioHelper) Base64ToFile(base64Str string, filenamePath string) error {
	decodeData, err := base64.StdEncoding.DecodeString(base64Str)
	if err != nil {
		return err
	}
	f, err := os.OpenFile(filenamePath, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}
	defer f.Close()
	if _, err := f.Write(decodeData); err != nil {
		return err
	}

	return nil
}

// 创建文件
func (ioHelper *ioHelper) CreateFile(filenamePath string) error {
	f, err := os.Create(filenamePath)
	defer f.Close()
	return err
}
