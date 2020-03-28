.PHONY: build

clean:
	@rm -rf ./build
	
build-mac: clean
	env GOOS=darwin GOARCH=amd64 go build -o ./build/devsup main.go
	
build-linux: clean
	env GOOS=linux GOARCH=amd64 go build -o ./build/devsup main.go

dep:
	dep ensure