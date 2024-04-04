# Worker Engine
Worker engine to execute TSF tasks in a sandbox.

# Dependencies

The following system libraries are required for build:
 * `libdevmapper-dev`
 * `libbtrfs-dev`
 * `libgpgme-dev`

# Development

 * Simply install dependencies with `go get`.
 * Build with `go build`
 * Rename/Copy `env.example`, `regauth.json.example` to remove `.example` suffix. Then update the values to suit your environment.
 * After updating the `swaggo` comments, run `swag init -pd -q` to generate swagger documentation files.

# Build

 * You can build a Debian package (`apt`, `dpkg`) with `make build`.
 * Cleanup package files with `make clean`.