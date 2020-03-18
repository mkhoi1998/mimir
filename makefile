.PHONY: build

clean:
	@rm -rf ./build
	
build: clean
	go build -o ./build/devsup main.go

dep:
	dep ensure