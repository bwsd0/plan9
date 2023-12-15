GO := go
GOOS := $(shell go env GOOS)
GOARCH := $(shell go env GOARCH)
GO_LDFLAGS := -s
GO_FLAGS := -a -ldflags="$(GO_LDFLAGS) $(GO_EXTLDFLAGS)"
VERSION := 0.0.1

TARG := $(shell go list ./acme/cmd/acme)

ifndef V
  V = 0
endif

ifeq ($(V),1)
  Q =
else
  Q = @
endif

.PHONY: build
build:
	$(Q)mkdir -p bin
	$(Q)GOARCH=$(GOARCH) $(GO) build -o bin/$(shell basename $(TARG)) $(TARG)

.PHONY: release
release:
	$(GO) mod tidy
	$(call cross,arm64)

install: release
	$(Q)GOARCH=$(GOARCH) $(GO) build -o bin/$(shell basename $(TARG)) $(TARG)

.PHONY: cross
cross:
	$(Q)$(GO) mod tidy
	$(Q)$(call crosscompile,linux,amd64)

define crosscompile
$(Q)mkdir -p bin
$(Q)GOOS=$(1) GOARCH=$(2) $(GO) build -mod=readonly \
	-o bin/$(shell basename $(TARG))-$(VERSION)-$(1)-$(2) \
	$(GO_FLAGS) \
	$(TARG);
endef

clean:
	$(Q)$(RM) -r bin/
