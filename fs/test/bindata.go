// Code generated for package test by go-bindata DO NOT EDIT. (@generated)
// sources:
// bin/start.sh
// bin/stop.sh
package test

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func bindataRead(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	clErr := gz.Close()

	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}
	if clErr != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type asset struct {
	bytes []byte
	info  os.FileInfo
}

type bindataFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
}

// Name return file name
func (fi bindataFileInfo) Name() string {
	return fi.name
}

// Size return file size
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}

// Mode return file mode
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}

// Mode return file modify time
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}

// IsDir return file whether a directory
func (fi bindataFileInfo) IsDir() bool {
	return fi.mode&os.ModeDir != 0
}

// Sys return file is sys mode
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

var _binStartSh = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xcc\x8f\x4b\x0e\xc2\x30\x0c\x44\xf7\x3e\xc5\x20\xd6\x4d\x58\xc3\x51\x80\x45\x4b\x1c\xd5\x22\x69\xab\x38\x2d\x3d\x3e\x6a\xa0\x7c\xc4\x05\x90\x57\x7e\xa3\x19\x8f\xb7\x1b\xdb\x48\x67\xb5\x25\xf1\x38\xa2\x62\x18\x9b\xe3\x80\xf3\x01\xb9\xe5\x8e\x80\x3d\x71\x50\x26\x20\x5e\x9d\xa4\x87\x4c\x5e\xe8\x69\x98\x61\x47\x4d\x36\xf4\x97\x3a\x94\xa8\x1b\xb3\xfb\xb0\x97\x55\x39\x4d\x9c\x70\x22\x00\xa8\x5e\x31\x2b\x88\xb5\x66\x4e\xe6\x87\x4f\x7d\x18\x23\x2f\xdc\x88\x9b\xbf\x35\x19\x50\x6e\xb6\xbd\xe6\x37\x33\x8d\x74\x0e\x3b\x53\x66\xc5\x5e\x02\xa7\xf5\x0b\x63\x95\xeb\xa5\x94\x57\xfb\x37\xdd\xbc\xd0\x3d\x00\x00\xff\xff\x09\xef\x85\x28\x89\x01\x00\x00")

func binStartShBytes() ([]byte, error) {
	return bindataRead(
		_binStartSh,
		"bin/start.sh",
	)
}

func binStartSh() (*asset, error) {
	bytes, err := binStartShBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "bin/start.sh", size: 393, mode: os.FileMode(493), modTime: time.Unix(1621153670, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _binStopSh = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x52\x56\xd4\x4f\xca\xcc\xd3\x2f\xce\xe0\xe2\xca\xce\xcc\xc9\x51\xd0\xb5\x54\x48\x28\x28\x56\x48\xac\x28\xad\x49\x2f\x4a\x2d\x50\x28\x4f\x4d\x4d\x81\xb0\x74\xcb\x14\x40\x74\x4d\x62\x79\xb6\x82\x7a\x75\x41\x51\x66\x5e\x89\x42\x8c\x8a\x51\xad\x7a\x4d\x45\x62\x51\x7a\x71\x02\x20\x00\x00\xff\xff\x6b\x5f\x02\x48\x4a\x00\x00\x00")

func binStopShBytes() ([]byte, error) {
	return bindataRead(
		_binStopSh,
		"bin/stop.sh",
	)
}

func binStopSh() (*asset, error) {
	bytes, err := binStopShBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "bin/stop.sh", size: 74, mode: os.FileMode(493), modTime: time.Unix(1621154552, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// MustAsset is like Asset but panics when Asset would return an error.
// It simplifies safe initialization of global variables.
func MustAsset(name string) []byte {
	a, err := Asset(name)
	if err != nil {
		panic("asset: Asset(" + name + "): " + err.Error())
	}

	return a
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
}

// AssetNames returns the names of the assets.
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() (*asset, error){
	"bin/start.sh": binStartSh,
	"bin/stop.sh":  binStopSh,
}

// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"}
// AssetDir("data/img") would return []string{"a.png", "b.png"}
// AssetDir("foo.txt") and AssetDir("notexist") would return an error
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		cannonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(cannonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, fmt.Errorf("Asset %s not found", name)
			}
		}
	}
	if node.Func != nil {
		return nil, fmt.Errorf("Asset %s not found", name)
	}
	rv := make([]string, 0, len(node.Children))
	for childName := range node.Children {
		rv = append(rv, childName)
	}
	return rv, nil
}

type bintree struct {
	Func     func() (*asset, error)
	Children map[string]*bintree
}

var _bintree = &bintree{nil, map[string]*bintree{
	"bin": &bintree{nil, map[string]*bintree{
		"start.sh": &bintree{binStartSh, map[string]*bintree{}},
		"stop.sh":  &bintree{binStopSh, map[string]*bintree{}},
	}},
}}

// RestoreAsset restores an asset under the given directory
func RestoreAsset(dir, name string) error {
	data, err := Asset(name)
	if err != nil {
		return err
	}
	info, err := AssetInfo(name)
	if err != nil {
		return err
	}
	err = os.MkdirAll(_filePath(dir, filepath.Dir(name)), os.FileMode(0755))
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(_filePath(dir, name), data, info.Mode())
	if err != nil {
		return err
	}
	err = os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
	if err != nil {
		return err
	}
	return nil
}

// RestoreAssets restores an asset under the given directory recursively
func RestoreAssets(dir, name string) error {
	children, err := AssetDir(name)
	// File
	if err != nil {
		return RestoreAsset(dir, name)
	}
	// Dir
	for _, child := range children {
		err = RestoreAssets(dir, filepath.Join(name, child))
		if err != nil {
			return err
		}
	}
	return nil
}

func _filePath(dir, name string) string {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	return filepath.Join(append([]string{dir}, strings.Split(cannonicalName, "/")...)...)
}
