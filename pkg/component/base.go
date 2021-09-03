package component

type BaseComponent struct{}

func (c BaseComponent) ShouldBeLeader() bool {
	return false
}

func (c BaseComponent) baseMethod() {}

func (c BaseComponent) Start(stop <-chan struct{}) error {
	panic("not implemented")
}
