all:
	cd cmd/releasetool && go build -o ../../bin/releasetool


# This sets the build parameters in the binary at build time
version := $(shell bin/releasetool version)
commit := $(shell bin/releasetool hash)
name := $(shell bin/releasetool hashname)
build_time := $(shell date +"%Y-%m-%dT%H:%M")

ldflags := -X github.com/ExploratoryEngineering/releasetool/pkg/version.Number=$(version) \
	-X github.com/ExploratoryEngineering/releasetool/pkg/version.Name=$(name) \
	-X github.com/ExploratoryEngineering/releasetool/pkg/version.CommitHash=$(commit) \
	-X github.com/ExploratoryEngineering/releasetool/pkg/version.BuildTime=$(build_time)

local-rel:
	cd cmd/releasetool && go build -ldflags "$(ldflags)" -o ../../bin/releasetool

builds: local-rel
	cd cmd/releasetool && GOOS=darwin GOARCH=amd64 go build -ldflags "$(ldflags)" -o ../../bin/releasetool.darwin-amd64	
	cd cmd/releasetool && GOOS=linux GOARCH=amd64 go build -ldflags "$(ldflags)" -o ../../bin/releasetool.linux-amd64
	cd cmd/releasetool && GOOS=windows GOARCH=amd64 go build -ldflags "$(ldflags)" -o ../../bin/releasetool.linux-amd64
	cd cmd/releasetool && GOOS=linux GOARCH=arm go build -ldflags "$(ldflags)" -o ../../bin/releasetool.linux-arm