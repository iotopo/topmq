package utils

import (
	"os"
	"os/exec"
	"path/filepath"
)

//goland:noinspection GoUnusedExportedFunction
func PathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	return !os.IsNotExist(err)
}

func ChangeWorkDir() {
	file, _ := exec.LookPath(os.Args[0]) // [1]
	_ = os.Chdir(filepath.Dir(file))     // [2] 更新当前工作目录(darwin/linux)
}