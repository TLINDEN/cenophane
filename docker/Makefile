upd = ../upd/upd
version   = $(shell egrep "= .v" ../upd/cfg/config.go | cut -d'=' -f2 | cut -d'"' -f 2)
BRANCH    = $(shell git branch --show-current)
COMMIT    = $(shell git rev-parse --short=8 HEAD)
BUILD     = $(shell date +%Y.%m.%d.%H%M%S) 
VERSION  := $(if $(filter $(BRANCH), development),$(version)-$(BRANCH)-$(COMMIT)-$(BUILD),$(version))

v:
	@echo $(VERSION)

build:
	cp -a ../updlocal upd
	make -C upd clean
	docker-compose build --build-arg VERSION=$(VERSION)
