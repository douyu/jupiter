package constant

type Module int

const (
	ModuleInvalid = Module(iota)

	ModuleClientGrpc
	ModuleClientResty
	ModuleClientRedisStub
	ModuleClientRedisCluster
	ModuleClientRocketMQ
	ModuleClientEtcd

	ModuleRegistryEtcd

	ModuleStoreMongoDB
	ModuleStoreGorm
)
