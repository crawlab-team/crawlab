package utils

import (
	"archive/zip"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/apex/log"
	"github.com/crawlab-team/crawlab/core/constants"
	"github.com/crawlab-team/crawlab/core/entity"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"runtime/debug"
)

func OpenFile(fileName string) *os.File {
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_RDWR, os.ModePerm)
	if err != nil {
		log.Errorf("create file error: %s, file_name: %s", err.Error(), fileName)
		debug.PrintStack()
		return nil
	}
	return file
}

func Exists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		return os.IsExist(err)
	}
	return true
}

func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

// ListDir Add: 增加error类型作为第二返回值
// 在其他函数如 /task/log/file_driver.go中的 *FileLogDriver.cleanup()函数调用时
// 可以通过判断err是否为nil来判断是否有错误发生
func ListDir(path string) ([]fs.FileInfo, error) {
	list, err := os.ReadDir(path)
	if err != nil {
		log.Errorf(err.Error())
		debug.PrintStack()
		return nil, err
	}

	var res []fs.FileInfo
	for _, item := range list {
		info, err := item.Info()
		if err != nil {
			log.Errorf(err.Error())
			debug.PrintStack()
			return nil, err
		}
		res = append(res, info)
	}
	return res, nil
}

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
		dirPath := filepath.Join(dstPath, filepath.Dir(innerFile.Name))
		if !Exists(dirPath) {
			if err = os.MkdirAll(dirPath, os.ModeDir|os.ModePerm); err != nil {
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
		newFilePath := filepath.Join(dstPath, innerFile.Name)
		newFile, err := os.OpenFile(newFilePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, info.Mode())
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

// Compress 压缩文件
// files 文件数组，可以是不同dir下的文件或者文件夹
// dest 压缩文件存放地址
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

func TrimFileData(data []byte) (res []byte) {
	if string(data) == constants.EmptyFileData {
		return res
	}
	return data
}

func ZipDirectory(dir, zipfile string) error {
	zipFile, err := os.Create(zipfile)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	baseDir := filepath.Dir(dir)

	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		relPath, err := filepath.Rel(baseDir, path)
		if err != nil {
			return err
		}

		zipFile, err := zipWriter.Create(relPath)
		if err != nil {
			return err
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		_, err = io.Copy(zipFile, file)
		if err != nil {
			return err
		}

		return nil
	})

	return err
}

// CopyFile File copies a single file from src to dst
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

// CopyDir Dir copies a whole directory recursively
func CopyDir(src string, dst string) error {
	var err error
	var fds []os.DirEntry
	var srcInfo os.FileInfo

	if srcInfo, err = os.Stat(src); err != nil {
		return err
	}

	if err = os.MkdirAll(dst, srcInfo.Mode()); err != nil {
		return err
	}

	if fds, err = os.ReadDir(src); err != nil {
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

func GetFileHash(filePath string) (res string, err error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

func ScanDirectory(dir string) (res map[string]entity.FsFileInfo, err error) {
	files := make(map[string]entity.FsFileInfo)

	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		var hash string
		if !info.IsDir() {
			hash, err = GetFileHash(path)
			if err != nil {
				return err
			}
		}

		relPath, err := filepath.Rel(dir, path)
		if err != nil {
			return err
		}

		files[relPath] = entity.FsFileInfo{
			Name:      info.Name(),
			Path:      relPath,
			FullPath:  path,
			Extension: filepath.Ext(path),
			IsDir:     info.IsDir(),
			FileSize:  info.Size(),
			ModTime:   info.ModTime(),
			Mode:      info.Mode(),
			Hash:      hash,
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return files, nil
}
