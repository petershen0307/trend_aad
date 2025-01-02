.PHONY: clean test build_apps build@$(1)

# static link net package: -tags netgo
# https://go.dev/doc/go1.2
GO_BUILD_FLAGS ?=

ADC_GITHUB_PERSONAL_TOKEN ?=

APPS := $(wildcard ./cmd/*/)
OUTDIR := out
VERSION ?= v1.0.0
GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)

APP_NAME ?= trend_aad
PACK_APP_NAME ?= $(APP_NAME).tar.gz

clean:
	rm -rf ./$(OUTDIR)/*

test:
	go test -v ./...

build_apps: clean $(foreach app,$(APPS),build@$(app))

define BUILD_APPS
build@$(1):
	mkdir -p $(OUTDIR)/$(shell basename $(1))
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build $(GO_BUILD_FLAGS) -o $(OUTDIR)/$(shell basename $(1)) $(1)
endef
$(foreach app,$(APPS),$(eval $(call BUILD_APPS,$(app))))

install_chromium_dependencies:
	sudo apt install -y libnss3 libgbm-dev libasound2

pack_app:
	echo "$(PACK_APP_NAME)" > ./$(OUTDIR)/$(APP_NAME)/VERSION
	tar -zcv -f ./$(OUTDIR)/$(PACK_APP_NAME) -C ./$(OUTDIR)/$(APP_NAME) .