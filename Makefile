MOD_NAME:=airer_driver

CMD_INSTALL:=/opt/piremd

CMD_BIN_DIR:=bin

CMD_PACKAGES:=$(shell go list ./cmd/...)

GO_FILES:=$(shell find . -type f -name '*.go' -print)

CMD_BIN:=$(CMD_PACKAGES:$(MOD_NAME)/cmd/%=$(CMD_BIN_DIR)/%)

BUILD_OPT:=-ldflags="-s -w" -trimpath

.PHONY: clean
clean:
	rm $(CMD_BIN_DIR)/**

.PHONY: build
build: $(CMD_BIN)

$(CMD_BIN): $(GO_FILES)
	$(BUILD_ENV) go build $(BUILD_OPT) -o $(CMD_BIN_DIR) $(@:$(CMD_BIN_DIR)/%=$(MOD_NAME)/cmd/%)


.PHONY: install
install:
	cp $(CMD_BIN).so $(CMD_INSTALL)

.PHONY: update
update: install

.PHONY: remove
remove:
	rm $(CMD_BIN:$(CMD_BIN_DIR)/%=$(CMD_INSTALL)/%).so

.PHONY: purge
purge: remove
