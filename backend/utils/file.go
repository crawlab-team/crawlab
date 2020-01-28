package utils

import (
	"archive/zip"
	"bufio"
	"fmt"
	"github.com/apex/log"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"runtime/debug"
	"strings"
)

// 删除文件
func RemoveFiles(path string) {
	if err := os.RemoveAll(path); err != nil {
		log.Errorf("remove files error: %s, path: %s", err.Error(), path)
		debug.PrintStack()
	}
}

// 读取文件一行
func ReadFileOneLine(fileName string) string {
	file := OpenFile(fileName)
	defer Close(file)
	buf := bufio.NewReader(file)
	line, err := buf.ReadString('\n')
	if err != nil {
		log.Errorf("read file error: %s", err.Error())
		return ""
	}
	return line
}

func GetSpiderMd5Str(file string) string {
	md5Str := ReadFileOneLine(file)
	// 去掉空格以及换行符
	md5Str = strings.Replace(md5Str, " ", "", -1)
	md5Str = strings.Replace(md5Str, "\n", "", -1)
	return md5Str
}

// 创建文件
func OpenFile(fileName string) *os.File {
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_RDWR, os.ModePerm)
	if err != nil {
		log.Errorf("create file error: %s, file_name: %s", err.Error(), fileName)
		debug.PrintStack()
		return nil
	}
	return file
}

// 创建文件夹
func CreateFilePath(filePath string) {
	if !Exists(filePath) {
		if err := os.MkdirAll(filePath, os.ModePerm); err != nil {
			log.Errorf("create file error: %s, file_path: %s", err.Error(), filePath)
			debug.PrintStack()
		}
	}
}

// 判断所给路径文件/文件夹是否存在
func Exists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		return os.IsExist(err)
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

func ListDir(path string) []os.FileInfo {
	list, err := ioutil.ReadDir(path)
	if err != nil {
		log.Errorf(err.Error())
		debug.PrintStack()
		return nil
	}
	return list
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
	defer Close(srcFile)
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
	defer Close(zipFile)

	// 遍历zip内所有文件和目录
	for _, innerFile := range zipFile.File {
		// 获取该文件数据
		info := innerFile.FileInfo()

		// 如果是目录，则创建一个
		if info.IsDir() {
			err = os.MkdirAll(filepath.Join(dstPath, innerFile.Name), os.ModeDir|os.ModePerm)
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
			err = os.MkdirAll(filepath.Join(dstPath, dirPath), os.ModeDir|os.ModePerm)
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
		newFile, err := os.OpenFile(filepath.Join(dstPath, innerFile.Name), os.O_RDWR|os.O_CREATE|os.O_TRUNC, info.Mode())
		if err != nil {
			log.Errorf("Unzip File Error : " + err.Error())
			debug.PrintStack()
			continue
		}

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
	defer Close(d)
	w := zip.NewWriter(d)
	defer Close(w)
	for _, file := range files {
		if err := _Compress(file, "", w); err != nil {
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
		if err != nil {
			debug.PrintStack()
			return err
		}
		header.Name = prefix + "/" + header.Name
		writer, err := zw.CreateHeader(header)
		if err != nil {
			debug.PrintStack()
			return err
		}
		_, err = io.Copy(writer, file)
		Close(file)
		if err != nil {
			debug.PrintStack()
			return err
		}
	}
	return nil
}

func GetFilesFromDir(dirPath string) ([]*os.File, error) {
	var res []*os.File
	for _, fInfo := range ListDir(dirPath) {
		f, err := os.Open(filepath.Join(dirPath, fInfo.Name()))
		if err != nil {
			return res, err
		}
		res = append(res, f)
	}
	return res, nil
}

func GetAllFilesFromDir(dirPath string) ([]*os.File, error) {
	var res []*os.File
	if err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if !IsDir(path) {
			f, err2 := os.Open(path)
			if err2 != nil {
				return err
			}
			res = append(res, f)
		}
		return nil
	}); err != nil {
		log.Error(err.Error())
		debug.PrintStack()
		return res, err
	}
	return res, nil
}

// File copies a single file from src to dst
func CopyFile(src, dst string) error {
	var err error
	var srcFd *os.File
	var dstFd *os.File
	var srcInfo os.FileInfo

	if srcFd, err = os.Open(src); err != nil {
		return err
	}
	defer srcFd.Close()

	if dstFd, err = os.Create(dst); err != nil {
		return err
	}
	defer dstFd.Close()

	if _, err = io.Copy(dstFd, srcFd); err != nil {
		return err
	}
	if srcInfo, err = os.Stat(src); err != nil {
		return err
	}
	return os.Chmod(dst, srcInfo.Mode())
}

// Dir copies a whole directory recursively
func CopyDir(src string, dst string) error {
	var err error
	var fds []os.FileInfo
	var srcInfo os.FileInfo

	if srcInfo, err = os.Stat(src); err != nil {
		return err
	}

	if err = os.MkdirAll(dst, srcInfo.Mode()); err != nil {
		return err
	}

	if fds, err = ioutil.ReadDir(src); err != nil {
		return err
	}
	for _, fd := range fds {
		srcfp := path.Join(src, fd.Name())
		dstfp := path.Join(dst, fd.Name())

		if fd.IsDir() {
			if err = CopyDir(srcfp, dstfp); err != nil {
				fmt.Println(err)
			}
		} else {
			if err = CopyFile(srcfp, dstfp); err != nil {
				fmt.Println(err)
			}
		}
	}
	return nil
}

// 设置文件变量值
// 可以理解为将文件中的变量占位符替换为想要设置的值
func SetFileVariable(filePath string, key string, value string) error {
	// 占位符标识
	sep := "###"

	// 读取文件到字节
	contentBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	// 将字节转化为文本
	content := string(contentBytes)

	// 替换文本
	content = strings.Replace(content, fmt.Sprintf("%s%s%s", sep, key, sep), value, -1)

	// 打开文件
	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_TRUNC, 0777)
	if err != nil {
		return err
	}

	// 将替换后的内容写入文件
	if _, err := f.Write([]byte(content)); err != nil {
		return err
	}

	f.Close()

	return nil
}
