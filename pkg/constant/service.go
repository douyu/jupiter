package constant

type ServiceKind uint8

const (
	ServiceUnknown ServiceKind = iota
	ServiceProvider
	ServiceGovernor
	ServiceConsumer
)

func (sk ServiceKind) String() string {
	switch sk {
	case ServiceProvider:
		return "providers"
	case ServiceGovernor:
		return "governors"
	case ServiceConsumer:
		return "consumers"
	default:
		return "unknown"
	}
}
