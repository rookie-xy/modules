package file

import (
    "fmt"
    "os"
    "syscall"
)

type Key struct {
    Inode  uint64 `json:"inode,"`
    Device uint64 `json:"device,"`
}

// GetOSState returns the FileStateOS for non windows systems
func GetOSState(info os.FileInfo) Key {
    stat := info.Sys().(*syscall.Stat_t)

    // Convert inode and dev to uint64 to be cross platform compatible
    fileState := Key{
        Inode:  uint64(stat.Ino),
        Device: uint64(stat.Dev),
    }

    return fileState
}

// IsSame file checks if the files are identical
func (r Key) IsSame(state Key) bool {
    return r.Inode == state.Inode && r.Device == state.Device
}

func (r Key) String() string {
    return fmt.Sprintf("%d-%d", r.Inode, r.Device)
}

// ReadOpen opens a file for reading only
func ReadOpen(path string) (*os.File, error) {
    flag := os.O_RDONLY
    perm := os.FileMode(0)

    return os.OpenFile(path, flag, perm)
}
