# Health

Пакет для работы с health check.

API соответствует [RFC](https://tools.ietf.org/id/draft-inadarei-api-health-check-01.html)

## Описание структур и интерфейсов

### Checkable interface

Интерфейс, который используется для HealthCheck.

#### Ping

Метод, который выполняет основной health check.

#### ComponentID

Метод, который возвращает ComponentID для необходимого health check.

#### ComponentType

Метод, который возвращает ComponentType для необходимого health check.

#### Name

Метод, который возвращает имя health check.

#### Пример использования

```go
ch := NewChecker(CheckerOptions{
    Version: "42", // может быть commithash.
    ReleaseID: "42", // может быть git tag.
    ServiceID: "42", // может быть hostname или ip-адресом.
})

ch.AddCallback("some_callback_name", CheckCallback(func() *CheckResult{
    return &CheckResult{} // какой-то ответ.
}))
```

## Описание методов

### NewHandler

Создаёт новый HTTP Handler для проверки healthcheck.

В случае, если проверка `Checker` вернула хотя бы один `Fail`, то статусом HTTP ответа будет `500 
Internal Server Error`.

В любом другом случае(`pass` или `warn`) возвращается `HTTP 200 OK`.

#### Пример использования

```go
ch := NewChecker(CheckerOptions{
    ReleaseID: "42",
    ServiceID: "42",
    Version: "42",
})
ch.AddCallback("callback_name", CheckCallback(func() *CheckResult{
    return &CheckResult{}
}))

http.Handle("/health", NewHandler(ch))
http.ListenAndServe(":4242", nil)
```
