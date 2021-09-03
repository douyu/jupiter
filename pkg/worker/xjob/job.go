package job

type JobFunc func()

func (jn JobFunc) Run() {
	jn()
}

// Runner ...
type Runner interface {
	Run()
}

type JobComponent struct {
}
