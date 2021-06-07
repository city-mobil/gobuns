# retry

Предоставляет конфигурируемое поведение выполнения повторных запросов в случае ошибок.

```go
cfg := retry.NewDefRetryConfig()
retrier := retry.New(cfg)

action := func() error {
    return retry.Unrecoverable(errors.New("fatal error"))
}

onRetry := func(n uint, err error) {
    logger.Error(fmt.Sprintf("attempt: %d, err: %v", n, err))
}

err := retrier.Do(context.Background(), action, onRetry)
```

Предоставляет возможность получения интервалов времени согласно указанной стратегии.

```go
cfg := retry.NewDefWaitConfig()
waiter := retry.NewWaiter(cfg)

time := waiter.Get(iter)
```

## Конфигурирование

```go
import (
    "github.com/city-mobil/gobuns/config"
    "github.com/city-mobil/gobuns/retry"
)

getRetryCfgFn := retry.GetRetryConfig("my_service")

defWait := 10 * time.Millisecond
getWaitCfgFn := retry.GetWaitConfig("my_service", defWait)

err := config.InitOnce()
if err != nil {
    panic(err)
}

retrier := retry.New(getRetryCfgFn())
waiter := retry.NewWaiter(getWaitCfgFn())
```

## Стратегии ожидания

Для всех стратегий можно задать значение `MaxWait`, которое ограничивает время ожидания независимо от других настроек.

### Fixed

Стратегия использует константное время ожидания.

Пример конфигурации:

```yaml
# Ожидаем 25ms перед повторным выполнением запроса.
retries: 
  wait_type: 'fixed'
  base_wait: 25ms
```

Отключить ожидание и выполнять повторный запрос сразу:

```yaml
retries:
  wait_type: 'fixed'
  base_wait: 0
```

### Random

Стратегия использует случайное время ожидания, но не более `MaxJitter`.

Пример конфигурации:

```yaml
# Ожидаем случайное время перед повторным выполнением запроса, но не более 100ms.
retries: 
  wait_type: 'random'
  max_jitter: 100ms
```

### Backoff

Стратегия увеличивает время ожидания перед каждым повтором.

Пример конфигурации:

```yaml
# 0: 25ms
# 1: 50ms
# 2: 100ms
# 3: 200ms
# ....
retries: 
  wait_type: 'backoff'
  base_wait: 25ms
```

### Combine

Стратегия объединяет стратегии Random и Backoff.

Пример конфигурации:

```yaml
# delay := random(max_jitter) + backoff(base_wait, attempt)
retries: 
  wait_type: 'combine'
  base_wait: 25ms
  max_jitter: 100ms
```