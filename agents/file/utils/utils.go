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

func GetState(state file.State) file.State {
	// refreshes the values in State with the values from the collector itself
    state.ID = file.Id(state.Fileinfo)
    return state
}

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

func CreateState(path string, fi os.FileInfo, log log.Factory) (file.State, error) {
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

func Match(file string, paths types.Value, logf log.Factory) bool {
    file = filepath.Clean(file)
    for _, path := range paths.GetArray() {
        path := filepath.Clean(path.(string))

        match, err := filepath.Match(path, file)
        if err != nil {
            logf(ERROR, "Finder error matching glob: %s", err)
            continue
        }

        if match {
            return true
        }
    }

    return false
}

func Excluded(file string) bool {
    return false
}

func Open(path string) (*os.File, error) {
    flag := os.O_RDONLY
    perm := os.FileMode(0)
    return os.OpenFile(path, flag, perm)
}

/*
old := f.states.FindPrevious(new)
        if old.IsEmpty() {
            //err := f.startCollector(new, 0)
            err := collector.Start(new, 0, f.input, f.decoder, f.jobs, f.states, f.log, f.conf)
            if err != nil {
                f.logf(ERROR,"Collector could not be started on new source: %s, Err: %s", new.Source, err)
            }
        } else {
            //f.keepCollector(new, old)
            collector.Keep(new, old, f.logf)
        }
*/

