# IPLookup

Позволяет извлечь IP адрес клиента из HTTP запроса согласно заданным правилам.

По умолчанию IP адрес извлекается в следующем порядке пока не будет найден:

* `RemoteAddr` - смотрите [документацию](https://golang.org/pkg/net/http/#Request),
* HTTP заголовок `X-Real-IP`,
* HTTP заголовок `X-Forwarded-For`.

В заголовке `X-Forwarded-For` могут быть перечислены несколько IP адресов, поэтому индекс нужного адреса можно настроить
через параметр `ForwardedForIndex`.

## Пример использования

```go
lookup := NewIPLookup()
ip := lookup.GetRemoteIP(req)
if ip != "" {
    ...
}
```

Использование пользовательской конфигурации:

```go
places := []string{"X-Real-IP", "RemoteAddr", "X-IP-Header"}
forwardedForIndex := 2
lookup := NewCustomIPLookup(places, forwardedForIndex)
ip := lookup.GetRemoteIP(req)
if ip != "" {
    ...
}
```
