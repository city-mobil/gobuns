# Prometheus metrics

Библиотека-обертка для работы с метриками prometheus.

# Основные метрики

## Event

Метрика-событие. Аналогично `RecordCustomEvent` в NR.

```go
event := &promlib.Event{
    Name:      "retry_count",
    Namespace: "myapp",
    Subsystem: "mysql",
    Help:      "Total number of SQL query retries",
}

promlib.IncCntEvent(event)
promlib.AddCntEvent(event, 2.0)
```

Либо регистрируем новое событие с дополнительной мета-информацией.

```go
event := &promlib.Event{
    Name:      "retry_count",
    Namespace: "myapp",
    Subsystem: "mysql",
    Help:      "Total number of SQL query retries",
}

promlib.IncCntEventWithLabels(event, promlib.Labels{
    "query": "SelectAll",
})
```

### Вспомогательные методы

#### IncCnt

Увеличиваем счётчик для конкретного события. Если счётчик не существует, он будет создан.

```go
promlib.IncCnt("some_event")
```

#### IncCntWithLabels

Увеличиваем счётчик для конкретного события с указанными лейблами. Если счётчик не существует, он будет создан.

```go
promlib.IncCntWithLabels("some_event", promlib.Labels{
    "a": "a",
    "b": "b",
    "c": "c",
    "d": "d",
})
```

#### AddCnt

Увеличиваем счётчик на заданное значение для конкретного события. Если счётчик не существует, он будет создан.

```go
promlib.AddCnt("some_event", 2.0)
```

#### AddCntWithLabels

Увеличиваем счётчик на заданное значение для конкретного события с указанными лейблами. Если счётчик не существует, он
будет создан.

```go
promlib.AddCntWithLabels("some_event", promlib.Labels{
    "a": "a",
    "b": "b",
    "c": "c",
    "d": "d",
}, 2.0)
```

#### SetGlobalNamespace

Устанавливает глобальный namespace, который будет использоваться по умолчанию для метрик.

```go
promlib.SetGlobalNamespace("my_app")
```

## Transaction

Транзакции поддерживают метрики вида Histogram и Summary.

Данную метрику удобно использовать для замера времени исполнения транзакции, в том числе и для походов в базы данных.

```go
opts := &promlib.HistogramOpts{
    MetaOpts: MetaOpts{
        Name:      "query_time",
        Subsystem: "redis",
        Help:      "Redis queries response time",
    },
    Labels:  []string{"addr", "query"},
    Buckets: []float64{.002, .005, .01, .015, .025, .05, .1, .25, .5, 1, 2, 10},
}
txn := promlib.NewTransaction(opts)
txn.Start("localhost:8000", "GetMyKey")
...
txn.End()
```

Метрику можно использовать для замера времени исполнения асинхронных процессов:

```go
opts := &promlib.HistogramOpts{
    MetaOpts: MetaOpts{
        Name:      "cross_time",
        Subsystem: "order",
        Help:      "Cross time communication",
    },
    Labels:  []string{"process", "system"},
    Buckets: []float64{.002, .005, .01, .015, .025, .05, .1, .25, .5, 1, 2, 10},
}
txn := promlib.NewTransaction(opts)
txn.Observe(11.4, "bill", "finance")
```

Библиотека предоставляет набор подготовленных конфигураций для транзакций (смотри [metrics](metrics.go)). Пример
использования:

```go
func (r *repo) exec(ctx context.Context, query *tarantool.Call) (*tarantool.Result, error) {
	shard, err := r.cluster.ChooseShard(tntcluster.ChooseFirstShard())
	if err != nil {
		return nil, err
	}

	txn := promlib.NewTransaction(promlib.TarantoolOpTime)
	txn.Start("client_auth", query.Name)
	defer txn.End()

	return shard.CallTolerant(ctx, query)
}
```

# HTTPMiddleware

Для того чтобы автоматически собирать метрики по всем HTTP запросам с разбивкой по путям, методам (GET, POST...)
нужно использовать HTTPMiddleware.

```go
promMW := promlib.NewMiddleware(promlib.DefHTTPRequestDurBuckets)
router := mux.NewRouter()
router.Use(promMW.Handler)
```

**Важно**: если URL запросы строятся динамически, то необходимо настроить свою функцию формирования значения для
метки `path`, иначе метрика будет иметь высокое значение cardinality, что плохо скажется на производительности
Prometheus.

Настраивается через опцию `promlib.WithCustomPath`:

```go
normalizer := func (req *http.Request) string {
    return r.URL.Path
}
mw := promlib.NewMiddleware(promlib.DefHTTPRequestDurBuckets, promlib.WithCustomPath(normalizer))
```

Если динамическая часть URL состоит только из чисел, то используйте `promlib.URLNumberNormalizer`:

```go
normalizer := func (req *http.Request) string {
	return promlib.URLNumberNormalizer(req, []string{":id", ":uid"})
}
```

Также доступны нормализаторы чисел в системе HEX `URLHEXNumberNormalizer` и по порядковому
индексу `URLIndicesNormalizer`.

# Отдача метрик

Следующий пример демонстрирует каким образом можно отдавать все метрики, которые собирает сервис.

```go
router := mux.NewRouter()
router.Handle("/metrics", promhttp.Handler())
```

# InstrumentRoundTripper

Используйте эту middleware для снятия метрик HTTP клиента.

```go
client := &http.Client{}
client.Timeout = 1 * time.Second
client.Transport = promlib.InstrumentRoundTripper("ask_google", http.DefaultTransport)
```
