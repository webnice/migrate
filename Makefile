DIR=$(strip $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST)))))

GOPATH := $(GOPATH)
GOBIN  := $(DIR)/bin
DATE=$(shell date -u +%Y%m%d.%H%M%S.%Z)
TESTPACKETS=$(shell if [ -f .testpackages ]; then cat .testpackages; fi)
BENCHPACKETS=$(shell if [ -f .benchpackages ]; then cat .benchpackages; fi)

PRJ01=gsmigrate
BIN01=$(DIR)/bin/$(PRJ01)
VER01="1.0.0-build.$(DATE)"
VERN01=$(shell echo "$(VER01)" | awk -F '-' '{ print $$1 }' )
VERB01=$(shell echo "$(VER01)" | awk -F 'build.' '{ print $$2 }' )

default: lint test

dep:
	@mkdir -p ${DIR}/{bin,pkg} 2>/dev/null; true
	@go get -u ./...
	@go mod download
	@go mod tidy
	@go mod vendor
.PHONY: dep

build:
	@cd ${DIR}/gsmigrate; go build -o ${BIN01}
.PHONY: build

rpm:
	mkdir -p ${DIR}/rpmbuild/{BUILD,BUILDROOT,RPMS,SOURCES,SPECS,SRPMS}; true
	cp -v ${DIR}/conf/${PRJ01}.spec ${DIR}/rpmbuild/SPECS/${PRJ01}.spec
	cp -v ${DIR}/bin/${PRJ01} ${DIR}/rpmbuild/SOURCES/${PRJ01}
	rpmbuild \
		--define "_topdir $(DIR)/rpmbuild" \
		--define "_app_version_number $(VERN01)" \
		--define "_app_version_build $(VERB01)" \
		-bb "$(DIR)/rpmbuild/SPECS/$(PRJ01).spec"
.PHONY: rpm

test:
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

bench:
	@for PACKET in $(BENCHPACKETS); do GOPATH=${GOPATH} go test -race -bench=. -benchmem $$PACKET; done
	@make clean
.PHONY: bench

lint:
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
	@rm -rf ${DIR}/rpmbuild; true
	@rm -rf ${DIR}/*.log; true
	@rm -rf ${DIR}/*.db; true
.PHONY: clean
