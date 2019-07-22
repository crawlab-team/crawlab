package utils

import (
	"archive/zip"
	"github.com/apex/log"
	"io"
	"os"
	"path/filepath"
	"runtime/debug"
)

// 判断所给路径文件/文件夹是否存在
func Exists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

// 判断所给路径是否为文件夹
func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

// 判断所给路径是否为文件
func IsFile(path string) bool {
	return !IsDir(path)
}

/**
@tarFile：压缩文件路径
@dest：解压文件夹
*/
func DeCompressByPath(tarFile, dest string) error {
	srcFile, err := os.Open(tarFile)
	if err != nil {
		return err
	}
	defer srcFile.Close()
	return DeCompress(srcFile, dest)
}

/**
@zipFile：压缩文件
@dstPath：解压之后文件保存路径
*/
func DeCompress(srcFile *os.File, dstPath string) error {
	// 如果保存路径不存在，创建一个
	if !Exists(dstPath) {
		if err := os.MkdirAll(dstPath, os.ModePerm); err != nil {
			debug.PrintStack()
			return err
		}
	}

	// 读取zip文件
	zipFile, err := zip.OpenReader(srcFile.Name())
	if err != nil {
		log.Errorf("Unzip File Error：" + err.Error())
		debug.PrintStack()
		return err
	}
	defer zipFile.Close()

	// 遍历zip内所有文件和目录
	for _, innerFile := range zipFile.File {
		// 获取该文件数据
		info := innerFile.FileInfo()

		// 如果是目录，则创建一个
		if info.IsDir() {
			err = os.MkdirAll(filepath.Join(dstPath, innerFile.Name), os.ModePerm)
			if err != nil {
				log.Errorf("Unzip File Error : " + err.Error())
				debug.PrintStack()
				return err
			}
			continue
		}

		// 如果文件目录不存在，则创建一个
		dirPath := filepath.Dir(innerFile.Name)
		if !Exists(dirPath) {
			err = os.MkdirAll(filepath.Join(dstPath, dirPath), os.ModePerm)
			if err != nil {
				log.Errorf("Unzip File Error : " + err.Error())
				debug.PrintStack()
				return err
			}
		}

		// 打开该文件
		srcFile, err := innerFile.Open()
		if err != nil {
			log.Errorf("Unzip File Error : " + err.Error())
			debug.PrintStack()
			continue
		}

		// 创建新文件
		newFile, err := os.Create(filepath.Join(dstPath, innerFile.Name))
		if err != nil {
			log.Errorf("Unzip File Error : " + err.Error())
			debug.PrintStack()
			continue
		}
		defer newFile.Close()

		// 拷贝该文件到新文件中
		if _, err := io.Copy(newFile, srcFile); err != nil {
			debug.PrintStack()
			return err
		}

		// 关闭该文件
		if err := srcFile.Close(); err != nil {
			debug.PrintStack()
			return err
		}

		// 关闭新文件
		if err := newFile.Close(); err != nil {
			debug.PrintStack()
			return err
		}
	}
	return nil
}

//压缩文件
//files 文件数组，可以是不同dir下的文件或者文件夹
//dest 压缩文件存放地址
func Compress(files []*os.File, dest string) error {
	d, _ := os.Create(dest)
	defer d.Close()
	w := zip.NewWriter(d)
	defer w.Close()
	for _, file := range files {
		err := _Compress(file, "", w)
		if err != nil {
			return err
		}
	}
	return nil
}

func _Compress(file *os.File, prefix string, zw *zip.Writer) error {
	info, err := file.Stat()
	if err != nil {
		debug.PrintStack()
		return err
	}
	if info.IsDir() {
		prefix = prefix + "/" + info.Name()
		fileInfos, err := file.Readdir(-1)
		if err != nil {
			debug.PrintStack()
			return err
		}
		for _, fi := range fileInfos {
			f, err := os.Open(file.Name() + "/" + fi.Name())
			if err != nil {
				debug.PrintStack()
				return err
			}
			err = _Compress(f, prefix, zw)
			if err != nil {
				debug.PrintStack()
				return err
			}
		}
	} else {
		header, err := zip.FileInfoHeader(info)
		header.Name = prefix + "/" + header.Name
		if err != nil {
			debug.PrintStack()
			return err
		}
		writer, err := zw.CreateHeader(header)
		if err != nil {
			debug.PrintStack()
			return err
		}
		_, err = io.Copy(writer, file)
		file.Close()
		if err != nil {
			debug.PrintStack()
			return err
		}
	}
	return nil
}
