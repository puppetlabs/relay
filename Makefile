#
# Commands
#

API_SPEC_CONVERTER ?= api-spec-converter
MKDIR ?= mkdir -p
GIT ?= git
GO ?= go
RM ?= rm -f
SHA256SUM ?= shasum -a 256
SWAGGER ?= swagger

#
# Variables
#

NEBULA_API_REPO ?= $(if $(GITHUB_TOKEN),https://$(GITHUB_TOKEN):,ssh://git)@github.com/puppetlabs/nebula-api.git
NEBULA_API_REF ?= master

GOFLAGS ?= -mod=vendor

CLI_DIST_TARGETS ?= linux-amd64 linux-386 linux-arm64 linux-ppc64le linux-s390x windows-amd64 darwin-amd64

#
#
#

CLI_DIST_NAME := nebula-cli
CLI_DIST_VERSION := $(shell $(GIT) describe --tags --always --dirty)

DEPEND_DIR := .depend
ARTIFACTS_DIR := artifacts
BIN_DIR := bin

NEBULA_API_DIR := $(DEPEND_DIR)/nebula-api
NEBULA_API_SPEC_FILENAME := $(NEBULA_API_DIR)/openapi/swagger.yaml

CLI_EXT_linux :=
CLI_EXT_windows := .exe
CLI_EXT_darwin :=

CLI_DIST_PREFIX := $(ARTIFACTS_DIR)/$(CLI_DIST_NAME)-$(CLI_DIST_VERSION)-
CLI_DIST_BINS := $(foreach TARGET,$(CLI_DIST_TARGETS),$(TARGET)$(CLI_EXT_$(word 1,$(subst -, ,$(TARGET)))))
CLI_DIST_BINS := $(addprefix $(CLI_DIST_PREFIX),$(CLI_DIST_BINS))
CLI_DIST_SHA256 := $(addsuffix .sha256,$(CLI_DIST_BINS))

#
# Targets
#

.PHONY: all
all: build

$(DEPEND_DIR) $(ARTIFACTS_DIR) $(BIN_DIR):
	$(MKDIR) $@

ifneq (,$(NEBULA_API_REPO))
$(NEBULA_API_DIR)/.git:
	$(GIT) clone --depth 1 --branch $(NEBULA_API_REF) $(NEBULA_API_REPO) $(NEBULA_API_DIR)

$(NEBULA_API_SPEC_FILENAME): $(NEBULA_API_DIR)/.git

$(DEPEND_DIR)/swagger.yaml: $(NEBULA_API_SPEC_FILENAME) $(DEPEND_DIR)
	$(API_SPEC_CONVERTER) -f openapi_3 -t swagger_2 -s yaml $^ >$@

pkg/client/api/nebula_client.go: $(DEPEND_DIR)/swagger.yaml
	$(RM) -r pkg/client/api
	$(SWAGGER) generate client -f $^ -c pkg/client/api -m pkg/client/api/models --skip-validation
endif

.PHONY: depend-client
depend-client: pkg/client/api/nebula_client.go

.PHONY: depend
depend: depend-client

.PHONY: build
build: depend $(BIN_DIR)
	$(GO) build $(GOFLAGS) -o $(BIN_DIR)/$(CLI_DIST_NAME) ./cmd/nebula

.PHONY: test
test: depend
	$(GO) test $(GOFLAGS) ./...

.PHONY: dist
dist: $(CLI_DIST_SHA256)

.PHONY: clean
clean:
	$(RM) -r $(DEPEND_DIR)/
	$(RM) -r $(ARTIFACTS_DIR)/
	$(RM) -r $(BIN_DIR)/

.PHONY: $(CLI_DIST_BINS)
$(CLI_DIST_BINS): GOFLAGS += -a
$(CLI_DIST_BINS): GOOS = $(word 1,$(subst -, ,$*))
$(CLI_DIST_BINS): GOARCH = $(subst $(CLI_EXT_$(GOOS)),,$(word 2,$(subst -, ,$*)))
$(CLI_DIST_BINS): LDFLAGS += -extldflags "-static"
$(CLI_DIST_BINS): $(CLI_DIST_PREFIX)%: depend $(ARTIFACTS_DIR)
	env CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) \
		$(GO) build $(GOFLAGS) -o $@ -ldflags '$(LDFLAGS)' ./cmd/nebula

$(ARTIFACTS_DIR)/%.sha256: $(ARTIFACTS_DIR)/%
	cd $(dir $^) && $(SHA256SUM) $(notdir $^) >$(notdir $@)
