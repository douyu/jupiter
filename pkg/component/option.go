package component

type options struct {
	priority int64
}

type Option func(*options)

func Priority(priority int64) func(*options) {
	return func(opts *options) { opts.priority = priority }
}
