
## Changelog 1.3.0: infatuated-burk

### Features

Git integration. Check for untracked files before creating a release. All changes
must be committed before a release is done. Also use to go-git library for git
operations.

Tag the version when making a release. This ensures that the version name is
consistent. The tag is created before the changelog is moved so the current
version is in `release/VERSION` and the current changelog will be in
`release/changelog.md`.

Updated status command with -V flag to show verbose output.
