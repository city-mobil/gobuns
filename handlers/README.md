# Handlers

Пакет для работы с HTTP-хендлерами на стороне сервера.

## AccessLogHandler

Handler для логирования входящих HTTP запросов и ответов.

### Пример использования

```go
// endpoint is used to init NR transaction,
// will be deprecated in future.
endpoint := "/api/v1/rider/auth"
logger := zlog.New(os.Stdout) // use your instance of zlog
accessLogger := handlers.NewAccessLogger(endpoint, WithLogger(logger))
// if you use hlog from zlog, 
// you might extract logger from the request context.
// accessLogger := handlers.NewAccessLogger(endpoint, WithLoggerFromReq())

http.Handle(endpoint, http.TimeoutHandler(handlers.AccessLogHandler(
        accessLogger,
        myHandler,
    ), time.Second, time.Second,
))

http.Handle(endpoint, AccessLogHandleFunc(
    accessLogger,
    myHandler,
))
```

### Фильтрация логирования входящих запросов

По умолчанию AccessLogger логирует все входящие запросы, но его можно настроить на запись только нужных.

```go
// Создаем логер с фильтрацией всех запросов со статусом ответа 2xx.
accessLogger := handlers.NewAccessLogger(endpoint, WithFilter(LogExcept2xx), WithLoggerFromReq())

// Создаем логер с пользовательским фильтром.
func myFilter(int code, err error) bool {
    if code == 401 {
        return true
    }
   
    return err != nil
}
accessLogger := handlers.NewAccessLogger(endpoint, WithFilter(LogExcept2xx), WithLoggerFromReq())
```

## Обработка ошибок

### RequestError

Является индикатором ошибки запроса. В случае, если произошла какая-то логическая ошибка(например, не найден заказ в
базе), необходимо возвращать `RequestError`.

```go
func HandleRequest(ctx handlers.Context) (interface{}, error) {
    if ctx.HTTPRequest().Method != http.MethodGet {
        return nil, &RequestError{
            Status: 404,
            Message: "invalid",
            Code: "invalid",
        }
    }
}
```
