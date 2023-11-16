# Worker Engine
Worker engine to execute TSF tasks in a sandbox.

# Development

 * Simply install dependencies with `go get`.
 * Build with `go build`
 * Copy `deb/env.example`, `deb/regauth.json.example` without the `.example` suffix to the same directory the `workerengine` binary is. Then update the values to suit your environment.
 * After updating the `swaggo` comments, run `swag init -pd -q` to generate swagger documentation files.

# Deploy

 * You can package for APT package manager with `make build`. Cleanup package files with `make clean`.