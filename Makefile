DIR=$(strip $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST)))))

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
	@mkdir -p ${DIR}/bin ${DIR}/pkg 2>/dev/null; true
	@go get -u ./...
	@go mod download
	@go mod tidy
	@go mod vendor
.PHONY: dep

## Встраивание внешней статики в приложение.
embed:
.PHONY: embed

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
	@echo "mode: count" > $(DIR)/coverage.log
	@for PACKET in $(TESTPACKETS); do \
		touch $(DIR)/coverage-tmp.log; \
		go test -v -covermode=count -coverprofile=$(DIR)/coverage-tmp.log $$PACKET || exit 1; \
		tail -n +2 $(DIR)/coverage-tmp.log | sort -r | awk '{if($$1 != last) {print $$0;last=$$1}}' >> $(DIR)/coverage.log; \
		rm -f $(DIR)/coverage-tmp.log; true; \
	done
.PHONY: test

cover: test
	go tool cover -html=$(DIR)/coverage.log
	@make clean
.PHONY: cover

bench:
	@if [ -n "$(BENCHPACKETS)" ]; then \
		for PACKET in $(BENCHPACKETS); do go test -race -bench=. -benchmem $$PACKET; done; \
	fi
	@make clean
.PHONY: bench

lint:
	go vet ./...
.PHONY: lint

clean:
	@echo "cleaning..."
	@rm -rf ${DIR}/bin; true
	@rm -rf ${DIR}/pkg; true
	@rm -rf ${DIR}/rpmbuild; true
	@rm -rf ${DIR}/*.log; true
	@rm -rf ${DIR}/*.db; true
.PHONY: clean
