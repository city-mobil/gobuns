# Config

Пакет для работы с конфигурацией приложения. Конфигурация может быть передана через файл в формате yaml или json,
аргументы командной строки.

Аргументы командной строки имеют приоритет над значениями из файла.

Имя файла конфигурации передается через аргумент `--config` или переменную окружения `CONFIG`.

## Пример вызова приложения

```bash
# command line argument
./app --config=/opt/conf.yaml --app.hostname="localhost"
# env variable
CONFIG=/opt/conf.yaml ./app
```

## Пример объявления флагов конфигурации

```go
package main

import (
	"log"

	"github.com/city-mobil/gobuns/config"
)

var (
	appAddr     = config.String("app.addr", ":8080", "TCP address for the service to listen on")
	appHostname = config.String("app.hostname", "localhost", "Application hostname")
	appEnv      = config.String("app.env", "dev", "Application environment")
)

func main() {
	err := config.InitOnce()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("addr: %s", *appAddr)
	log.Printf("hostname: %s", *appHostname)
	log.Printf("env: %s", *appEnv)
}
```

## Пример объявления структуры конфигурации

```go
package app

import (
	"strings"

	"github.com/city-mobil/gobuns/config"
)

const (
	defAddr     = ":8080"
	defHostname = "localhost"
)

type Config struct {
	Addr     string
	Hostname string
}

func NewConfig(prefix string) func() *Config {
	pf := prefix
	if !strings.HasSuffix(prefix, ".") {
		pf += "."
	}

	var (
		addr     = config.String(pf+"addr", defAddr, "TCP address for the service to listen on")
		hostname = config.String(pf+"hostname", defHostname, "Application hostname")
	)

	return func() *Config {
		return &Config{
			Addr:     *addr,
			Hostname: *hostname,
		}
	}
}
```

Использование в main.go:

```go
package main

import (
	"log"

	"github.com/city-mobil/gobuns/config"
)

func main() {
	getAppConfigFn := app.NewConfig("app")

	err := config.InitOnce()
	if err != nil {
		log.Fatal(err)
	}

	appCfg := getAppConfigFn()

	log.Printf("addr: %s", appCfg.Addr)
	log.Printf("hostname: %s", appCfg.Hostname)
}
```

Такие структуры можно создавать для каждой изолированной области приложения, либо создать единую конфигурацию для всего
сервиса.

Можно комбинировать объявление глобальных флагов и структур конфигурации.

## Особенности использования

Методы объявления флагов аналогичны методам из пакета `flag` стандартной библиотеки.

## InitOnce

Метод для инициализации конфига приложения. Возвращаемое значение `error` обозначает наличие ошибки во время
инициализации конфигов.

Все флаги конфигурации должны быть объявлены строго до вызова метода `InitOnce`.

Пример вызова

```go
err := config.InitOnce()
if err != nil {
    log.Fatal(err)
}
```

## StringSlice

Флаг типа `StringSlice` ожидает значения разделенные запятой, при этом разделитель будет учитываться и в самих
значениях.

```go
[]string{"val1, val2", "val3"} -> []string{"val1", "val2", "val3"}
```
