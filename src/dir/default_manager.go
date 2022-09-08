package dir

import (
	"fmt"
	"path/filepath"
)

type DefaultManager struct {
	Path, FileName, Context string
}

func (dm *DefaultManager) Filepath() string {
	return filepath.Join(dm.Abs(), dm.FileName)
}

func (dm *DefaultManager) Abs() string {
	return filepath.Join(dm.Path, dm.createDir())
}

func (dm *DefaultManager) createDir() string {
	return fmt.Sprintf("%s/%s/%s/%s/", dm.Context, dm.FileName[0:2], dm.FileName[2:4], dm.FileName[4:6])
}
