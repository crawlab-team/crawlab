package log

import (
	"bufio"
	"bytes"
	"errors"
	"github.com/apex/log"
	"github.com/crawlab-team/crawlab/core/utils"
	"github.com/crawlab-team/crawlab/trace"
	"github.com/spf13/viper"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

type FileLogDriver struct {
	// settings
	logFileName string
	rootPath    string

	// internals
	mu sync.Mutex
}

func (d *FileLogDriver) Init() (err error) {
	go d.cleanup()

	return nil
}

func (d *FileLogDriver) Close() (err error) {
	return nil
}

func (d *FileLogDriver) WriteLine(id string, line string) (err error) {
	d.initDir(id)

	d.mu.Lock()
	defer d.mu.Unlock()
	filePath := d.getLogFilePath(id, d.logFileName)

	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, os.FileMode(0760))
	if err != nil {
		return trace.TraceError(err)
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Errorf("close file error: %s", err.Error())
		}
	}(f)

	_, err = f.WriteString(line + "\n")
	if err != nil {
		return trace.TraceError(err)
	}

	return nil
}

func (d *FileLogDriver) WriteLines(id string, lines []string) (err error) {
	linesString := strings.Join(lines, "\n")
	if err := d.WriteLine(id, linesString); err != nil {
		return err
	}
	return nil
}

func (d *FileLogDriver) Find(id string, pattern string, skip int, limit int) (lines []string, err error) {
	if pattern != "" {
		return lines, errors.New("not implemented")
	}
	if !utils.Exists(d.getLogFilePath(id, d.logFileName)) {
		return nil, nil
	}

	f, err := os.Open(d.getLogFilePath(id, d.logFileName))
	if err != nil {
		return nil, trace.TraceError(err)
	}
	defer f.Close()

	sc := bufio.NewReaderSize(f, 1024*1024*10)

	i := -1
	for {
		line, err := sc.ReadString(byte('\n'))
		if err != nil {
			break
		}
		line = strings.TrimSuffix(line, "\n")

		i++

		if i < skip {
			continue
		}

		if i >= skip+limit {
			break
		}

		lines = append(lines, line)
	}

	return lines, nil
}

func (d *FileLogDriver) Count(id string, pattern string) (n int, err error) {
	if pattern != "" {
		return n, errors.New("not implemented")
	}
	if !utils.Exists(d.getLogFilePath(id, d.logFileName)) {
		return 0, nil
	}

	f, err := os.Open(d.getLogFilePath(id, d.logFileName))
	if err != nil {
		return n, trace.TraceError(err)
	}
	return d.lineCounter(f)
}

func (d *FileLogDriver) Flush() (err error) {
	return nil
}

func (d *FileLogDriver) getLogPath() (logPath string) {
	return viper.GetString("log.path")
}

func (d *FileLogDriver) getBasePath(id string) (filePath string) {
	return filepath.Join(d.getLogPath(), id)
}

func (d *FileLogDriver) getMetadataPath(id string) (filePath string) {
	return filepath.Join(d.getBasePath(id), MetadataName)
}

func (d *FileLogDriver) getLogFilePath(id, fileName string) (filePath string) {
	return filepath.Join(d.getBasePath(id), fileName)
}

func (d *FileLogDriver) getLogFiles(id string) (files []os.FileInfo) {
	// 增加了对返回异常的捕获
	files, err := utils.ListDir(d.getBasePath(id))
	if err != nil {
		trace.PrintError(err)
		return nil
	}
	return
}

func (d *FileLogDriver) initDir(id string) {
	if !utils.Exists(d.getBasePath(id)) {
		if err := os.MkdirAll(d.getBasePath(id), os.FileMode(0770)); err != nil {
			trace.PrintError(err)
		}
	}
}

func (d *FileLogDriver) lineCounter(r io.Reader) (n int, err error) {
	buf := make([]byte, 32*1024)
	count := 0
	lineSep := []byte{'\n'}

	for {
		c, err := r.Read(buf)
		count += bytes.Count(buf[:c], lineSep)

		switch {
		case err == io.EOF:
			return count, nil

		case err != nil:
			return count, err
		}
	}
}

func (d *FileLogDriver) getTtl() time.Duration {
	ttl := viper.GetString("log.ttl")
	if ttl == "" {
		return DefaultLogTtl
	}

	if strings.HasSuffix(ttl, "s") {
		ttl = strings.TrimSuffix(ttl, "s")
		n, err := strconv.Atoi(ttl)
		if err != nil {
			return DefaultLogTtl
		}
		return time.Duration(n) * time.Second
	} else if strings.HasSuffix(ttl, "m") {
		ttl = strings.TrimSuffix(ttl, "m")
		n, err := strconv.Atoi(ttl)
		if err != nil {
			return DefaultLogTtl
		}
		return time.Duration(n) * time.Minute
	} else if strings.HasSuffix(ttl, "h") {
		ttl = strings.TrimSuffix(ttl, "h")
		n, err := strconv.Atoi(ttl)
		if err != nil {
			return DefaultLogTtl
		}
		return time.Duration(n) * time.Hour

	} else if strings.HasSuffix(ttl, "d") {
		ttl = strings.TrimSuffix(ttl, "d")
		n, err := strconv.Atoi(ttl)
		if err != nil {
			return DefaultLogTtl
		}
		return time.Duration(n) * 24 * time.Hour
	} else {
		return DefaultLogTtl
	}
}

func (d *FileLogDriver) cleanup() {
	if d.getLogPath() == "" {
		return
	}
	if !utils.Exists(d.getLogPath()) {
		if err := os.MkdirAll(d.getLogPath(), os.FileMode(0770)); err != nil {
			log.Errorf("failed to create log directory: %s", d.getLogPath())
			trace.PrintError(err)
			return
		}
	}
	for {
		// 增加对目录不存在的判断
		dirs, err := utils.ListDir(d.getLogPath())
		if err != nil {
			trace.PrintError(err)
			time.Sleep(10 * time.Minute)
			continue
		}
		for _, dir := range dirs {
			if time.Now().After(dir.ModTime().Add(d.getTtl())) {
				if err := os.RemoveAll(d.getBasePath(dir.Name())); err != nil {
					trace.PrintError(err)
					continue
				}
				log.Infof("removed outdated log directory: %s", d.getBasePath(dir.Name()))
			}
		}

		time.Sleep(10 * time.Minute)
	}
}

var logDriver Driver

func newFileLogDriver() (driver Driver, err error) {
	// driver
	driver = &FileLogDriver{
		logFileName: "log.txt",
		mu:          sync.Mutex{},
	}

	// init
	if err := driver.Init(); err != nil {
		return nil, err
	}

	return driver, nil
}

func GetFileLogDriver() (driver Driver, err error) {
	if logDriver != nil {
		return logDriver, nil
	}
	logDriver, err = newFileLogDriver()
	if err != nil {
		return nil, err
	}
	return logDriver, nil
}
