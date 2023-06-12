# go-upload-progress-backend

[![GoDoc](https://godoc.org/github.com/thewizardplusplus/go-upload-progress-backend?status.svg)](https://godoc.org/github.com/thewizardplusplus/go-upload-progress-backend)
[![Go Report Card](https://goreportcard.com/badge/github.com/thewizardplusplus/go-upload-progress-backend)](https://goreportcard.com/report/github.com/thewizardplusplus/go-upload-progress-backend)

Back-end of the service that implements a simple file manager with a display of file upload progress.

The main challenge of the project was to write it without using third-party libraries.

## Features

- RESTful API:
  - models:
    - file info:
      - fields:
        - name;
        - size (in bytes);
        - modification time (in the [RFC 3339](https://www.rfc-editor.org/rfc/rfc3339.html) format with sub-second precision);
      - operations:
        - getting info about all uploaded files:
          - sort the results by modification time in descending order;
        - uploading a file:
          - generate unique filenames;
        - deleting a file by a filename;
        - deleting all uploaded files;
  - representing:
    - in a JSON:
      - payloads of responses;
    - as a plain text:
      - errors;
- generating unique filenames:
  - add a random suffix to a duplicated file name:
    - use a sequence of random bytes as the suffix;
    - format the suffix in base 16, lower-case, two characters per byte;
  - restrict a number of tries to generate a unique filename;
- server:
  - additional routing:
    - serving static files;
    - serving uploaded files;
  - storing settings in environment variables;
  - supporting graceful shutdown;
  - logging:
    - logging requests;
    - logging errors;
  - panics:
    - recovering on panics;
    - logging of panics;
- distributing:
  - [Docker](https://www.docker.com/) image;
  - [Docker Compose](https://docs.docker.com/compose/) configuration;
- utilities:
  - utility for generating a dummy file of a specified size, filled with random bytes.

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
- `STATIC_FILE_DIR` &mdash; path to static files (default: `./static`);
- `UPLOADED_FILE_DIR` &mdash; path to uploaded files (default: `./files`);
- `MAXIMUM_NUMBER_OF_TRIES` &mdash; maximum number of tries to generate a unique filename (default: `10000`);
- `RANDOM_SUFFIX_BYTE_COUNT` &mdash; byte count in a random suffix of duplicate filenames (default: `4`).

## API Description

API description:

- in the [Swagger](http://swagger.io/) format: [docs/swagger.yaml](docs/swagger.yaml);
- in the format of a [Postman](https://www.postman.com/) collection: [docs/postman_collection.json](docs/postman_collection.json).

## Utilities

- [generate-dummy](tools/generate-dummy) &mdash; the utility for generating a dummy file of a specified size, filled with random bytes

## License

The MIT License (MIT)

Copyright &copy; 2022-2023 thewizardplusplus
