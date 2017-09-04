package scanner

import (
    "fmt"
    "os"
    "path/filepath"

    "github.com/rookie-xy/hubble/src/types"
    "github.com/rookie-xy/modules/agents/src/log/collector"
)

type Scanner struct {
    paths    types.Value
    exclude  types.Value
}

func New(paths types.Value) *Scanner {
    return &Scanner{
        paths: paths,
    }
}

func (r *Scanner) Run() {
    r.scan()
}

func getKeys(files map[string]os.FileInfo) []string {
    paths := make([]string, 0)
    for file := range files {
        paths = append(paths, file)
    }

    return paths
}

func getFileState(path string, info os.FileInfo, p *Scanner) (file.State, error) {
    var err error
    var absolutePath string

    absolutePath, err = filepath.Abs(path)
    if err != nil {
        return file.State{}, fmt.Errorf("could not fetch abs path for file %s: %s", absolutePath, err)
    }

    fmt.Println("scanner", "Check file for collecting: %s", absolutePath)

    // Create new state for comparison
    newState := file.NewState(info, absolutePath, p.config.Type)

    return newState, nil
}

func (r *Scanner) scan() {
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

            err := collector.New(r, newState, 0)
            if err != nil {
                fmt.Println("collector could not be started on new file: %s, Err: %s", newState.Source, err)
            }

        } else {
            r.collectExistingFile(newState, lastState)
        }
    }
}

// startCollector starts a new collector with the given offset
// In case the CollectorLimit is reached, an error is returned
func (r *Scanner) startCollector(sta state, offset int64) error {
    if r.config.CollectorLimit > 0 && r.collectors.Len() >= r.config.CollectorLimit {
        //collectorSkipped.Add(1)
        return fmt.Errorf("collector limit reached")
    }

    // Set state to "not" finished to indicate that a collector is running
    state.Finished = false
    state.Offset = offset

    // Create collector with state
    //c, err := r.createCollector(state)
    c, err := collector.New(state)
    if err != nil {
        return err
    }

    err = r.Setup()
    if err != nil {
        return fmt.Errorf("Error setting up collector: %s", err)
    }

    // Update state before staring collector
    // This makes sure the states is set to Finished: false
    // This is synchronous state update as part of the scan
    c.SendStateUpdate()

    r.collectors.Start(c)

    return nil
}

func (r *Scanner) Stop() {

}

func (r *Scanner) Wait() {

}

func (r *Scanner) getFiles() map[string]os.FileInfo {
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
func (r *Scanner) isFileExcluded(file string) bool {
    patterns := r.exclude
    return len(patterns) > 0 && collector.MatchAny(patterns, file)
}
