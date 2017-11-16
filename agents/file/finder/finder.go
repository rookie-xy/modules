package finder

import (
    "fmt"
    "path/filepath"


    "github.com/rookie-xy/hubble/job"
    "github.com/rookie-xy/hubble/log"
    "github.com/rookie-xy/hubble/adapter"

    "github.com/rookie-xy/modules/agents/file/collector"
    "github.com/rookie-xy/modules/agents/file/configure"
    "github.com/rookie-xy/hubble/input"
    "github.com/rookie-xy/hubble/models/file"
	"github.com/rookie-xy/modules/agents/file/finder/utils"
	"github.com/rookie-xy/hubble/codec"
	"github.com/rookie-xy/hubble/prototype"
)

type Finder struct {
    conf    *configure.Configure

    states  *file.States
    jobs    *job.Jobs
    done     chan struct{}

    codec    codec.Codec
    input    input.Input

    log      log.Log
}

func New(log log.Log) *Finder {
    return &Finder{
        log:  log,
        jobs: job.New(log),
    }
}

func (f *Finder) Init(input input.Input, codec codec.Codec,
	                  conf *configure.Configure, sinceDB adapter.SinceDB) error {
	f.conf = conf
    f.states = file.News()

    f.codec = codec
    f.input = input

    if states := sinceDB.Load(); states != nil {
        if err := f.load(states); err != nil {
            return err
        }
    }

    return nil
}

func (f *Finder) load(states []file.State) error {
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

func (f *Finder) update(state file.State) error {
    if f.conf.Expire > 0 && state.TTL != 0 {
    	state.TTL = f.conf.Expire
    }

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
            fmt.Printf("Finder error matching glob: %s\n", err)
            continue
        }

        if match {
            return true
        }
    }

    return false
}

func (f *Finder) Find() {
    files := utils.GetFiles(f.conf.Paths)
    paths := utils.GetPaths(files)

    for i := 0; i < len(files); i++ {
        path := paths[i]
        info := files[path]

        select {

        case <-f.done:
            fmt.Println("Find aborted because scanner stopped.")
            return
        default:

        }

        new, err := utils.GetState(path, info)
        if err != nil {
            fmt.Printf("Skipping source %s due to error %s\n", path, err)
        }

        old := f.states.FindPrevious(new)
        if old.IsEmpty() {
            err := f.startCollector(new, 0)
            if err != nil {
                fmt.Printf("Collector could not be started on new source: %s, Err: %s\n", new.Source, err)
            }

        } else {
            f.keepCollector(new, old)
        }
    }

    return
}

func (f *Finder) startCollector(state file.State, offset int64) error {
    if f.conf.Limit > 0 && f.jobs.Len() >= f.conf.Limit {
        return fmt.Errorf("collector limit reached")
    }

    state.Finished = false
    state.Offset = offset

    input := prototype.Input(f.input)
    codec := prototype.Codec(f.codec)

    collector := collector.New(f.log)
    if err := collector.Init(input, codec, state, f.states, f.conf); err != nil {
        return err
    }

    collector.Update(state)

    f.jobs.Start(collector)

    return nil
}

func (f *Finder) keepCollector(new, old file.State) {
 	fmt.Printf("Finder Update existing file for collecting: %s, offset: %v, finish:%v, newFileSize:%d\n", new.Source, old.Offset, old.Finished, new.Fileinfo.Size())

	// No collector is running for the file, start a new collector
	// It is important here that only the size is checked and not modification time, as modification time could be incorrect on windows
	// https://blogs.technet.microsoft.com/asiasupp/2010/12/14/file-date-modified-property-are-not-updating-while-modifying-a-file-without-closing-it/
	if old.Finished && new.Fileinfo.Size() > old.Offset {
		// Resume harvesting of an old file we've stopped harvesting from
		// This could also be an issue with force_close_older that a new collector is started after each scan but not needed?
		// One problem with comparing modTime is that it is in seconds, and scans can happen more then once a second
		fmt.Printf("hhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhFinder Resuming collecting of file: %s, offset: %d, new size: %d\n", new.Source, old.Offset, new.Fileinfo.Size())
		err := f.startCollector(new, old.Offset)
		if err != nil {
            fmt.Printf("Collector could not be started on existing file: %s, Err: %s\n", new.Source, err)
		}
		return
	}

	// File size was reduced -> truncated file
	if old.Finished && new.Fileinfo.Size() < old.Offset {
		fmt.Printf("Finder old file was truncated. Starting from the beginning: %s, offset: %d, new size: %d\n", new.Source, new.Fileinfo.Size())
		err := f.startCollector(new, 0)
		if err != nil {
			fmt.Printf("Collector could not be started on truncated file: %s, Err: %s\n", new.Source, err)
		}
		return
	}

	// Check if file was renamed
	if old.Source != "" && old.Source != new.Source {
		// This does not start a new collector as it is assume that the older collector is still running
		// or no new lines were detected. It sends only an event status update to make sure the new name is persisted.
		fmt.Printf("Finder file rename was detected: %s -> %s, Current offset: %v\n", old.Source, new.Source, old.Offset)

		if old.Finished {
			fmt.Printf("Finder updating state for renamed file: %s -> %s, Current offset: %v\n", old.Source, new.Source, old.Offset)
			// Update state because of file rotation
			old.Source = new.Source
			err := f.update(old)
			if err != nil {
				fmt.Printf("File rotation state update error: %s\n", err)
			}
		} else {
			fmt.Println("Finder file rename detected but collector not finished yet.")
		}
	}

	if !old.Finished {
		// Nothing to do. Collector is still running and file was not renamed
		fmt.Printf("Finder collector for file is still running: %s\n", new.Source)
	} else {
		fmt.Printf("Finder file didn't change: %s\n", new.Source)
    }
}

func (f *Finder) Stop() {

}

func (f *Finder) Wait() {

}


func (f *Finder) isExcluded(file string) bool {
	/*
    patterns := r.excludes.GetArray()
    if len(patterns) > 0 {
        for _, pattern := range patterns {
            if matched, err := regexp.MatchString(pattern.(string), source); err != nil {
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
