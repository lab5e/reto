# Release tool

This is a tool to package and manage both releases and versions.

## Random ramblings

Basic assumptions:

* Binaries are located in bin/
* Release notes and documentation are placed in release/
* Individual releases are placed in release/<version>
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

Archive <name>-<version>-<platform>.zip
SHA256 checksum: <name>-<version>-<platform>.sha256.txt

bin/ - binaries
doc/ - documentation
<version>-changelog.md - release note/aggregated change log
<version>-sha256.txt - checksums for binaries

Change logs - internal and external. Internal change log documents
management plus external visible features, external change log
documents externally visible changes.

Each version is

<version>: <code name> (date)

Internal

External

