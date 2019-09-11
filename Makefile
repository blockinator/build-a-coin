export GOPATH := ${PWD}:${GOPATH}

PACKAGE=buildacoin
BINDIR=bin
PKGDIR=pkg
LOCALDIRS=${BINDIR} ${PKGDIR}

.PHONY: compile
compile:
	go build ${PACKAGE}/...

.PHONY: webserver
webserver:
	go install ${PACKAGE}/web/main
	mv ${BINDIR}/main ${BINDIR}/buildacoin-web

.PHONY: tool
tool:
	go install ${PACKAGE}/tool/main
	mv ${BINDIR}/main ${BINDIR}/buildacoin-tool

# use | between deps here to avoid races
.PHONY: bin
bin: webserver
	@echo "done"

.PHONY: run
run: bin
	./bin/buildacoin-web -conf testing/web_test_conf.json

.PHONY: test
test:
	go test ${PACKAGE}/...

.PHONY: bench
bench:
	go test -bench . ${PACKAGE}/...

.PHONY: clean
clean:
	go clean ${PACKAGE}/...
	rm -rf ${LOCALDIRS}

.PHONY: doc
doc:
	go doc ${PACKAGE}/...
