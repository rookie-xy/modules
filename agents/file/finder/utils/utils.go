package utils

import (
	"os"
	"path/filepath"
	"fmt"

	"github.com/rookie-xy/modules/agents/file/id"
	"github.com/rookie-xy/hubble/types"
    "github.com/rookie-xy/hubble/models/file"
)

func GetFiles(v types.Value) map[string]os.FileInfo {
    files := map[string]os.FileInfo{}
    paths := v.GetArray()

    for _, path := range paths {
        matches, err := filepath.Glob(path.(string))
        if err != nil {
            fmt.Printf("glob(%s) failed: %v\n", path, err)
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
                fmt.Printf("Finder lstat(%s) failed: %s\n", file, err)
                continue
            }

            if fileInfo.IsDir() {
                fmt.Printf("Finder Skipping directory: %s\n", file)
                continue
            }

            isSymlink := fileInfo.Mode() & os.ModeSymlink > 0
            if isSymlink {
                fmt.Printf("Finder ile %s skipped as it is a symlink.\n", file)
                continue
            }

            // Fetch Stat source info which fetches the inode.
								    // In case of a symlink, the original inode is fetched
            fileInfo, err = os.Stat(file)
            if err != nil {
                fmt.Printf("Finder stat(%s) failed: %s\n", file, err)
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

func GetState(path string, fi os.FileInfo) (file.State, error) {
    var err error
    var absolutePath string

    absolutePath, err = filepath.Abs(path)
    if err != nil {
        return file.State{}, fmt.Errorf("could not fetch abs path for source %s: %s", absolutePath, err)
    }

    fmt.Printf("Finder check source for collecting: %s\n", absolutePath)

    state := file.New()
    if err := state.Init(id.GetID(fi).String(), fi, absolutePath, "file"); err != nil {
        return state, err
    }

    return state, nil
}

//import "github.com/rookie-xy/modules/agents/log/match"

// MatchAny checks if the text matches any of the regular expressions
/*
func MatchAny(matchers []match.Matcher, source string) bool {
    for _, m := range matchers {
        if m.MatchString(source) {
            return true
        }
    }

    return false
}
*/
