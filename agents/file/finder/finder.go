package finder

import (
    "fmt"
    "os"
    "path/filepath"
    "encoding/json"

    "github.com/rookie-xy/hubble/types"

    "github.com/rookie-xy/hubble/job"
    "github.com/rookie-xy/hubble/log"
    "github.com/rookie-xy/hubble/adapter"
    "github.com/rookie-xy/hubble/types/value"

    "github.com/rookie-xy/modules/agents/file/collector"
    "github.com/rookie-xy/modules/agents/file/state"
    "github.com/rookie-xy/modules/agents/file/configure"
    "github.com/rookie-xy/hubble/input"
)

type Finder struct {
    conf    *configure.Configure

    states  *state.States
    jobs    *job.Jobs
    done     chan struct{}

    log      log.Log
}

func New(log log.Log) *Finder {
    return &Finder{
        log:  log,
        jobs: job.New(log),
    }
}

func (f *Finder) Init(conf *configure.Configure, db adapter.SinceDB) error {
	f.conf = conf
    f.states = state.News()

    var states state.States
    if v := db.Get(); v != nil {
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

    return nil
}

func (f *Finder) load(states []state.State) error {
    for _, state := range states {
        source := state.Source
        if f.match(source) {
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
    paths := f.conf.Paths.GetArray()

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

func getFiles(v types.Value) map[string]os.FileInfo {
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

func getPaths(files map[string]os.FileInfo) []string {
    keys := make([]string, 0)
    for file := range files {
        keys = append(keys, file)
    }

    return keys
}

func getState(path string, fi os.FileInfo, f *Finder) (state.State, error) {
    var err error
    var absolutePath string

    absolutePath, err = filepath.Abs(path)
    if err != nil {
        return state.State{}, fmt.Errorf("could not fetch abs path for file %s: %s", absolutePath, err)
    }

    fmt.Printf("Finder check file for collecting: %s\n", absolutePath)

    state := state.New()
    if err := state.Init(fi, absolutePath, "file"); err != nil {
        return state, err
    }

    return state, nil
}

func (f *Finder) Find() {
    files := getFiles(f.conf.Paths)
    paths := getPaths(files)

    for i := 0; i < len(files); i++ {
        path := paths[i]
        info := files[path]

        select {

        case <-f.done:
            fmt.Println("Find aborted because scanner stopped.")
            return
        default:

        }

        new, err := getState(path, info, f)
        if err != nil {
            fmt.Printf("Skipping file %s due to error %s\n", path, err)
        }

        old := f.states.FindPrevious(new)

        if old.IsEmpty() {
            fmt.Printf("Finder start collector for new file: %s\n", new.Source)
            err := f.startCollector(new, 0, f.conf.Input)
            if err != nil {
                fmt.Printf("collector could not be started on new file: %s, Err: %s\n", new.Source, err)
            }

        } else {
            f.collectExistingFile(new, old)
        }
    }

    return
}

func (f *Finder) startCollector(state state.State, offset int64, input input.Input) error {
    if f.conf.Limit > 0 && f.jobs.Len() >= f.conf.Limit {
        return fmt.Errorf("collector limit reached")
    }

    state.Finished = false
    state.Offset = offset

    collector := collector.New(f.log)
    if err := collector.Init(input, f.conf.Output, state); err != nil {
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
	/*
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
	*/

    return false
}
