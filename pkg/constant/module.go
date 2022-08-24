package constant

type Module int

const (
	ModuleInvalid = Module(iota)

	ModuleClientGrpc
	ModuleClientResty
	ModuleClientRedis
	ModuleClientRocketMQ
	ModuleClientETCDV3

	ModuleStoreMongoDB
)
