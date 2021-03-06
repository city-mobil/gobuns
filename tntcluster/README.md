# Tntcluster

Содержит API для работы с Tarantool.

## Описание пакетов

### Pool

Содержит connection pool для работы с тарантулом. Стандартный коннектор использует одно соединение. Это может приводить
к проблемам.

<details><summary>Развернуть</summary>
<p>
Изначально стандартный коннектор устанавливает только одно соединение к тарантулу.

Предположим, что есть какой-то тяжёлый запрос, возвращающийся из тарантула.

Так как протокол тарантула асинхронный, то пока не тяжёлый ответ не будет записан в ожидающего его клиента, то все
остальные запросы могут привести к таймауту

Поэтому добавляется пул соединений, выбор соединения из которого осуществляется при помощи Round-Robin. Это позволяет
слегка уменьшить эффект тяжёлого запроса.
</p>
</details>

### tntclustercfg

Содержит конфигурацию для кластера и шардов.

#### OldClusterConfig

Создаёт конфигурацию для кластера в старом формате.

Старый формат отличается от нового отсутствием префиксов для шардов. Конфиг выглядит следующим образом:

```yaml
addrs = ["1", "2"];
slaves = [["1_slave_1", "1_slave_2"], ["2_slave_1", "2_slave_2"]]
```

#### NewClusterConfig

Инициализация конфига нового образца. Пример использования в yaml:

```yaml
auth:
  tntcluster:
    shard0:
      addr: 'auth1.ddk:3301'
      user: 'guest'
      password: 'guest'
    shard1:
      addr: 'auth2.ddk:3301'
      user: 'guest'
      password: 'guest'
```

#### Пример использования

```go
func main() {
    // Регистрируем конфигурационные параметры.
    getCfgFn := tntclustercfg.NewClusterConfig("prefix1")
    getOldCfgFn := tntclustercfg.OldClusterConfig("prefix2")

    // Чтение конфигурации из файла/переменных окружения/параметров командной строки
    err := config.InitOnce()
    if err != nil {
        log.Fatal(err)
    }

    // Создаём новый кластер
    cluster := tntcluster.NewCluster(nil, getCfgFn())
    cluster2 := tntcluster.NewCluster(nil, getOldCfgFn())
}
```

## Описание методов

### NewShard

Создаёт новый инстанс шарда для тарантула.

В случае, если тарантул одиночный(отсутствует шардирование или шардирование выполнено на стороне vshard), рекомендуется
использовать Shard вместо Cluster.

### NewShardWithTracing

Создаёт новый инстанс шарда для тарантула с включённым OpenTracing.

В случае, если тарантул одиночный (отсутствует шардирование или шардирование выполнено на стороне vshard), рекомендуется
использовать Shard вместо Cluster.

### NewCluster

Создаёт новый инстанс кластера для тарантула.

### NewClusterWithTracing

Создаёт новый инстанс кластера для тарантула с включённым OpenTracing.
