# Sheep Farming Backend

This repository contains a Go backend for managing sheep farming operations. The project includes a small JavaScript front-end under the `front/` directory.

## Docker

A `Dockerfile` is provided for building a production container. Build and run the container with:

```bash
docker build -t sheep-farm .
docker run -p 8080:8080 sheep-farm
```

Environment variables such as `FIREBASE_CREDENTIALS_PATH` should be provided when running the container.

## CI/CD

Continuous integration is configured via GitHub Actions in `.github/workflows/ci.yml`. The workflow performs the following on each push or pull request:

- Checks out the code.
- Sets up Go 1.23.
- Builds the application and runs `go vet`.
- Builds the Docker image to ensure the container configuration is valid.

