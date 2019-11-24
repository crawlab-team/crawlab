package utils

import (
	"archive/zip"
	"bufio"
	"errors"
	"fmt"
	"github.com/apex/log"
	"io"
	"io/ioutil"
	"os"
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
		defer Close(newFile)

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

/**
 * 拷贝文件夹,同时拷贝文件夹中的文件
 * @param srcPath  		需要拷贝的文件夹路径: D:/test
 * @param destPath		拷贝到的位置: D:/backup/
 */
func CopyDir(srcPath string, destPath string) error {
	// 检测目录正确性
	if srcInfo, err := os.Stat(srcPath); err != nil {
		fmt.Println(err.Error())
		return err
	} else {
		if !srcInfo.IsDir() {
			e := errors.New("srcPath不是一个正确的目录！")
			fmt.Println(e.Error())
			return e
		}
	}
	if destInfo, err := os.Stat(destPath); err != nil {
		fmt.Println(err.Error())
		return err
	} else {
		if !destInfo.IsDir() {
			e := errors.New("destInfo不是一个正确的目录！")
			fmt.Println(e.Error())
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
			_, _ = CopyFile(path, destNewPath)
		}
		return nil
	})
	if err != nil {
		fmt.Printf(err.Error())
	}
	return err
}

// 生成目录并拷贝文件
func CopyFile(src, dest string) (w int64, err error) {
	srcFile, err := os.Open(src)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer srcFile.Close()
	// 分割path目录
	destSplitPathDirs := strings.Split(dest, "/")

	// 检测时候存在目录
	destSplitPath := ""
	for index, dir := range destSplitPathDirs {
		if index < len(destSplitPathDirs)-1 {
			destSplitPath = destSplitPath + dir + "/"
			if !Exists(destSplitPath) {
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
	content = strings.ReplaceAll(content, fmt.Sprintf("%s%s%s", sep, key, sep), value)

	// 打开文件
	f, err := os.OpenFile(filePath, os.O_WRONLY, 0777)
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
