# Change Log

## [v1.5.0](https://github.com/thewizardplusplus/go-upload-progress-backend/tree/v1.5.0) (2023-06-12)

Use the [github.com/thewizardplusplus/go-writable-fs](https://github.com/thewizardplusplus/go-writable-fs) package and add the utility for generating a dummy file of a specified size, filled with random bytes.

- Service API:
  - Upload a file:
    - Restrict a number of tries to generate a unique filename
- Utilities:
  - Add the utility for generating a dummy file of a specified size, filled with random bytes:
    - Parse the options
    - Parse a file size
    - Calculate a file size in bytes
    - Generate a filename
    - Generate a file
    - Display progress
- Perform refactoring:
  - Use the [github.com/thewizardplusplus/go-writable-fs](https://github.com/thewizardplusplus/go-writable-fs) package:
    - Use the `fsutils.ReadDirEntriesByKind()` function
    - Use the `writablefs.WritableFS` interface
    - Use the `writablefs.DirFS` structure
  - Add the `main.getIntEnv()` function

## [v1.4.0](https://github.com/thewizardplusplus/go-upload-progress-backend/tree/v1.4.0) (2023-05-12)

Improve uploading files, unique filename generating, and logging, implement graceful server shutdown and perform refactoring.

- Service API:
  - Load a file list:
    - Improve field naming in the `FileInfo` structure
  - Upload a file:
    - Fix the bug with uploading large files
    - Make sure that the file will not be overwritten on uploading
    - Return saved file information on file uploading
    - generating unique filenames:
      - Read random suffix byte count from the environment variable
      - Use the `crypto/rand` package for unique filename generating
      - Improve the format of a random suffix
- Implement logging:
  - Add the `middlewares.responseWriterWrapper` structure
- Implement graceful server shutdown
- Rename the `public` directory to `static`
- Perform refactoring:
  - Add the `entities.FileInfoGroup` type:
    - Add the `entities.NewFileInfoFromDirEntry()` function
    - Add the `entities.ByModificationTime` type
  - Extract the `fsutils.ReadDirFiles()` function
  - Extract the `usecases/generators.FilenameGenerator` structure
  - Add the `handlers.FileHandler.writeAsJSON()` method
  - of the `gateways/writablefs` package:
    - Rename the `writablefs.WritableFS` structure to `DirFS`
    - Improve the implementation of the `fs.FS` interface with the `writablefs.DirFS` structure

## [v1.3.0](https://github.com/thewizardplusplus/go-upload-progress-backend/tree/v1.3.0) (2022-11-12)

Add the [Docker](https://www.docker.com/) image.

- Read the service parameters from environment variables
- Add the [Docker](https://www.docker.com/) image:
  - Add the [Docker Compose](https://docs.docker.com/compose/) configuration
- Add the [Swagger](http://swagger.io/) documentation for the service API
- Perform refactoring:
  - Move the main code to a separate directory

## [v1.2.0](https://github.com/thewizardplusplus/go-upload-progress-backend/tree/v1.2.0) (2022-11-10)

Implement logging.

- Implement logging:
  - Implement error logging in the `gateways/handlers.FileHandler` structure
  - Add the `gateways/handlers/middlewares` package:
    - Add the `middlewares.ApplyMiddlewares()` function
    - Add the `middlewares.LoggingMiddleware()` function
    - Add the `middlewares.RecoveringMiddleware()` function
    - Use the middlewares in the main code
- Perform refactoring:
  - Add the `main.makeLogger()` function
  - Add the `main.makeFileServer()` function

## [v1.1.0](https://github.com/thewizardplusplus/go-upload-progress-backend/tree/v1.1.0) (2022-11-08)

Perform refactoring and improve unique filename generating.

- Service API:
  - Upload a file:
    - Improve collision processing on unique filename generating
    - Set up the random number generator seed
- Perform refactoring:
  - Extract the `entities.FileInfo` structure
  - Extract the `usecases` package:
    - Extract the `usecases.FileUsecase.GetFiles()` method
    - Extract the `usecases.FileUsecase.SaveFile()` method
    - Extract the `usecases.FileUsecase.DeleteFile()` method
    - Extract the `usecases.FileUsecase.DeleteFiles()` method
    - Add the `usecases.FileUsecase.readDirFiles()` method
    - Add the `usecases.makeRandomFilename()` function
  - Extract the `gateways/handlers` package:
    - Extract the `handlers.FileHandler.GetFiles()` method
    - Extract the `handlers.FileHandler.SaveFile()` method
    - Extract the `handlers.FileHandler.DeleteFile()` method
    - Extract the `handlers.FileHandler.DeleteFiles()` method
    - Add the `handlers.FileHandler.ServeHTTP()` method
  - Add the `gateways/writablefs.WritableFS` structure

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
