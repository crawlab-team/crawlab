package test

import (
	"github.com/apex/log"
	"github.com/cenkalti/backoff/v4"
	"github.com/crawlab-team/crawlab/trace"
	"github.com/google/uuid"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"time"
)

func init() {
	var err error
	TmpDir, err = filepath.Abs("tmp")
	if err != nil {
		panic(err)
	}
	if _, err := os.Stat(TmpDir); err != nil {
		if err := os.MkdirAll(TmpDir, os.ModePerm); err != nil {
			panic(err)
		}
	}
	//TmpDir = getTmpDir()
}

var TmpDir string

func StartTestSeaweedFs() (err error) {
	// skip if CRAWLAB_IGNORE_WEED is set true
	if os.Getenv("CRAWLAB_IGNORE_WEED") != "" {
		return nil
	}

	// write to start.sh and stop.sh
	if err := writeShFiles(TmpDir); err != nil {
		return trace.TraceError(err)
	}

	// run weed
	go runCmd(exec.Command("sh", "./start.sh"), TmpDir)

	// wait for containers to be ready
	time.Sleep(5 * time.Second)
	f := func() error {
		_, err := T.m.ListDir("/", true)
		if err != nil {
			return err
		}
		return nil
	}
	b := backoff.WithMaxRetries(backoff.NewConstantBackOff(5*time.Second), 5)
	nt := func(err error, duration time.Duration) {
		log.Infof("seaweedfs services not ready, re-attempt in %.1f seconds", duration.Seconds())
	}
	err = backoff.RetryNotify(f, b, nt)
	if err != nil {
		return trace.TraceError(err)
	}

	return nil
}

func StopTestSeaweedFs() (err error) {
	// skip if CRAWLAB_IGNORE_WEED is set true
	if os.Getenv("CRAWLAB_IGNORE_WEED") != "" {
		return nil
	}

	// stop seaweedfs
	if err := runCmd(exec.Command("sh", "./stop.sh"), TmpDir); err != nil {
		return trace.TraceError(err)
	}
	time.Sleep(5 * time.Second)

	// remove tmp folder
	if err := os.RemoveAll(TmpDir); err != nil {
		return trace.TraceError(err)
	}

	return nil
}

func writeShFiles(dirPath string) (err error) {
	fileNames := []string{
		"start.sh",
		"stop.sh",
	}

	for _, fileName := range fileNames {
		data, err := Asset("bin/" + fileName)
		if err != nil {
			return trace.TraceError(err)
		}
		filePath := path.Join(dirPath, fileName)
		if err := ioutil.WriteFile(filePath, data, os.FileMode(0766)); err != nil {
			return trace.TraceError(err)
		}
	}

	return nil
}

func runCmd(cmd *exec.Cmd, dirPath string) (err error) {
	log.Infof("running cmd: %v", cmd)
	cmd.Dir = dirPath
	//cmd.Stdout = os.Stdout
	//cmd.Stderr = os.Stdout
	return cmd.Run()
}

func getTmpDir() string {
	id, _ := uuid.NewUUID()
	tmpDir := path.Join(os.TempDir(), id.String())
	if _, err := os.Stat(tmpDir); err != nil {
		if err := os.MkdirAll(tmpDir, os.FileMode(0766)); err != nil {
			panic(err)
		}
	}
	return tmpDir
}
