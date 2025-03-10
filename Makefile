BINARY=ai-knowledge

default:
	@echo 'Usage of make: [ build | linux | windows | run | clean ]'

build: 
	go build -buildvcs=false -o ./bin/${BINARY} ./

linux: 
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./bin/${BINARY} ./

windows: 
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o ./bin/${BINARY}.exe ./

run: build
	cd bin && ./${BINARY}

clean: 
	cd bin && rm -f ./${BINARY}*

.PHONY: default build linux run docker docker_push clean