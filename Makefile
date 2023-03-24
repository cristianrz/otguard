all: bin bin/keytoaccess-server bin/keytoaccess-manager bin/keytoaccess-mkkey bin/keytoaccess-purgerules
	
bin:
	mkdir -p bin/

bin/keytoaccess-server:
	go build -o bin/keytoaccess-server keytoaccess-server/main.go

bin/keytoaccess-manager:
	go build -o bin/keytoaccess-manager keytoaccess-manager/main.go

bin/keytoaccess-mkkey:
	go build -o bin/keytoaccess-mkkey keytoaccess-mkkey/main.go

bin/keytoaccess-purgerules:
	cp scripts/keytoaccess-purgerules bin/

clean:
	rm -rf ./bin
