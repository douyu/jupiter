go:
	@echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>making build<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<"
	./script/gobuild.sh hello 12345 bin/hello

version:export APP_REGION=wuhan
version:export APP_ZONE=f3
version:
	./bin/hello --version

run:export APP_REGION=wuhan
run:export APP_ZONE=f3
run:
	./bin/hello --config=config.toml

governor:
	curl http://127.0.0.1:19090/routes
