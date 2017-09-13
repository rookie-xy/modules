package finder

import (
    "fmt"
    "os"
    "errors"
    "path/filepath"

    "github.com/rookie-xy/hubble/types"

    "github.com/rookie-xy/hubble/job"
    "github.com/rookie-xy/hubble/log"

    "github.com/rookie-xy/modules/agents/log/collector"
    "github.com/rookie-xy/modules/agents/log/file/state"
)

type Finder struct {
    log        log.Log
    jobs      *job.Jobs
    paths      types.Value
    excludes   types.Value
    from       string
    states    *state.States
    done       chan struct{}

    collector *collector.Collector
    limit      uint64
}

func New(log log.Log) *Finder {
    return &Finder{
        log:  log,
        jobs: job.New(log),
    }
}

func (r *Finder) Init(from string,
                       paths, excludes types.Value,
                       cc *collector.Collector, limit uint64) error {
    r.paths    = paths
    r.excludes = excludes
    r.from = from
    r.collector = cc

    if limit > 0 {
        r.limit = limit
    } else {
        return errors.New("limit is error")
    }

    return nil
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

    fmt.Println("scanner", "Check file for collecting: %s", absolutePath)

    // Create new state for comparison
    newState := state.New(fi, absolutePath, s.from)

    return newState, nil
}

func (r *Finder) Find() bool {
    var paths []string

    files := r.getFiles()
		  paths = getKeys(files)

    for i := 0; i < len(files); i++ {

        var path string
        var info os.FileInfo

        path = paths[i]
        info = files[path]

        select {

        case <-r.done:
            fmt.Println("Scan aborted because scanner stopped.")
            return

        default:

        }

        newState, err := getFileState(path, info, r)
        if err != nil {
            fmt.Println("Skipping file %s due to error %s", path, err)
        }

        // Load last state
        lastState := r.states.FindPrevious(newState)

        // Ignores all files which fall under ignore_older
        if r.isIgnoreOlder(newState) {
            err := r.handleIgnoreOlder(lastState, newState)
            if err != nil {
                fmt.Println("Updating ignore_older state error: %s", err)
												}

            continue
        }

        // Decides if previous state exists
        if lastState.IsEmpty() {

            fmt.Println("scanner", "Start collector for new file: %s", newState.Source)
            err := r.startTextFinder(newState, 0)
            if err != nil {
                fmt.Println("collector could not be started on new file: %s, Err: %s", newState.Source, err)
            }

        } else {
            r.collectExistingFile(newState, lastState)
        }
    }

    return true
}

// startCollector starts a new collector with the given offset
// In case the CollectorLimit is reached, an error is returned
func (r *Finder) startCollector(state state.State, offset int64) error {
    if r.limit > 0 && r.jobs.Len() >= r.limit {
        //collectorSkipped.Add(1)
        return fmt.Errorf("collector limit reached")
    }

    // Set state to "not" finished to indicate that a collector is running
    state.Finished = false
    state.Offset = offset

    // Create collector with state
    job := r.collector.Job(r.states)

    // Update state before staring collector
    // This makes sure the states is set to Finished: false
    // This is synchronous state update as part of the scan
    job.Update(state)

    r.jobs.Start(job)

    return nil
}

func (r *Finder) Stop() {

}

func (r *Finder) Wait() {

}

func (r *Finder) getFiles() map[string]os.FileInfo {
    files := map[string]os.FileInfo{}
    paths := r.paths.GetArray()

    for _, path := range paths {

        matches, err := filepath.Glob(path.(string))
        if err != nil {
            fmt.Println("glob(%s) failed: %v", path, err)
            continue
        }

    //OUTER:
        // Check any matched files to see if we need to start a collector
        for _, file := range matches {

             // check if the file is in the exclude_files list
             if r.isFileExcluded(file) {
                 fmt.Println("scanner", "Exclude file: %s", file)
                 continue
             }

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
            if p.config.Symlinks {
                for _, finfo := range paths {
                    if os.SameFile(finfo, fileInfo) {
                        fmt.Println("Same file found as symlink and originap. Skipping file: %s", file)
																				    continue OUTER
                    }
                }
            }
            */

            files[file] = fileInfo
        }
    }

    return files
}

// isFileExcluded checks if the given path should be excluded
func (r *Finder) isFileExcluded(file string) bool {
    patterns := r.excludes
    return len(patterns) > 0 && MatchAny(patterns, file)
}
