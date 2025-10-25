# JobDoe

JobDoe is a sandbox Job/Task runner using Docker and K8s APIs.

# Development

 * Simply install dependencies with `go get`.
 * Build with `go build`
 * Rename/Copy `env.example`, `regauth.json.example` to remove `.example` suffix. Then update the values to suit your environment.
 * After updating the `swaggo` comments, run `swag init -pd -q` to generate swagger documentation files.
