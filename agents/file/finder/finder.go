package finder

import (
    "fmt"

    "github.com/rookie-xy/hubble/job"
    "github.com/rookie-xy/hubble/log"
    "github.com/rookie-xy/hubble/adapter"

    "github.com/rookie-xy/modules/agents/file/collector"
    "github.com/rookie-xy/modules/agents/file/configure"
    "github.com/rookie-xy/hubble/input"
    "github.com/rookie-xy/hubble/models/file"
	"github.com/rookie-xy/modules/agents/file/utils"
	"github.com/rookie-xy/hubble/codec"
	"github.com/rookie-xy/hubble/prototype"
	"sync"
  . "github.com/rookie-xy/hubble/log/level"
	"errors"
)

type Finder struct {
    conf    *configure.Configure

    states  *file.States
    jobs    *job.Jobs
    done     chan struct{}
    once     sync.Once

    decoder  codec.Decoder
    input    input.Input
    level    Level
    log      log.Log
    logf     log.Factory
}

func New(log log.Log) *Finder {
    return &Finder{
        log:  log,
        level: adapter.ToLevelLog(log).Get(),
        jobs: job.New(log),
        done: make(chan struct{}),
    }
}

func (f *Finder) Init(input input.Input, decoder codec.Decoder,
	                  conf *configure.Configure, sinceDB adapter.SinceDB, logf log.Factory) error {
	f.logf = logf
	f.conf = conf
    f.states = file.News(f.logf)

    f.decoder = decoder
    f.input   = input

    if states := sinceDB.Load(); states != nil {
        if err := f.load(states); err != nil {
            return err
        }
    }

    return nil
}

func (f *Finder) load(states []file.State) error {
    for _, state := range states {
        if utils.Match(state.Source, f.conf.Paths, f.logf) {
            state.TTL = -1

            if !state.Finished {
                return fmt.Errorf("can only start a finder when all related "+
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

func (f *Finder) Find() {
    files := utils.GetFiles(f.conf.Paths, f.logf)
    paths := utils.GetPaths(files)

    for i := 0; i < len(files); i++ {
        path := paths[i]
        info := files[path]

        select {

        case <-f.done:
            f.logf(INFO,"Find aborted because scanner stopped.")
            return
        default:
        }

        new, err := utils.CreateState(path, info, f.logf)
        if err != nil {
            f.logf(ERROR,"Skipping source %s due to error %s", path, err)
            continue
        }

        old := f.states.FindPrevious(new)
        if old.IsEmpty() {
            err := f.startCollector(new, 0)
            if err != nil {
                f.logf(ERROR,"Collector could not be started on new source: %s, Err: %s", new.Source, err)
            }
        } else {
            f.keepCollector(new, old)
        }
    }
}

func (f *Finder) startCollector(state file.State, offset int64) error {
    if f.conf.Limit > 0 && f.jobs.Len() >= f.conf.Limit {
        return errors.New("collector limit reached")
    }

    state.Finished = false
    state.Offset = offset

    input   := prototype.Input(f.input)
    decoder := prototype.Decoder(f.decoder)

    collector := collector.New(f.log)
    if err := collector.Init(input, decoder, state, f.states, f.conf); err != nil {
        return err
    }

    collector.Update(state)
    f.jobs.Start(collector)
    return nil
}

func (f *Finder) keepCollector(new, old file.State) {
 	f.logf(DEBUG,"Finder Update existing file for collecting: %s, offset: %v, finish:%v, newFileSize:%d\n",
 		            new.Source, old.Offset, old.Finished, new.Fileinfo.Size())

	if old.Finished && new.Fileinfo.Size() > old.Offset {
		f.logf(DEBUG,"Finder Resuming collecting of file: %s, offset: %d, new size: %d\n",
			            new.Source, old.Offset, new.Fileinfo.Size())
		err := f.startCollector(new, old.Offset)
		if err != nil {
            f.logf(ERROR,"Collector could not be started on existing file: %s, Err: %s\n",
            	            new.Source, err)
		}
		return
	}

	if old.Finished && new.Fileinfo.Size() < old.Offset {
		f.logf(DEBUG,"Finder old file was truncated. Starting from the beginning: %s, offset: %d, new size: %d\n",
			            new.Source, new.Fileinfo.Size())

		err := f.startCollector(new, 0)
		if err != nil {
			f.logf(ERROR,"Collector could not be started on truncated file: %s, Err: %s\n",
				            new.Source, err)
		}
		return
	}

	if old.Source != "" && old.Source != new.Source {
		f.logf(DEBUG,"Finder file rename was detected: %s -> %s, Current offset: %v\n",
			            old.Source, new.Source, old.Offset)

		if old.Finished {
			f.logf(DEBUG,"Finder updating state for renamed file: %s -> %s, Current offset: %v\n",
				            old.Source, new.Source, old.Offset)

			old.Source = new.Source
			err := f.update(old)
			if err != nil {
				f.logf(ERROR,"File rotation state update error: %s\n", err)
			}

		} else {
			f.logf(WARN,"Finder file rename detected but collector not finished yet.")
		}
	}

	if !old.Finished {
		f.logf(DEBUG,"Finder collector for file is still running: %s\n", new.Source)
	} else {
		f.logf(DEBUG,"Finder file didn't change: %s\n", new.Source)
    }
}

func (f *Finder) Wait() {
	f.jobs.WaitForCompletion()
}

func (f *Finder) Stop() {
    close(f.done)

	if length := f.jobs.Len(); length > 0 {
		f.jobs.Stop()
	}
}
