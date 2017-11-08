package source

import "os"

func Open(path string) (*os.File, error) {
    flag := os.O_RDONLY
    perm := os.FileMode(0)
    return os.OpenFile(path, flag, perm)
}
