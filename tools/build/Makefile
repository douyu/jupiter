# Jupiter Golang Application Standard Makefile
SHELL:=/bin/bash
BASE_PATH:=$(shell dirname $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST)))))
BUILD_PATH:=$(BASE_PATH)/build
TITLE:=$(shell basename $(BASE_PATH))
BUILD_TIME:=$(shell date +%Y-%m-%d--%T)
JUPITER:=github.com/douyu/jupiter

all:print fmt lint build zip

alltar:print fmt lint build tar

npm:
	@cd static && npm start a

print:
	@echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>making print<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<"
	@echo SHELL:$(SHELL)
	@echo BASE_PATH:$(BASE_PATH)
	@echo BUILD_PATH:$(BUILD_PATH)
	@echo TITLE:$(TITLE)
	@echo BUILD_TIME:$(BUILD_TIME)
	@echo JUPITER:$(JUPITER)
	@echo -e "\n"

fmt:
	@echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>making fmt<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<"
	go fmt $(TITLE)/internal/...
	@echo -e "\n"

lint:
	@echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>making lint<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<"
	@OUT=`$(GOPATH)/bin/golint $(BASE_PATH)/app/... | grep -v pb.go`;\
	if [ "$$OUT" != "" ]; then\
		echo $$OUT;\
		exit 1;\
	else\
		echo "lint success";\
	fi
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

build:
	@echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>making build<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<"
	chmod +x $(BUILD_PATH)/script/shell/*.sh
	$(BUILD_PATH)/script/shell/build.sh $(TITLE) $(JUPITER)/pkg $(BASE_PATH)/cmd/main.go
	@echo -e "\n"

ent:
	go generate ./ent

tar:
	@echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>making tar<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<"
	$(eval TAR=$(BASE_PATH)/$(TITLE)_$(VCS_INFO).tar.gz)
	@cd $(BASE_PATH) && tar zcf $(TAR) bin config build lib Gopkg.* >/dev/null
	@echo $(TAR)
	@echo -e "\n"

zip:
	@echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>making zip<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<"
	@ZIP=`cd $(BASE_PATH) && zip -r $(PUB_ZIP) bin config lib template* Gopkg.* >/dev/null`
	@if [ "$(PUB_ZIP)" == "" ]; then\
		echo "no pub zip";\
	else\
		echo $$ZIP;\
	fi
	@echo $(PUB_ZIP)
	@echo -e "\n"
