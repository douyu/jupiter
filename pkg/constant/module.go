package constant

type Module int

const (
	ModuleInvalid = Module(iota)

	ModuleClientGrpc
	ModuleClientResty
	ModuleClientRedis
	ModuleClientRocketMQ
	ModuleClientEtcd

	ModuleStoreMongoDB
	ModuleStoreGorm
)
