PID      = /tmp/go-bumper.pid
SOURCEDIR=.
SOURCES := $(shell find $(SOURCEDIR) -name '*.go')
TEST_SOURCES := $(shell find $(SOURCEDIR) -name '*_test.go')

BINARY=go-bumper

VERSION=1.0.0
BUILD_TIME=`date +%FT%T%z`

LDFLAGS=

.DEFAULT_GOAL: $(BINARY)

$(BINARY): $(SOURCES)
	go build ${LDFLAGS} -o $@ $?

.PHONY: watch
watch: restart
	@fswatch -o $(SOURCES) | xargs -n1 -I{}  make restart || make kill

.PHONY: watchTests
watchTests:
	@fswatch -o $(TEST_SOURCES) | xargs -n1 -I{}  make test

.PHONY: test
	go test .

.PHONY: kill
kill:
	@kill `cat $(PID)` || true

.PHONY: restart
restart: kill $(BINARY)
	./$(BINARY) & echo $$! > $(PID)

.PHONY: install
install:
	go install ${LDFLAGS} ./...

.PHONY: clean
clean:
	if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi
