package scanner

type Scanner struct {

}

func New() *Scanner {
    return &Scanner{}
}

func (r *Scanner) Run() {
    r.scan()
}

func (r *Scanner) scan() {

}

func (r *Scanner) Stop() {

}

func (r *Scanner) Wait() {

}
