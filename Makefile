# Author: Olaf Conradi <olaf@conradi.org>
#
# Type make help to see a list of targets.

all: build  ##- Build everything.

build: webd grpcd  ##- Build all programs.

# Source files
SRC_CMD := $(wildcard cmd/**/*.go)
SRC_INT := $(wildcard internal/**/*.go)
SRC_PKG := $(wildcard pkg/**/*.go)

SRC_PROTO := $(wildcard api/proto/**/*.proto)
API_OUT := $(subst api/proto/,pkg/api/,$(patsubst %.proto,%.pb.go,$(SRC_PROTO)) $(patsubst %.proto,%_grpc.pb.go,$(SRC_PROTO)))

SOURCES := $(SRC_CMD) $(SRC_INT) $(SRC_PKG)

ifeq ($(OS),Windows_NT)
BIN_WEBD := bin/webd.exe
BIN_GRPCD := bin/grpcd.exe
else
BIN_WEBD := bin/webd
BIN_GRPCD := bin/grpcd
endif

# OS independent binaries
GO ?= go
PROTOC ?= protoc
MKDIR ?= mkdir

# OS specific binaries
ifeq ($(OS),Windows_NT)
RM := PowerShell -ExecutionPolicy ByPass -File ./scripts/make-helper.ps1 rm
RMDIR := PowerShell -ExecutionPolicy ByPass -File ./scripts/make-helper.ps1 rmdir
MAKE_HELP := PowerShell -ExecutionPolicy ByPass -File ./scripts/make-helper.ps1 help
else
RM ?= rm -f
RMDIR ?= rm -f -r
MAKE_HELP := awk 'BEGIN {FS=":.*\#{2}-[ \t]"; print "Usage: make [TARGET]\n\nTargets:"}; /^[^\#].*:.*\#{2}-.*$$/ {printf "  %-28s%s\n", $$1, $$2}'
endif

## Target recipes to build each program

.PHONY: grpcd
grpcd: $(BIN_GRPCD)  ##- Build grpcd program.

$(BIN_GRPCD): cmd/grpcd/main.go bin $(SOURCES) $(API_OUT)
	$(go-build-cmd)

.PHONY: webd
webd: $(BIN_WEBD)  ##- Build webd program.

$(BIN_WEBD): cmd/webd/main.go bin $(SOURCES) $(API_OUT)
	$(go-build-cmd)

run-grpcd: $(BIN_GRPCD)  ##- Run grpcd program.
	$(BIN_GRPCD)

run-webd: $(BIN_WEBD)  ##- Run webd program.
	$(BIN_WEBD)

## Generic targets

.PHONY: all api build clean help

api: $(API_OUT)  ##- Generate protobuf GRPC API files.

$(API_OUT): $(SRC_PROTO)
	$(PROTOC) -I api/proto --go_out=pkg/api --go_opt=paths=source_relative --go-grpc_out=pkg/api --go-grpc_opt=paths=source_relative $(SRC_PROTO)

# Go recipe to build a command program
define go-build-cmd =
	$(GO) build -o $@ $(firstword $^)
endef

bin:
	$(MKDIR) bin

clean:  ##- Remove build artifacts.
	$(RM) $(API_OUT)
	$(RMDIR) bin

help:  ##- Show this help.
	@$(MAKE_HELP) $(MAKEFILE_LIST)
