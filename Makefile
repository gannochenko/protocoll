run:
	@go run cmd/main.go generate --folder ./sample-data --name "Sample import" 1> ./coll.json

build:
	@./script/build.sh

release:
	@git tag v$(version)
	@git push origin v$(version)
