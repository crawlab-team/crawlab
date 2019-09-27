package model

import (
	"crawlab/database"
	"crawlab/utils"
	"github.com/apex/log"
	"github.com/globalsign/mgo/bson"
	"os"
	"runtime/debug"
	"time"
)

type GridFs struct {
	Id         bson.ObjectId `json:"_id" bson:"_id"`
	ChunkSize  int32         `json:"chunk_size" bson:"chunkSize"`
	UploadDate time.Time     `json:"upload_date" bson:"uploadDate"`
	Length     int32         `json:"length" bson:"length"`
	Md5        string        `json:"md_5" bson:"md5"`
	Filename   string        `json:"filename" bson:"filename"`
}

type File struct {
	Name  string `json:"name"`
	Path  string `json:"path"`
	IsDir bool   `json:"is_dir"`
	Size  int64  `json:"size"`
}

func (f *GridFs) Remove() {
	s, gf := database.GetGridFs("files")
	defer s.Close()
	if err := gf.RemoveId(f.Id); err != nil {
		log.Errorf("remove file id error: %s, id: %s", err.Error(), f.Id.Hex())
		debug.PrintStack()
	}
}

func GetAllGridFs() []*GridFs {
	s, gf := database.GetGridFs("files")
	defer s.Close()

	var files []*GridFs
	if err := gf.Find(nil).All(&files); err != nil {
		log.Errorf("get all files error: {}", err.Error())
		debug.PrintStack()
		return nil
	}
	return files
}

func GetGridFs(id bson.ObjectId) *GridFs {
	s, gf := database.GetGridFs("files")
	defer s.Close()

	var gfFile GridFs
	err := gf.Find(bson.M{"_id": id}).One(&gfFile)
	if err != nil {
		log.Errorf("get gf file error: %s, file_id: %s", err.Error(), id.Hex())
		debug.PrintStack()
		return nil
	}
	return &gfFile
}

func RemoveFile(path string) error {
	if !utils.Exists(path) {
		log.Info("file not found: " + path)
		debug.PrintStack()
		return nil
	}
	if err := os.Remove(path); err != nil {
		return err
	}
	return nil
}
