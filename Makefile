all: build build/bin/otguard-server build/bin/otguard-manager build/bin/otguard-mkkey build/bin/otguard-purgerules build/usr/share/otguard/login.html build/etc/otguard/cert.pem build/etc/otguard/key.pem build/etc/otguard/secrets
	
build:
	mkdir -p build/bin/
	mkdir -p build/usr/share/otguard/
	mkdir -p build/etc/otguard/

build/bin/otguard-server: otguard-server/main.go
	go build -o build/bin/otguard-server otguard-server/main.go

build/bin/otguard-manager: otguard-manager/main.go
	go build -o build/bin/otguard-manager otguard-manager/main.go

build/bin/otguard-mkkey: otguard-mkkey/main.go
	go build -o build/bin/otguard-mkkey otguard-mkkey/main.go

build/bin/otguard-purgerules: scripts/otguard-purgerules
	cp scripts/otguard-purgerules build/bin/

build/usr/share/otguard/login.html: login.html
	cp login.html build/usr/share/otguard/login.html

build/etc/otguard/cert.pem:
	cp cert.pem build/etc/otguard/

build/etc/otguard/key.pem:
	cp key.pem build/etc/otguard/

build/etc/otguard/secrets:
	cp secrets build/etc/otguard/

clean:
	rm -rf ./build
