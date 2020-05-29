# Jupiter Golang Application Standard Makefile
SHELL:=/bin/bash
BASE_PATH:=$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
BUILD_PATH:=$(BASE_PATH)/build
TITLE:=$(shell basename $(BASE_PATH))
BUILD_TIME:=$(shell date +%Y-%m-%d--%T)
JUPITER:=/vendor/github.com/dsouyu/jupiter
LDFLAGS:="-X $(JUPITER).buildTime=$(BUILD_TIME) -X $(JUPITER).name=$(APP_NAME) -X $(JUPITER).id=$(APP_ID)"

all:print fmt lint buildDemo
alltar:print fmt lint buildDemo

print:
	@echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>making print<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<"
	@echo SHELL:$(SHELL)
	@echo BASE_PATH:$(BASE_PATH)
	@echo BUILD_PATH:$(BUILD_PATH)
	@echo TITLE:$(TITLE)
	@echo BUILD_TIME:$(BUILD_TIME)
	@echo JUPITER:$(JUPITER)
	@echo BINS:$(BINS)
	@echo APP_NAME:$(APP_NAME)
	@echo LDFLAGS:$(LDFLAGS)
	@echo -e "\n"

fmt:
	@echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>making fmt<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<"
	go fmt $(TITLE)/internal/...
	@echo -e "\n"

lint:
	@echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>making lint<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<"
ifndef REVIVE
	go get -u github.com/mgechev/revive
endif
	revive -formatter stylish internal/...
	@echo -e "\n"

errcheck:
	@echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>making errcheck<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<"
ifndef ERRCHCEK
	go get -u github.com/kisielk/errcheck
endif
	@errcheck internal/...
	@echo -e "\n"

test:
	@echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>making test<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<"
	@echo testPath ${BAST_PATH}
	go test -v .${BAST_PATH}/...

buildDemo:
	@echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>making build<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<"
	chmod +x $(BUILD_PATH)/script/shell/*.sh
	$(BUILD_PATH)/script/shell/build.sh $(LDFLAGS) $(BINS)
	@echo -e "\n"

run:
	go run cmd/main.go --config=config/config.toml


