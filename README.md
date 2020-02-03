# reto - The release tool

## Installing

Install by running `go get -u github.com/ExploratoryEngineering/reto/cmd/reto` or
download the binaries from the release page.

## About

This is a tool to manage releases and change logs for releases. Strictly
speaking is't just a zip file generator with versioning.

Some of the the goals for the tool are:

* Work within the source tree without imposing any weird build steps.
* Be non-intrusive.
* Work with (almost) any source code.
* Automate release notes/change logs for the different version.
* Make releases that are easily mapped to source code revisions.
* Semantic versioning. It's not necessarily something that end users care
  about (they are mostly focused on features and if it works or not)
* Support software on multiple platforms and architectures.
* Handle sets of binaries and resources (such as documentation, upgrade scripts
  and so on)
* Possible to automate on build servers and in makefiles.

We're assuming the following:

* Git is the SCM of choice.
* Semantic versioning (or a scheme with major-minor-patch version numbering) is
  used.
* There's a changelog accompanying each release. The complete change log is
  included with each release, the newest is kept at the top.
* You can set the version numbers during the build.
* When you do a release everything is commited to Git.
* There's no roll back semantics. If something goes wrong when you generate a
  release you can use Git to check out the previous version.
  The contents of the `release` directory are kept under source control and
  in sync with the rest of the source. If something breaks (or a release goes
  terribly pear-shaped) you can roll back the changes via Git.
* Changelogs are written in markdown formats.

These are the things we *don't* care about:

* Nothing is built by the tool except the zip files.
* The source code. Or if there *is* any source code.
* The build process. You build the sofware the way *you* want to.
* The tool is a *tool* not a life style choice.
* How your commits look like. You can do one giant commit for each release or
  you can to thousands.
* What the change log for each version says. The tool does not write a change
  log for you. You *can* automate this by looking at the commits but that's not
  somthing the tool cares about. As long as the change log is free of "TODO"
  (it's a basic safeguard)

## What it does

The tool keeps track of the current version. The current version is always
incremented, either by patch number, minor or major version. You *can* edit the
version file manually to step backwards.

The changelog in `release/changelog.md` documents the new features, bugs and
other information you want to include with the current release.

When a a release is built a new directory is created (`release/n.n.n` where
`n.n.n` is the current version) and the current changelog (in
`release/changelog.md`) is copied to the release directory. The changelog from
*all* releases are merged into one big changelog.md and the built binaries and
associated files are copied into zip archives under `release/archives/n.n.n`, one
archive for each target set up in the configuration file. If the configuration
includes additionla files these are added to the archive as well.

The SHA256 hash for each file included in the archives are computed and written
to the file `release/archives/n.n.n/sha256sum_[name]_n.n.n.txt`

Once the release is completed the template changelog is copied from
`release/templates/changelog_template.md` to `release/changelog.md`, ready for
the next release.

## Configuration file

The configuration file is a JSON file in `release/config.json`. It looks like
the following (the example is for the tool itself):

```json
{
    "sourceRoot": ".",
    "name": "reto",
    "committerName": "Reto release tool",
    "committerEmail": "reto@exploratory.engineering",
    "targets": [
        "darwin-amd64", "linux-amd64", "windows-amd64", "linux-arm"
    ],
    "templates": [
        { "name": "changelog.md", "action": "concatenate" },
        { "name": "verify.txt",   "action": "include" },
        { "name": "releases.txt", "action": "concatenate" }
    ]
    "files": [
        { "id": "tool", "name": "bin/reto.darwin-amd64","target": "darwin-amd64" },
        { "id": "tool", "name": "bin/reto.linux-amd64", "target": "linux-amd64" },
        { "id": "tool", "name": "bin/reto.linux-arm", "target": "linux-arm" },
        { "id": "tool", "name": "bin/reto.windows-amd64", "target": "windows-amd64" },
        { "id": "readme", "name": "README.md", "target": "-" }
    ]
}
```

* `name` is the name of the product.
* `committerName` and `committerEmail` is the name and email used when tagging
  git releases.
* `sourceRoot` points to the source code root. Usually you'd want the release
  files to reside next to the tool itself but if you want to put the release
  files elsewhere you can set the path to point to the root of the source code.
* `targets` is a list of target strings that you plan to release. This can be a
  single string (if you have a single target) or multiple strings depending on
  what (and how) you build your released binaries. The strings can be anything
  as long as they are unique.
* `templates` are templates for files included in the release archives. They can
  either be included as is (with `action` set to "include") or they can be
  aggregated logs (like a changelog) that are aggregates of this and all previous
  releases.
* The `files` section contains a list of files. Each file has an `id` property
  that groups the differnt files into one. If a file is common for all targets
  use the target `'-'`, typically text files, PDFs or images. If you list three
  different targets in the `targets` property you must have the same three targets
  for each `id`. The `name` field is the file name of the file to include.

## Fancy naming of commit hashes

Commit hashes gets user-friendly names. The first three bytes are translated
into an adjective-name string that identifies the commit. This is enough to
get relatively unique version names for different releases.

Names can get quite interesting and someone might be offended if they put their
mind to it but it's a nice option to have if you have frequent releases and
use the semantic version for internal features as well as external features.

(Your end users might not care one iota about the new backwards incompatible
flag `--upsert-before-commit-with-logstash-logging' but you still have a version
that the end users can relate to)

Commit names are `withering-gustie` `well-spoken-amin` `timid-chesley`,
`effervescent-laurie` and `appropriate-yulissa`.

## Setting versions in Go builds

You can use the `ldflags` parameter to set variables in Go packages. If you look
at the tool itself the file `pkg/version/version.go` contains the variables
`Number`, `Name`, `CommitHash` and `BuildTime` that has suitable defaults.

If you build the executable with the `ldflags` parameter set like this:

```shell
go build -ldflags "-X github.com/ExploratoryEngineering/reto/pkg/version.Number=1.0.0 -X github.com/ExploratoryEngineering/reto/pkg/version.Name=psychological-karson -X github.com/ExploratoryEngineering/reto/pkg/version.CommitHash=e02bc6cd8ab3edd3d6e8874d1e97c08ef6c5db49 -X github.com/ExploratoryEngineering/reto/pkg/version.BuildTime=2020-02-01T22:35"
```

The variables will be set to `1.0.0`, `psychological-karson`,
`e02bc6cd8ab3edd3d6e8874d1e97c08ef6c5db49` and `2020-02-01T22:35` respectively.

The Makefile in the project contains a sample configuration where the value is
retrieved from the tool itself. For obvious reasons you should built a non-versioned
binary before you can set the version in the build.

## C/C++

For gcc you can use the `-D` flag to set defines at compile time. Use the tool
to extract the version info. The Makefile in the project builds a Go executable
and a `gcc` build should be quite similar.
