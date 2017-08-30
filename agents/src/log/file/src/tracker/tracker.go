package tracker


type tracker struct {

}

func New() *tracker {
    return &tracker{}
}

func (r *tracker) Start() {
    r.scan()
}

func (r *tracker) scan() {

}

func (r *tracker) Stop() {

}

func (r *tracker) Wait() {

}
