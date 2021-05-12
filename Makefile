REPO = github.com/lixinio/weixin
BINARIES:=wxwork_agent wxwork_oauth

all: $(BINARIES)
build: $(BINARIES)

$(BINARIES):
	mkdir -p build
	echo "$@"
	find ./examples/ -mindepth 1 -maxdepth 1 -type d -name "$@" | \
		grep "$@" |\
		xargs go build -mod=vendor -o build

.PHONY: mod
mod:
	go mod tidy && go mod vendor
