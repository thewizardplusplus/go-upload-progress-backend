# go-upload-progress-backend

[![GoDoc](https://godoc.org/github.com/thewizardplusplus/go-upload-progress-backend?status.svg)](https://godoc.org/github.com/thewizardplusplus/go-upload-progress-backend)
[![Go Report Card](https://goreportcard.com/badge/github.com/thewizardplusplus/go-upload-progress-backend)](https://goreportcard.com/report/github.com/thewizardplusplus/go-upload-progress-backend)

Back-end of the service that implements a simple file manager with a display of file upload progress.

## Installation

```
$ go install github.com/thewizardplusplus/go-upload-progress-backend@latest
```

## Usage

```
$ go-upload-progress-backend
```

Environment variables:

- `SERVER_ADDRESS` &mdash; server URI (default: `:8080`);
- `PUBLIC_FILE_DIR` &mdash; path to the project's front-end (default: `./public`);
- `UPLOADED_FILE_DIR` &mdash; path to uploaded files (default: `./files`).

## API Description

API description:

- in the [Swagger](http://swagger.io/) format: [docs/swagger.yaml](docs/swagger.yaml);
- in the format of a [Postman](https://www.postman.com/) collection: [docs/postman_collection.json](docs/postman_collection.json).

## License

The MIT License (MIT)

Copyright &copy; 2022 thewizardplusplus
