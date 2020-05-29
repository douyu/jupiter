########################################################
grpc-direct-server: ## Run server
	@cd grpc/direct/direct-server &&  go run main.go --config=config.toml

grpc-direct-client: ## Run client
	@cd grpc/direct/direct-client && go run main.go --config=config.toml

config-filewatch: ## Run client
	@cd config/filewatch && go run main.go --config=config.toml
