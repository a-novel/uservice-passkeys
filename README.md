# Service: Passkeys

![GitHub Actions Workflow Status](https://img.shields.io/github/actions/workflow/status/a-novel/uservice-passkeys/main.yaml)
[![codecov](https://codecov.io/gh/a-novel/uservice-passkeys/graph/badge.svg?token=Tyo7MYuQ75)](https://codecov.io/gh/a-novel/uservice-passkeys)

![GitHub repo file or directory count](https://img.shields.io/github/directory-file-count/a-novel/uservice-passkeys)
![GitHub code size in bytes](https://img.shields.io/github/languages/code-size/a-novel/uservice-passkeys)

![Coverage graph](https://codecov.io/gh/a-novel/uservice-passkeys/graphs/sunburst.svg?token=Tyo7MYuQ75)

Passwords and secret keys manager.

### Prerequisites

- [Go](https://go.dev/doc/install)
- Make
    - macOS:
      ```bash
      brew install make
      ```
    - Ubuntu:
      ```bash
      sudo apt-get install make
      ```
    - Windows: Install [chocolatey](https://chocolatey.org/install) (from a PowerShell with admin privileges), then run:
      ```bash
      choco install make
      ```

Install the project dependencies.

```bash
go get ./... && go mod tidy
```

## Run the project locally

### From command line

```bash
make run
```

### From GitHub packages

You can get a working version of the service from the GitHub packages, using this image:

```
ghcr.io/a-novel/uservice-passkeys/master:latest
```

> You can replace the `master` part with the name of any branch, to retrieve the image built from that branch. Or
> replace `latest` with the sha of a commit to get the image built from that commit.

The image needs 2 environment variables to work:

- `PORT`: The port the service will listen to.
- `DSN`: The connection string to a postgres database.

### Make test queries

You can run queries on the go from a terminal using [grpcurl](https://github.com/fullstorydev/grpcurl). Below is an
example for the global health check (available on all services).

```bash
grpcurl -plaintext -d '{"service": ""}' localhost:4003 grpc.health.v1.Health/Check
```

## Work on the project

Make sure the project files are properly formatted.

```bash
make format
```

Run tests.

```bash
make test
```

Make sure your code is compliant with the linter.

```bash
make lint
```

If you create / update interfaces signatures, make sure to update the mocks.

```bash
make mocks
```
