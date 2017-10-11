package open

import (
	"os"
)

type File struct {
	*os.File
}

func (File) Continuable() bool { return true }

func ReadOpen(path string) (*os.File, error) {
    flag := os.O_RDONLY
    perm := os.FileMode(0)
    return os.OpenFile(path, flag, perm)
}
