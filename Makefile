# This Makefile is meant to be used by people that do not usually work
# with Go source code. If you know what GOPATH is then you probably
# don't need to bother with make.

.PHONY: all bottos bcli wallet clean

GOBIN = build/bin

bottos:
	build/vmlib.sh
	build/bottos.sh
	go build
	@echo "Done building."
	@echo "Run \"./bottos --help\" for help."

bcli:
	build/bcli.sh go install
	@echo "Done building."
	@echo "Run \"$(GOBIN)/bcli\" to launch command line tool."

wallet:
	build/wallet.sh go install
	@echo "Done building."
	@echo "Run \"$(GOBIN)/wallet\" to launch wallet."

all: bottos bcli wallet

clean:
	rm -fr build/_workspace/pkg/ vendor/pkg $(GOBIN)/*
