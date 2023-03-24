DESTDIR ?= /
PREFIX ?= /usr/local/otguard
BUILD_DIR = build
BIN_DIR = $(BUILD_DIR)$(PREFIX)/bin
SHARE_DIR = $(BUILD_DIR)$(PREFIX)/share/otguard
ETC_DIR = $(BUILD_DIR)/etc/otguard
CRON_DIR = $(BUILD_DIR)/etc/cron.d
SYSTEMD_DIR = $(BUILD_DIR)/usr/local/lib/systemd/system
VERSION := $(shell date +%y%m%d)
GO_FLAGS = --ldflags '-extldflags "-static"'


all: build $(BIN_DIR)/otguard-web $(BIN_DIR)/otguardd $(BIN_DIR)/otguard-mksecret $(BIN_DIR)/otguard-purgerules $(SHARE_DIR)/login.html $(SYSTEMD_DIR)/otguardd.service $(CRON_DIR)/otguard-cron

build:
	@mkdir -p $(BIN_DIR) $(SHARE_DIR) $(ETC_DIR) $(SYSTEMD_DIR) $(CRON_DIR)

$(BIN_DIR)/otguard-web: otguard-web/main.go
	go build $(GO_FLAGS) -o $@ $<

$(BIN_DIR)/otguardd: otguardd/main.go
	go build $(GO_FLAGS) -o $@ $<

$(BIN_DIR)/otguard-mksecret: scripts/otguard-mksecret
	cp $< $@

$(BIN_DIR)/otguard-purgerules: scripts/otguard-purgerules
	cp $< $@

$(SHARE_DIR)/login.html: login.html
	cp $< $@

$(SYSTEMD_DIR)/otguardd.service: otguardd.service
	cp $< $@

$(CRON_DIR)/otguard-cron: otguard-cron
	cp $< $@

clean:
	rm -rf $(BUILD_DIR)

dist:
	tar -czvf otguard-$(VERSION).tgz build/ README.md LICENSE.md install scripts/

.PHONY: all build clean
