# MySQL

Пакет для работы с MySQL

## Описания интерфейсов

### Shard

Представляет собой одиночный шард с мастер-соединением и реплика-соединением.

#### Описание методов

##### GetMasterConn

Возвращает соединение к мастеру.

##### GetSlaveConn

Возвращает соединение к реплике.

Реплика выбирается согласно <code>ReplicaStrategy</code> в конфигурации.

##### ExecMaster

Выполняет <code>Exec</code> на мастере.

Для большей информации, посмотрите [официальную документацию](https://golang.org/pkg/database/sql/#DB.Exec)

##### ExecMasterTolerant

Выполняет `Exec` на мастере с использованием ретраев в случае ошибок. Смотри также [ExecMaster](#ExecMaster).

##### ExecReplica

Выполняет <code>Exec</code> на реплике.

Реплика выбирается согласно <code>ReplicaStrategy</code> в конфигурации.

Для большей информации, посмотрите [официальную документацию](https://golang.org/pkg/database/sql/#DB.Exec)

##### ExecReplicaTolerant

Выполняет `Exec` на реплике с использованием ретраев в случае ошибок. Смотри также [ExecReplica](#ExecReplica).

##### QueryMaster

Выполняет <code>Query</code> на мастере.

Для большей информации, посмотрите [официальную документацию](https://golang.org/pkg/database/sql/#DB.Query)

##### QueryMasterTolerant

Выполняет `Query` на мастере с использованием ретраев в случае ошибок. Смотри также [QueryMaster](#QueryMaster).

##### QueryReplica

Выполняет <code>Query</code> на реплике.

Реплика выбирается согласно <code>ReplicaStrategy</code> в конфигурации.

Для большей информации, посмотрите [официальную документацию](https://golang.org/pkg/database/sql/#DB.Query)

##### QueryReplicaTolerant

Выполняет `Query` на реплике с использованием ретраев в случае ошибок. Смотри также [QueryReplica](#QueryReplica).

##### QueryRowMaster

Выполняет <code>QueryRow</code> на мастере.

Для большей информации, посмотрите [официальную документацию](https://golang.org/pkg/database/sql/#DB.QueryRow)

##### QueryRowMasterTolerant

Выполняет `QueryRow` на мастере с использованием ретраев в случае ошибок. Смотри также [QueryRowMaster](#QueryRowMaster)
.

##### QueryRowReplica

Выполняет `QueryRow` на реплике.

Реплика выбирается согласно `ReplicaStrategy` в конфигурации.

Для большей информации, посмотрите [официальную документацию](https://golang.org/pkg/database/sql/#DB.QueryRow)

##### QueryRowReplicaTolerant

Выполняет `QueryRow` на реплике с использованием ретраев в случае ошибок. Смотри
также [QueryRowReplica](#QueryRowReplica).

##### BarberStats

Возвращает статистику Circuit Breaker"а.

##### Close

Завершает все соединения(к мастеру и к репликам), которые открыты.

##### Setup

Метод, который устанавливает соединения к мастеру и ко всем репликам из конфигурации.

Этот метод *ОБЯЗАТЕЛЬНО* должен быть вызван перед началом работы с Shard.

## Описания методов

## Примеры использования

### Использование внутри логики

```go
type Connector interface {
    GetStuff(context.Context)
}

type defaultConnector struct {
    client mysql.Shard
}

func (d *defaultConnector) GetStuff(_ context.Context) {

}

func NewConnector(config *mysqlconfig.ShardConfig) Connector {
    cirulink := barber.NewBarber([]int{1}, config.CircuitBreakerConfig)
    client := mysql.NewShard(config, mysqlconfig.ClusterConnectorTypeSQL, cirulink)
    if err := client.Setup(); err != nil {
        log.Fatal(err)
    }
    return &defaultConnector{
        client: client,
    }
}
```

## Инициализация конфигурации

```go
import (
    "github.com/city-mobil/gobuns/mysql"
    "github.com/city-mobil/gobuns/mysql/mysqlconfig"
)

sqlMasterCfgFn := mysqlconfig.NewDatabaseConfig("mysql.master")
sqlSlaveCfgFn := mysqlconfig.NewDatabaseConfig("mysql.slave")
sqlRetryCfgFn := mysqlconfig.NewRetryConfig("mysql.retries")
// Или используйте базовую конфигурацию
// sqlRetryCfg := mysqlconfig.NewDefaultRetryConfig()

err := config.InitOnce()
if err != nil {
    panic(err)
}

shardCfg := &mysqlconfig.ShardConfig{
		MasterConfig:    sqlMasterCfgFn(),
		SlaveConfigs:    []*mysqlconfig.DatabaseConfig{sqlSlaveCfgFn()}, // ожидаем HAproxy
		RetryConfig:     sqlRetryCfgFn(), // или sqlRetryCfg
		ReplicaStrategy: mysqlconfig.ReplicaStrategyRoundRobin,
}

shard := mysql.NewShard(shardCfg, mysqlconfig.ClusterConnectorTypeSQL)
```

Пример yml конфигурации:

```yml
mysql:
  master:
    addr: '127.0.0.1:3306'
    dbname: 'db'
    user: 'user'
    password: 'pass'
    timeout: 1s
    read_timeout: 1s
    write_timeout: 1s
    pool:
      max_open_connections: 200
      max_idle_connections: 200
      max_life_time: 0
  slave:
    addr: '127.0.0.1:3306'
    dbname: 'db'
    user: 'user'
    password: 'pass'
  retries:
    max: 5
    timeout: 100ms
```