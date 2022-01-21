
BINARY=syncer
GOBUILD=go build -ldflags "-s -w" -o ${BINARY}
GOCLEAN=go clean

# Build
build:
	$(GOCLEAN)
	$(GOBUILD)


clean:
	$(GOCLEAN)
	$(RMTARGZ)

inside=root@192.168.10.3
dev:
	go fmt
	$(GOCLEAN)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD)
	ssh -p22 ${inside} "mkdir -p /data/git-syncer/repo"
	scp ./${BINARY} ${inside}:/data/git-syncer
	$(GOCLEAN)