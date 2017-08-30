package api

type Tracker interface {
   	Start()
   	Stop()
   	Wait()
}
/*
func New(l log.Log) (*Collector, error) {
    c := &Collector{
				    Log: l,
    }

    c.Init()

    return c, nil
}

func (r *Collector) Init() {
    return
}

func (r *Collector) Start() {
    return
}
*/
