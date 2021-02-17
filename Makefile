all:
	cd cmd/reto && go build -o ../../bin/reto


# This sets the build parameters in the binary at build time
version := $(shell bin/reto version)
commit := $(shell bin/reto hash)
name := $(shell bin/reto hashname)
build_time := $(shell date +"%Y-%m-%dT%H:%M")

ldflags := -X github.com/lab5e/reto/pkg/version.Number=$(version) \
	-X github.com/lab5e/reto/pkg/version.Name=$(name) \
	-X github.com/lab5e/reto/pkg/version.CommitHash=$(commit) \
	-X github.com/lab5e/reto/pkg/version.BuildTime=$(build_time)

local-rel:
	cd cmd/reto && go build -ldflags "$(ldflags)" -o ../../bin/reto

builds: local-rel
	cd cmd/reto && GOOS=darwin GOARCH=amd64 go build -ldflags "$(ldflags)" -o ../../bin/reto.darwin-amd64
	cd cmd/reto && GOOS=linux GOARCH=amd64 go build -ldflags "$(ldflags)" -o ../../bin/reto.linux-amd64
	cd cmd/reto && GOOS=windows GOARCH=amd64 go build -ldflags "$(ldflags)" -o ../../bin/reto.windows-amd64
	cd cmd/reto && GOOS=linux GOARCH=arm go build -ldflags "$(ldflags)" -o ../../bin/reto.linux-arm
