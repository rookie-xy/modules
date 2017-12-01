package utils

import (
	"os"
	"path/filepath"
	"fmt"

	"github.com/rookie-xy/hubble/types"
    "github.com/rookie-xy/hubble/models/file"
    "github.com/rookie-xy/hubble/log"
  . "github.com/rookie-xy/hubble/log/level"
)

func GetFiles(v types.Value, log log.Factory) map[string]os.FileInfo {
    files := map[string]os.FileInfo{}
    paths := v.GetArray()

    for _, path := range paths {
        matches, err := filepath.Glob(path.(string))
        if err != nil {
            log(ERROR,"glob(%s) failed: %v", path, err)
            continue
        }

//    OUTER:
        // Check any matched files to see if we need to start a collector
        for _, file := range matches {
             // check if the source is in the exclude_files list
            /*
            if r.isExcluded(source) {
                fmt.Printf("Finder Exclude source: %s\n", source)
                continue
            }
            */

            // Fetch Lstat File info to detected also symlinks
            fileInfo, err := os.Lstat(file)
            if err != nil {
                log(ERROR,"Finder lstat(%s) failed: %s", file, err)
                continue
            }

            if fileInfo.IsDir() {
                log(WARN,"Finder Skipping directory: %s", file)
                continue
            }

            isSymlink := fileInfo.Mode() & os.ModeSymlink > 0
            if isSymlink {
                log(WARN,"Finder ile %s skipped as it is a symlink", file)
                continue
            }

            // Fetch Stat source info which fetches the inode.
								    // In case of a symlink, the original inode is fetched
            fileInfo, err = os.Stat(file)
            if err != nil {
                log(ERROR,"Finder stat(%s) failed: %s", file, err)
                continue
            }

            // If symlink is enabled, it is checked that original is not part of same scanner
            // It original is harvested by other scanner, states will potentially overwrite each other
            /*
                for _, finfo := range paths {
                    if id.SameFile(finfo, fileInfo) {
                        fmt.Println("Same source found as symlink and originap. Skipping source: %s", source)
																				    continue OUTER
                    }
                }
                */

            files[file] = fileInfo
        }
    }

    return files
}

func GetPaths(files map[string]os.FileInfo) []string {
    keys := make([]string, 0)
    for file := range files {
        keys = append(keys, file)
    }

    return keys
}

func GetState(path string, fi os.FileInfo, log log.Factory) (file.State, error) {
    var err error
    var absolutePath string

    absolutePath, err = filepath.Abs(path)
    if err != nil {
        return file.State{}, fmt.Errorf("could not fetch abs path for source %s: %s", absolutePath, err)
    }

    log(DEBUG,"Finder check source for collecting: %s", absolutePath)

    state := file.New()
    if err := state.Init(absolutePath, fi); err != nil {
        return state, err
    }

    return state, nil
}
