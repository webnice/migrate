DIR=$(strip $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST)))))

GOPATH := $(DIR):$(GOPATH)
DATE=$(shell date -u +%Y%m%d.%H%M%S.%Z)
TESTPACKETS=$(shell if [ -f .testpackages ]; then cat .testpackages; fi)
BENCHPACKETS=$(shell if [ -f .benchpackages ]; then cat .benchpackages; fi)

default: lint test

link:
	@echo "prepare..."
	@mkdir {src,bin} 2>/dev/null; true
	@if [ ! -L $(DIR)/src/goose ]; then ln -s $(DIR)/goose $(DIR)/src/goose 2>/dev/null; fi
	@if [ ! -L $(DIR)/src/gsmigrate ]; then ln -s $(DIR)/gsmigrate $(DIR)/src/gsmigrate 2>/dev/null; fi
	@if [ ! -L $(DIR)/src/vendor ]; then ln -s $(DIR)/vendor $(DIR)/src/vendor 2>/dev/null; fi
	@if [ ! -L $(DIR)/src/github.com ]; then ln -s $(DIR)/vendor/github.com $(DIR)/src/github.com 2>/dev/null; fi
	@if [ ! -L $(DIR)/src/golang.org ]; then ln -s $(DIR)/vendor/golang.org $(DIR)/src/golang.org 2>/dev/null; fi
	@if [ ! -L $(DIR)/src/google.golang.org ]; then ln -s $(DIR)/vendor/google.golang.org $(DIR)/src/google.golang.org 2>/dev/null; fi
	@if [ ! -L $(DIR)/src/alecthomas ]; then ln -s $(DIR)/vendor/gopkg.in/alecthomas $(DIR)/src/alecthomas 2>/dev/null; fi
	@cd ${DIR}/src && ln -s . gopkg.in 2>/dev/null; true
	@cd ${DIR}/src && ln -s . webnice 2>/dev/null; true
	@cd ${DIR}/src && ln -s . migrate.v1 2>/dev/null; true
.PHONY: link

test: link
	@echo "mode: set" > $(DIR)/coverage.log
	@for PACKET in $(TESTPACKETS); do \
		touch $(DIR)/coverage-tmp.log; \
		GOPATH=${GOPATH} go test -v -covermode=count -coverprofile=$(DIR)/coverage-tmp.log $$PACKET; \
		if [ "$$?" -ne "0" ]; then exit $$?; fi; \
		tail -n +2 $(DIR)/coverage-tmp.log | sort -r | awk '{if($$1 != last) {print $$0;last=$$1}}' >> $(DIR)/coverage.log; \
		rm -f $(DIR)/coverage-tmp.log; true; \
	done
.PHONY: test

cover: test
	GOPATH=${GOPATH} go tool cover -html=$(DIR)/coverage.log
	@make clean
.PHONY: cover

bench: link
	@for PACKET in $(BENCHPACKETS); do GOPATH=${GOPATH} go test -race -bench=. -benchmem $$PACKET; done
	@make clean
.PHONY: bench

lint: link
	GOPATH=${GOPATH} gometalinter \
	--vendor \
	--deadline=15m \
	--cyclo-over=20 \
	--disable=aligncheck \
	--linter="vet:go tool vet -printf {path}/*.go:PATH:LINE:MESSAGE" \
	--skip=src/vendor \
	--skip=github.com/mattn/go-sqlite3 \
	./...
	@make clean
.PHONY: lint

clean:
	@echo "cleaning..."
	@rm -rf ${DIR}/src; true
	@rm -rf ${DIR}/bin; true
	@rm -rf ${DIR}/pkg; true
	@rm -rf ${DIR}/*.log; true
	@rm -rf ${DIR}/*.db; true
.PHONY: clean
