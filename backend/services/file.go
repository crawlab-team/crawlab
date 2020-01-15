package services

import (
	"crawlab/model"
	"github.com/apex/log"
	"os"
	"path"
	"runtime/debug"
)

func GetFileNodeTree(dstPath string, level int) (f model.File, err error) {
	dstF, err := os.Open(dstPath)
	if err != nil {
		log.Errorf(err.Error())
		debug.PrintStack()
		return f, err
	}
	defer dstF.Close()
	fileInfo, err := dstF.Stat()
	if err != nil {
		log.Errorf(err.Error())
		debug.PrintStack()
		return f, nil
	}
	if !fileInfo.IsDir() { //如果dstF是文件
		return model.File{
			Label:     fileInfo.Name(),
			Name:     fileInfo.Name(),
			Path:     dstPath,
			IsDir:    false,
			Size:     fileInfo.Size(),
			Children: nil,
		}, nil
	} else { //如果dstF是文件夹
		dir, err := dstF.Readdir(0) //获取文件夹下各个文件或文件夹的fileInfo
		if err != nil {
			log.Errorf(err.Error())
			debug.PrintStack()
			return f, nil
		}
		f = model.File{
			Label:    path.Base(dstPath),
			Name:     path.Base(dstPath),
			Path:     dstPath,
			IsDir:    true,
			Size:     0,
			Children: nil,
		}
		for _, subFileInfo := range dir {
			subFileNode, err := GetFileNodeTree(path.Join(dstPath, subFileInfo.Name()), level+1)
			if err != nil {
				log.Errorf(err.Error())
				debug.PrintStack()
				return f, err
			}
			f.Children = append(f.Children, subFileNode)
		}
		return f, nil
	}
}
