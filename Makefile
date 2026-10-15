GO      ?= go
BINDIR  := bin
PREFIX  ?= /usr/local
BINDEST := $(DESTDIR)$(PREFIX)/bin

# One binary per cmd/*/ directory.
CMDS := $(notdir $(wildcard cmd/*))
BINS := $(addprefix $(BINDIR)/,$(CMDS))

# All binaries share the osctl library sources, so a library change rebuilds
# everything.
SRC := $(wildcard internal/osctl/*.go)

.PHONY: all test vet clean install uninstall

all: $(BINS)

$(BINDIR)/%: cmd/%/main.go $(SRC) go.mod
	@mkdir -p $(BINDIR)
	$(GO) build -o $@ ./cmd/$*

test:
	$(GO) test ./...

vet:
	$(GO) vet ./...

install: all
	install -d $(BINDEST)
	install -m 0755 $(BINS) $(BINDEST)/

uninstall:
	rm -f $(addprefix $(BINDEST)/,$(CMDS))

clean:
	rm -rf $(BINDIR)