# Change Log

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
