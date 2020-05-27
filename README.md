Introduction
============

Pondi is a helper tool to build release assets (binaries, source code archive),
create a release on GitHub and upload the files. It is specifically tailored to
the programming language Go and the style and requirements tools around the
[restic](https://restic.net) backup program have.

Goals (what pondi should do once it's finished):
 * Generate a changelog file using [`calens`](https://github.com/restic/calens)
 * Update the version in the source code as well as the `VERSION` file (if it exists)
 * Create a new release tag (if it does not exist yet)
 * Build binaries for all architectures and OS, in a reproducible way
 * Create `tar.gz` files for each architecture/OS
 * Build a source archive
 * Create a `SHA256SUMS` of all created files and and sign it
 * Generate release notes using `calens`
 * Create a release on GitHub with the generated release notes and upload the assets

The program should be idempotent, so if you re-run pondi it should continue where it left of instead of starting from the beginning again.

Process
=======

Pondi helps automating the following steps:

 1. Tag a new release, operate on the checkout
    1. Check that the target tag (`v1.2.3` for version `1.2.3`) does not exist yet.
    1. Run hooks (`go mod download`, `go generate`), abort if any command fails.
    1. Check that the source has not been modified (so all changes were already committed), abort if otherwise.
    1. Run checks on the source code.
 2. Build binaries and the source code archive in a reproducible way, operate on the tagged release
 3. Publish a new release on GitHub and upload all assets
