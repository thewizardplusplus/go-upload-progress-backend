# Change Log

## [v1.0.0](https://github.com/thewizardplusplus/go-upload-progress-backend/tree/v1.0.0) (2022-11-07)

The major version. Implement the service API.

- Serve files:
  - Serve static files
  - Serve uploaded files
- Service API:
  - Load a file list:
    - information about each file:
      - name
      - size (in bytes)
      - modification time
  - Upload a file:
    - add a random suffix to a duplicated file name
  - Delete a file
  - Delete all files
- Add the [Postman](https://www.postman.com/) documentation for the service API
