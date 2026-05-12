package file

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"os/exec"

	log "github.com/sirupsen/logrus"

	//"archive/zip"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/yeka/zip"
)

// ReadFile 读取文件到二进制
func ReadFile(path string) ([]byte, error) {
	content, err := os.ReadFile(path)
	return content, err
}

// WriteFile 写入二进制到文件
func WriteFile(fileName string, byteArray []byte) error {
	var f *os.File
	var err error
	if CheckPathExist(fileName) { //文件存在
		f, err = os.OpenFile(fileName, os.O_APPEND, 0666) //打开文件
		if err != nil {
			return err
		}
	} else { //文件不存在
		f, err = os.Create(fileName) //创建文件
		if err != nil {
			return err
		}
	}
	defer f.Close()
	//将文件写进去
	_, err1 := io.WriteString(f, string(byteArray))
	if err1 != nil {
		return err1
	}
	return nil
}

// CopyFile 复制单个文件
func CopyFile(src, dst string) error {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destination.Close()

	if _, err := io.Copy(destination, source); err != nil {
		return err
	}

	return os.Chmod(dst, sourceFileStat.Mode())
}

// 生成临时文件
func SaveTmpFile(filename, scriptContent string) (string, error) {
	f, err := ioutil.TempFile("", filename)
	defer f.Close()
	if err != nil {
		return "", err
	}
	_, err = f.WriteString(scriptContent)
	if err != nil {
		return "", err
	}
	return f.Name(), nil
}

// CheckPathExist 验证 文件/目录 是否存在
func CheckPathExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

// IsDir 判断所给路径是否为文件夹
func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

// IsFile 判断所给路径是否为文件
func IsFile(path string) bool {
	return !IsDir(path)
}

// CreateDir 调用os.MkdirAll递归创建文件夹
func CreateDir(path string) error {
	if !CheckPathExist(path) {
		err := os.MkdirAll(path, os.ModePerm)
		return err
	}
	return nil
}

// RemoveFile 删除文件
func RemoveFile(path string) error {
	return os.Remove(path)
}

// GetFileTmpName 获取一个随机的临时文件名
func GetFileTmpName(preString string, rand int) string {
	timeUnixNano := time.Now().UnixNano()
	timeString := strconv.FormatInt(timeUnixNano, 10)

	return preString + "_" + timeString + "_" + GetRandomString(rand)
}

// GetRandomString 获取随机n位字符串
func GetRandomString(n int) string {
	rand.Seed(time.Now().UnixNano())
	str := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	bytes := []byte(str)
	var result []byte
	for i := 0; i < n; i++ {
		result = append(result, bytes[rand.Intn(len(bytes))])
	}
	return string(result)
}

// GetTargetFileDir 获取指定文件位置
func GetTargetFileDir(dest string, targetFileName string) (string, error) {
	dir, err := os.Open(dest)
	if err != nil {
		return "", err
	}
	fileList, err := dir.Readdir(-1)
	if err != nil {
		return "", err
	}

	for _, file := range fileList {
		path := dest + string(os.PathSeparator) + file.Name()

		// 匹配到
		if file.Name() == targetFileName {
			return path, nil
		}

		if file.IsDir() {
			res, err := GetTargetFileDir(path, targetFileName)
			if err != nil {
				return "", err
			}
			if res != "" {
				return res, nil
			}
		}
	}

	return "", nil
}

func ScanFile(dest string, isNext bool, result []string) []string {
	dir, err := os.Open(dest)
	if err != nil {
		return nil
	}
	fileList, err := dir.Readdir(-1)
	if err != nil {
		return nil
	}

	for _, file := range fileList {
		if file.Name() == "__MACOSX" {
			continue
		}

		path := dest + string(os.PathSeparator) + file.Name()

		if isNext == true && file.IsDir() {
			res := ScanFile(path, isNext, make([]string, 0))
			if res != nil {
				result = append(result, res...)
			}
		} else {
			result = append(result, path)
		}
	}

	return result
}

// GetFilesFromDirectory 获取某个目录下的所有文件
func GetFilesFromDirectory(directory string) ([]string, error) {
	fileList := make([]string, 0)
	files, err := ioutil.ReadDir(directory)
	if err != nil {
		fmt.Println("error reading directory:", err)
		return fileList, err
	}
	for _, file := range files {
		//过滤
		if strings.Contains(file.Name(), "__MACOSX") {
			continue
		}
		//相对文件路径
		filePath := filepath.Join(directory, file.Name())
		//处理目录中的文件
		if !file.IsDir() {
			fileList = append(fileList, filePath)
			continue
		}
		//处理目录中的目录
		subFileList, err := GetFilesFromDirectory(filePath)
		if err != nil {
			continue
		}
		fileList = append(fileList, subFileList...)
	}
	return fileList, nil
}

// Unzip 解压缩设置密码的zip包
func Unzip(filePath, targetDir, passwd string) error {
	//清空目录
	if err := os.RemoveAll(targetDir); err != nil {
		return err
	}
	// 1、使用zip.OpenReader打开zip文件
	reader, err := zip.OpenReader(filePath)
	if err != nil {
		return err
	}
	defer reader.Close()

	// 2、循环访问 zip 中的文件 zip.File 切片
	for _, file := range reader.File {
		if file.IsEncrypted() {
			file.SetPassword(passwd)
		}
		if strings.Contains(file.Name, "__MACOSX") {
			continue
		}
		targetPath := filepath.Join(targetDir, file.Name)

		//创建目录
		if file.FileInfo().IsDir() {
			if err := os.MkdirAll(targetPath, os.ModePerm); err != nil {
				return err
			}
			continue
		}

		if err := os.MkdirAll(filepath.Dir(targetPath), os.ModePerm); err != nil {
			return err
		}

		// 3、使用 zip.File.Open 方法读取 zip 中文件的内容
		fileReader, err := file.Open()
		if err != nil {
			return err
		}
		defer fileReader.Close()

		targetFile, err := os.OpenFile(targetPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return err
		}
		defer targetFile.Close()

		// 4、使用 io.Copy 或 io.Writer.Write 保存解压后的文件内容
		if _, err := io.Copy(targetFile, fileReader); err != nil {
			return err
		}
	}

	return nil
}

// VerifyVersion 版本校验
func VerifyVersion(param map[string]string) error {
	oldVersion, _ := param["oldVersion"]
	unzipDir, _ := param["unzipDir"]
	filePtr, err := os.Open(filepath.Join(unzipDir, "config.json"))
	if err != nil {
		fmt.Println("Fail to read config file:" + err.Error())
		return err
	}
	defer filePtr.Close()
	// 创建json解码器
	decoder := json.NewDecoder(filePtr)
	var configData map[string]interface{}
	err = decoder.Decode(&configData)
	if err != nil {
		fmt.Println("JSON parsing err:" + err.Error())
		return err
	}

	version := configData["version"]
	versionStr := fmt.Sprintf("%v", version)
	// 验证版本号是否一致
	if versionStr != oldVersion {
		// 判断是否是升级包
		updateVersion := configData["update_version"]
		if updateVersion == nil {
			return errors.New("版本号不一致")
		}
		// 验证升级包版本号是否一致
		updateVersionStr := fmt.Sprintf("%v", updateVersion)
		if updateVersionStr != oldVersion {
			return errors.New("版本号不一致")
		}
	}

	return nil
}

// VerifyFileUniformity 验证文件一致性
func VerifyFileUniformity(fileList []string) error {
	for _, filePath := range fileList {
		fileObj, err := os.Open(filePath)
		if err != nil {
			return err
		}
		defer fileObj.Close()
		// 创建一个新的MD5哈希接口
		hash := md5.New()
		// 将文件内容复制到哈希接口中
		if _, err := io.Copy(hash, fileObj); err != nil {
			return err
		}
		// 计算哈希值
		hashInBytes := hash.Sum(nil)[:16]
		// 将字节切片转换为16进制字符串
		md5String := hex.EncodeToString(hashInBytes)

		// 获取文件名
		fileName := filepath.Base(filePath)
		if fileName == "smart" || fileName == "config.json" {
			continue
		}
		if !strings.Contains(fileName, md5String) {
			return errors.New("文件已被篡改")
		}
	}
	return nil
}

// InvokeUpgradeScript 调用升级脚本
func InvokeUpgradeScript(scriptPath string) error {
	// 给脚本赋予执行权限
	err := os.Chmod(scriptPath, 0755)
	if err != nil {
		log.Error("Failed to chmod script:", err)
		return err
	}

	cmd := exec.Command("bash", scriptPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		log.Error("Failed to execute script:", err)
		return err
	}
	return nil
}

// GetFileSize 获取文件大小
func GetFileSize(path string) (int64, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return 0, err
	}
	return fileInfo.Size(), nil
}
