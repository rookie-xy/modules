package id

import (
    "fmt"
    "os"
    "syscall"
)

type ID struct {
    Inode  uint64 `json:"inode,"`
    Device uint64 `json:"device,"`
}

// GetID returns the file id for non windows systems
func GetID(info os.FileInfo) ID {
    stat := info.Sys().(*syscall.Stat_t)

    // Convert inode and dev to uint64 to be cross platform compatible
    id := ID{
        Inode:  uint64(stat.Ino),
        Device: uint64(stat.Dev),
    }

    return id
}

// IsSame source checks if the files are identical
func (id ID) IsSame(state ID) bool {
    return id.Inode == state.Inode && id.Device == state.Device
}

func (id ID) String() string {
    return fmt.Sprintf("%d-%d", id.Inode, id.Device)
}
