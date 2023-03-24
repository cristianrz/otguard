DESTDIR ?= /
PREFIX ?= /usr/local
BUILD_DIR = build
BIN_DIR = $(BUILD_DIR)$(PREFIX)/bin
SHARE_DIR = $(BUILD_DIR)$(PREFIX)/share/otguard
ETC_DIR = $(BUILD_DIR)/etc/otguard
SYSTEMD_DIR = $(ETC_DIR)/systemd/system

all: build $(BIN_DIR)/otguard-web $(BIN_DIR)/otguardd $(BIN_DIR)/otguard-mksecret $(BIN_DIR)/otguard-purgerules $(SHARE_DIR)/login.html $(ETC_DIR)/cert.pem $(ETC_DIR)/key.pem $(ETC_DIR)/secrets $(SYSTEMD_DIR)/otguardd.service

init:
	./scripts/otguard-mkcert
	./scripts/otguard-mksecret
	
build:
	@mkdir -p $(BIN_DIR) $(SHARE_DIR) $(ETC_DIR) $(SYSTEMD_DIR)

$(BIN_DIR)/otguard-web: otguard-web/main.go
	go build -o $@ $<

$(BIN_DIR)/otguardd: otguardd/main.go
	go build -o $@ $<

$(BIN_DIR)/otguard-mksecret: scripts/otguard-mksecret
	cp $< $@

$(BIN_DIR)/otguard-purgerules: scripts/otguard-purgerules
	cp $< $@

$(SHARE_DIR)/login.html: login.html
	cp $< $@

$(ETC_DIR)/cert.pem: cert.pem
	cp $< $@

$(ETC_DIR)/key.pem: key.pem
	cp $< $@

$(ETC_DIR)/secrets: secrets
	cp $< $@

$(SYSTEMD_DIR)/otguardd.service: otguardd.service
	cp $< $@

clean:
	rm -rf $(BUILD_DIR)

dist:
	cd build/ && tar -czvf ../otguard.tgz .

.PHONY: all build clean
