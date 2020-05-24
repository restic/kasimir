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
 * Build a source archive
 * Create a `SHA256SUMS` and sign it
 * Create a release on GitHub and upload the assets
 * It should be idempotent, if you re-run pondi it should continue where it left of

Process
=======

Pondi helps automating the following steps:

 1. Tag a new release
 2. Build binaries and the source code release in a reproducible way
 3. Publish a new release on GitHub and upload all assets
