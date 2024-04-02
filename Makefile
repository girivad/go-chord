all: protos chord-docker
protos:
	../protoc-25.3-win64/bin/protoc.exe ./Proto/overlay.proto --go_out=. --go-grpc_out=. --go-grpc_opt=paths=source_relative --go_opt=paths=source_relative
chord:
	go build
chord-docker:
	docker build -t go-chord-node .