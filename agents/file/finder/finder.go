package finder

import (
    "fmt"
    "os"
    "errors"
    "regexp"
    "path/filepath"
    "encoding/json"

    "github.com/rookie-xy/hubble/types"

    "github.com/rookie-xy/hubble/job"
    "github.com/rookie-xy/hubble/log"
    "github.com/rookie-xy/hubble/adapter"
    "github.com/rookie-xy/hubble/factory"
    "github.com/rookie-xy/hubble/types/value"

    "github.com/rookie-xy/modules/agents/file/collector"
    "github.com/rookie-xy/modules/agents/file/state"
)

type Finder struct {
    log         log.Log

    paths       types.Value
    excludes    types.Value
    from        string
    states     *state.States
    sincedb     adapter.SinceDB
    done        chan struct{}

    jobs       *job.Jobs
    limit       uint64
}

func New(log log.Log) *Finder {
    return &Finder{
        log:  log,
        jobs: job.New(log),
    }
}

func (f *Finder) Init(from string,
                      paths, excludes types.Value, limit uint64) error {
    f.paths    = paths
    f.excludes = excludes
    f.from = from
    f.states = state.News()

    if client, err := factory.Forward("plugin.client.sincedb"); err != nil {
        return err
    } else {
        f.sincedb = adapter.ToSinceDB(client)
    }

    var states state.States
    if v := f.sincedb.Get(); v != nil {
        if val := value.New(v); val != nil {
            if err := json.Unmarshal(val.GetBytes(), &states); err != nil {
            	fmt.Println(err)
                return err
            }
        }
    }

    if err := f.load(states.States); err != nil {
        return err
    }

    if limit > 0 {
        f.limit = limit
    } else {
        return errors.New("limit is error")
    }

    return nil
}

func (f *Finder) load(states []state.State) error {
    for _, state := range states {
        if f.match(state.Source) {
        	state.TTL = -1

        	if !state.Finished {
        	    return fmt.Errorf("Can only start a finder when all related " +
        	                             "states are finished: %+v", state)
            }

            if err := f.update(state); err != nil {
                return err
            }
        }
    }

    return nil
}

func (f *Finder) update(state state.State) error {
	f.states.Update(state)
    return nil
}

func (f *Finder) match(file string) bool {
    file = filepath.Clean(file)
    paths := f.paths.GetArray()

    for _, path := range paths {
        path := filepath.Clean(path.(string))

        match, err := filepath.Match(path, file)
        if err != nil {
            fmt.Printf("finder", "Error matching glob: %s", err)
            continue
        }

        if match {
            return true
        }
    }

    return false
}

func getFiles(path types.Value) map[string]os.FileInfo {
    files := map[string]os.FileInfo{}
    paths := path.GetArray()

    for _, path := range paths {
        matches, err := filepath.Glob(path.(string))
        if err != nil {
            fmt.Printf("glob(%s) failed: %v\n", path, err)
            continue
        }

//    OUTER:
        // Check any matched files to see if we need to start a collector
        for _, file := range matches {
             // check if the file is in the exclude_files list
            /*
            if r.isExcluded(file) {
                fmt.Printf("Finder Exclude file: %s\n", file)
                continue
            }
            */

            // Fetch Lstat File info to detected also symlinks
            fileInfo, err := os.Lstat(file)
            if err != nil {
                fmt.Println("scanner", "lstat(%s) failed: %s", file, err)
                continue
            }

            if fileInfo.IsDir() {
                fmt.Println("scanner", "Skipping directory: %s", file)
                continue
            }

            isSymlink := fileInfo.Mode() & os.ModeSymlink > 0
            if isSymlink {
                fmt.Println("scanner", "File %s skipped as it is a symlink.", file)
                continue
            }

            // Fetch Stat file info which fetches the inode.
								    // In case of a symlink, the original inode is fetched
            fileInfo, err = os.Stat(file)
            if err != nil {
                fmt.Println("scanner", "stat(%s) failed: %s", file, err)
                continue
            }

            // If symlink is enabled, it is checked that original is not part of same scanner
            // It original is harvested by other scanner, states will potentially overwrite each other
            /*
                for _, finfo := range paths {
                    if os.SameFile(finfo, fileInfo) {
                        fmt.Println("Same file found as symlink and originap. Skipping file: %s", file)
																				    continue OUTER
                    }
                }
                */

            files[file] = fileInfo
        }
    }

    return files
}

func getKeys(files map[string]os.FileInfo) []string {
    paths := make([]string, 0)
    for file := range files {
        paths = append(paths, file)
    }

    return paths
}

func getFileState(path string, fi os.FileInfo, s *Finder) (state.State, error) {
    var err error
    var absolutePath string

    absolutePath, err = filepath.Abs(path)
    if err != nil {
        return state.State{}, fmt.Errorf("could not fetch abs path for file %s: %s", absolutePath, err)
    }

    fmt.Println("finder", "Check file for collecting: %s", absolutePath)

    newState := state.New(fi, absolutePath, s.from)

    return newState, nil
}

func (r *Finder) Find() {
    var paths []string

    files := getFiles(r.paths)
    paths = getKeys(files)

    for i := 0; i < len(files); i++ {
        var path string
        var info os.FileInfo

        path = paths[i]
        info = files[path]
        fmt.Println(info.Name())

        select {

        case <-r.done:
            fmt.Println("Find aborted because scanner stopped.")
            return
        default:

        }

        newState, err := getFileState(path, info, r)
        if err != nil {
            fmt.Println("Skipping file %s due to error %s", path, err)
        }

        oldState := r.states.FindPrevious(newState)

        // Decides if previous state exists
        if oldState.IsEmpty() {
            fmt.Println("finder", "Start collector for new file: %s", newState.Source)
            err := r.startCollector(newState, 0)
            if err != nil {
                fmt.Println("collector could not be started on new file: %s, Err: %s", newState.Source, err)
            }

        } else {
            r.collectExistingFile(newState, oldState)
        }
    }

    return
}

func (f *Finder) startCollector(state state.State, offset int64) error {
    if f.limit > 0 && f.jobs.Len() >= f.limit {
        return fmt.Errorf("collector limit reached")
    }

    state.Finished = false
    state.Offset = offset

    collector := collector.New(f.log)
    if err := collector.Init(group.GetString(), Type.GetString(),
                                                codec, client); err != nil {
        return err
    }

    if err := collector.Setup(); err != nil {
        return err
    }

    collector.Update(state)

    f.jobs.Start(collector)

    return nil
}

func (r *Finder) collectExistingFile(newState, oldState state.State) {
}

func (r *Finder) Stop() {

}

func (r *Finder) Wait() {

}

func (r *Finder) isExcluded(file string) bool {
    patterns := r.excludes.GetArray()
    if len(patterns) > 0 {
        for _, pattern := range patterns {
            if matched, err := regexp.MatchString(pattern.(string), file); err != nil {
                fmt.Println(err)
            } else {
                if matched {
                    return matched
                }
            }
        }
    }

    return false
}
