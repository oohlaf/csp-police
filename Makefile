# Author: Olaf Conradi <olaf@conradi.org>
#
# Type make help to see a list of targets.

all: build  ##- Build everything.

build: webd grpcd  ##- Build all programs.

# Source files
SRC_CMD := $(wildcard cmd/**/*.go)
SRC_CMD_GRPCD := $(wildcard cmd/grpcd/*.go)
SRC_CMD_WEBD := $(wildcard cmd/webd/*.go)
SRC_INT := $(wildcard internal/**/*.go)
SRC_PKG := $(wildcard pkg/**/*.go)

SRC_PROTO := $(wildcard api/proto/**/*.proto)
API_OUTDIR := pkg/api
API_OUT := $(subst api/proto/,$(API_OUTDIR)/,$(patsubst %.proto,%.pb.go,$(SRC_PROTO)) $(patsubst %.proto,%_grpc.pb.go,$(SRC_PROTO)))

SOURCES := $(SRC_INT) $(SRC_PKG)

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

# Go recipe to build a command program
GO_BUILD_CMD = $(GO) build -o $@ $(firstword $^)

## Target recipes to build each program

.PHONY: grpcd
grpcd: $(BIN_GRPCD)  ##- Build grpcd program.

$(BIN_GRPCD): cmd/grpcd/main.go $(SRC_CMD_GRPCD) $(SOURCES) $(API_OUT) | bin
	$(GO_BUILD_CMD)

.PHONY: webd
webd: $(BIN_WEBD)  ##- Build webd program.

$(BIN_WEBD): cmd/webd/main.go $(SRC_CMD_WEBD) $(SOURCES) $(API_OUT) | bin
	$(GO_BUILD_CMD)

run-grpcd: $(BIN_GRPCD)  ##- Run grpcd program.
	$(BIN_GRPCD)

run-webd: $(BIN_WEBD)  ##- Run webd program.
	$(BIN_WEBD)

## Generic targets

.PHONY: all api build clean help

api: $(API_OUT)  ##- Generate protobuf GRPC API files.

$(API_OUT): $(SRC_PROTO) | $(API_OUTDIR)
	$(PROTOC) -I api/proto --go_out=$(API_OUTDIR) --go_opt=paths=source_relative --go-grpc_out=$(API_OUTDIR) --go-grpc_opt=paths=source_relative $(SRC_PROTO)

bin:
	$(MKDIR) $@

$(API_OUTDIR):
	$(MKDIR) -p $@

clean:  ##- Remove build artifacts.
	$(RM) $(API_OUT)
	$(RMDIR) bin

help:  ##- Show this help.
	@$(MAKE_HELP) $(MAKEFILE_LIST)
