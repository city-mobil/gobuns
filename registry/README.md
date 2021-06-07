# Consul KV

Пакет для работы с KV хранилищем Consul.

## Инциализация клиента

```go
config := registry.Config{
    Addr:         "localhost:3000",
    QueryTimeout: 1 * time.Second,
    RetryConfig:  retry.NewDefRetryConfig(),
}
client, err := registry.NewClient(config)
```

Более подробно о стратегиях повторных запросов смотрите в [документации](../retry/README.md).

Пример использования YAML конфигурации:

```yaml
app:
  registry:
    addr: localhost:8500
    query_timeout: 1s
    retries:
      wait_type: 'backoff'
      base_wait: 25ms
      max_attempts: 3
```

```go
getRegistryCfgFn := registry.NewConfig("app")

err := config.InitOnce()
if err != nil {
    panic(err)
} 

config := getRegistryCfgFn()
client, err := registry.NewClient(config)
```

### Чтение значения из KV по ключу

```go
key := "service/key"
v, err := client.GetString(context.Background(), key)
if err != nil {
    panic(err)
}
```

### Подписка на изменения значения по ключу

```go
key := "service/key"
hd := registry.WatchHandleFunc(func(data *string) {
    if data == nil {
        log.Print("key has been removed")        
    } else {
        log.Printf("got updated value: %s", *data)
    }
})
wc := registry.WatchConfig{
    Addr: "localhost:8500",
    OnErr: func(err error) {
        log.Fatalf("failed to watch: %w", err)
    },
}

cancel, err := registry.Watch(key, hd, wc)
if err != nil {
    panic(err)
}
defer cancel()
```
