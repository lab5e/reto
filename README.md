# Release tool

This is a tool to package and manage both releases and versions. Both a version
number and a git hash is kept to keep track of the build. The git hash is encoded
to a human readable short hash.

## Goals

The goal of this tool is to aid in binary releases with archives. Releases can
be managed in the source tree itself or outside the source tree. The release
files should be under SCM (but the generated archives shouldn't)

The release tool can either be used as part of a manual release or as part
of a scripted release.

### What the releasetool does

* Keeps track of version numbers
* Create archives with releases
* Manages the changelog for each release

### What the releasetool doesn't do

* Upload releases
* Build release
* Writes a changelog for you
* Talks to external systems

This is something you want to do in your own scripts.

## Commands

init - set up the release tool and the release directory
version - dump the current version to stdout. Useful for build scripts
bump - increase major, minor or patch version
hashname - dump the name of the current git hash
namehash - convert name into git hash

## On hash naming

It works like this:

```shell
$ git rev-parse --short=6 HEAD
2b85aa
$ bin/releasetool hashname
frantic-bennie
$ bin/releasetool namehash frantic-bennie
2b85aa
```

To find the matching commit to a name use `git log -n 1 $(bin/releasetool namehash calculating-aldona)`

## Configuration file

```json
{
    "sourceRoot": ".",
    "architectures": [
        "amd64"
    ],
    "oses": [
        "darwin"
    ],
    "files": [
        {
            "id": "tool",
            "name": "bin/release",
            "os": "darwin",
            "arch": "amd64"
        }
    ]
}
```

## Release information in builds

### Go

```make
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
```

### C/C++

Use the -D[define] to set the configuration option at build time

### West/Zephyr

## Random ramblings

Basic assumptions:

* Binaries are located in bin/
* Release notes and documentation are placed in release/
* Individual releases are placed in release/[version]
* Git is used for SCM
* The current working version is in release/VERSION
* Binaries are signed
* The list of binaries are configurable
* Commit hash + version
* we can use several sets of change logs and release notes

Commands

* init - set up release tool
* version - show current working version

* bump - bump versions
* release - package a release. Move files into release structure, aggregate change log
* upgrade - do a database upgrade
* versions - list versions
* rebuild-changes - rebuild change log

## Templates

Release templates are in release/template

* note.md

## Release package structure

There's a release for each operating system
Each release contains all binaries

Archive [name]-[version]-[platform].zip
SHA256 checksum: [name]-[version]-[platform].sha256.txt

bin/ - binaries
doc/ - documentation
[version]-changelog.md - release note/aggregated change log
[version]-sha256.txt - checksums for binaries

Change logs - internal and external. Internal change log documents
management plus external visible features, external change log
documents externally visible changes.

Each version is

[version]: [code name] (date)

Internal

External

## Release process

1. Make sure it isn't already released (ie version isn't in release list)
1. Make release directory
1. Construct change log by appending current release note
1. Checksum binaries, write checksum file
1. Make archive with binaries, checksum file and change log
