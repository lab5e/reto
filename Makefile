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

rel:
	cd cmd/releasetool && go build -ldflags "$(ldflags)" -o ../../bin/releasetool