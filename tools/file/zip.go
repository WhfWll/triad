package file

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
	"strings"

	ezip "github.com/alexmullins/zip"
)

// Zip 压缩文件夹，path:待压缩的文件/文件夹 路径；zipFilePath:zip文件路径
func Zip(path, zipFilePath string) error {
	//第一步，创建 zip 文件
	zipFile, err := os.Create(zipFilePath)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	//第二步，创建一个新的 *Writer 对象
	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	//打开待压缩的文件/文件夹
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	//关闭 zip writer，将所有数据写入指向基础 zip 文件的数据流
	defer f.Close()

	err = compress(f, "", zipWriter)

	return err
}

// EncryptZip 加密压缩文件
func EncryptZip(src, desc, password string) error {
	zipFile, err := os.Create(desc)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	archive := ezip.NewWriter(zipFile)
	defer archive.Close()

	err = filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		header, err := ezip.FileInfoHeader(info)
		if err != nil {
			return err
		}
		header.Name = strings.TrimPrefix(path, filepath.Dir(src)+"/")
		if info.IsDir() {
			header.Name += "/"
		} else {
			header.Method = zip.Deflate
		}
		// 设置密码
		header.SetPassword(password)
		writer, err := archive.CreateHeader(header)
		if err != nil {
			return err
		}
		if !info.IsDir() {
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()
			_, err = io.Copy(writer, file)
		}
		return err
	})
	return err
}

func compress(file *os.File, prefix string, zw *zip.Writer) error {
	info, err := file.Stat()
	if err != nil {
		return err
	}
	if info.IsDir() {
		prefix = prefix + "/" + info.Name()
		fileInfos, err := file.Readdir(-1)
		if err != nil {
			return err
		}
		for _, fi := range fileInfos {
			f, err := os.Open(file.Name() + "/" + fi.Name())
			if err != nil {
				return err
			}
			err = compress(f, prefix, zw)
			if err != nil {
				return err
			}
		}
	} else {
		header, err := zip.FileInfoHeader(info)
		header.Name = prefix + "/" + header.Name
		if err != nil {
			return err
		}
		writer, err := zw.CreateHeader(header)
		if err != nil {
			return err
		}
		//将原文件内容复制到压缩文件中
		_, err = io.Copy(writer, file)
		file.Close()
		if err != nil {
			return err
		}
	}
	return nil
}
