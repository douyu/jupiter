package constant

//ServiceKind service kind
type ServiceKind uint8

const (
	//ServiceUnknown service non-name
	ServiceUnknown ServiceKind = iota
	//ServiceProvider service provider
	ServiceProvider
	//ServiceGovernor service governor
	ServiceGovernor
	//ServiceConsumer service consumer
	ServiceConsumer
)

var serviceKinds = make(map[ServiceKind]string)

func init() {
	serviceKinds[ServiceUnknown] = "unknown"
	serviceKinds[ServiceProvider] = "providers"
	serviceKinds[ServiceGovernor] = "governors"
	serviceKinds[ServiceConsumer] = "consumers"
}

func (sk ServiceKind) String() string {
	if s, ok := serviceKinds[sk]; ok {
		return s
	}
	return "unknown"
}
