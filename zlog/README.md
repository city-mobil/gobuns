# JSON logger based on Zerolog

Пакет предоставляет возможность логирования данных в формате JSON, используя
пакет [zerolog](https://github.com/rs/zerolog).

API пакета `zlog` практически полностью повторяет API пакета `zerolog`. Основные отличия:

* Пакет `glog` заменяет пакет `zerolog/log`,
* Пакет `hlog` заменяет пакет `zerolog/hlog`,
* Пакет `zwatch` позволяет настроить обновление глобального уровня логирования через KV хранилище Consul,
* Часть функций `zerolog` может отсутствовать (например, `hooks`, пользовательский формат логирования объектов и
  некоторые другие).

Далее приведены базовые примеры использования логера. Больше примеров и подробную документацию смотрите в описании
к `zerolog`.

## Инициализация логера

```go
// Создание простого логера без каких-либо полей в контексте.
logger := zlog.Raw(os.Stdout)
// Создание логера с timestamp полем.
logger = zlog.New(os.Stdout)
// Создание отключенного логера.
logger = zlog.Nop()
```

## Пример использования глобального логера

```go
package main

import (
    "github.com/city-mobil/gobuns/zlog"
    "github.com/city-mobil/gobuns/zlog/glog"
)

func main() {
    zlog.SetTimeFieldFormat(zlog.TimeFormatUnix)
    glog.Print("hello world")
}

// Output: {"time":1516134303,"level":"debug","message":"hello world"}
```

По умолчанию глобальный логер пишет в `stdout` с уровнем логирования Debug.

## Логирование с контекстом

```go
package main

import (
    "github.com/city-mobil/gobuns/zlog"
    "github.com/city-mobil/gobuns/zlog/glog"
)

func main() {
    zlog.SetTimeFieldFormat(zlog.TimeFormatUnix)

    glog.Debug().
        Str("Scale", "833 cents").
        Float64("Interval", 833.09).
        Msg("Fibonacci is everywhere")
    
    glog.Debug().
        Str("Name", "Tom").
        Send()
}

// Output: {"level":"debug","Scale":"833 cents","Interval":833.09,"time":1562212768,"message":"Fibonacci is everywhere"}
// Output: {"level":"debug","Name":"Tom","time":1562212768}
```

## Интеграция с http.Handler

Пакет `hlog` предоставляет набор вспомогательных методов для интеграции логера в `http.Handler`.

```go
package main

import (
	"net/http"
	"os"

	"github.com/city-mobil/gobuns/zlog"
	"github.com/city-mobil/gobuns/zlog/hlog"
)

func final(w http.ResponseWriter, r *http.Request) {
	log := hlog.FromRequest(r)
	log.Info().Str("status", "ok").Msg("request finished")

	_, _ = w.Write([]byte("OK"))
}

func main() {
	log := zlog.New(os.Stdout)
	mux := http.NewServeMux()

	fh := http.HandlerFunc(final)
	mw := hlog.RequestIDHandler("request_id", hlog.JaegerTraceHeaderName, true)(
		hlog.RemoteAddrHandler("remote")(
			hlog.MethodHandler("method")(fh),
		),
	)

	mux.Handle("/", hlog.NewHandler(log)(mw))

	if err := http.ListenAndServe(":8080", mux); err != http.ErrServerClosed {
		log.Fatal().Err(err)
	}
}

// Output: {"level":"info","request_id":"514bbe5bb5251c92bd07a9846f4a1ab6","method":"GET","status":"ok","time":"2020-12-14T21:22:56+03:00","message":"request finished"}
```

В `RequestIDHandler` можно настроить автоматическую генерацию уникального идентификатора запроса, если он не передан в
HTTP заголовке, через аргумент `generateReqID`.

## Изменение уровня логирования через Consul

```go
package main

import (
    "github.com/city-mobil/gobuns/registry"
    "github.com/city-mobil/gobuns/zlog"
    "github.com/city-mobil/gobuns/zlog/glog"
    "github.com/city-mobil/gobuns/zlog/zwatch"
)

func main() {
    zlog.SetTimeFieldFormat(zlog.TimeFormatUnix)
	
    wc := registry.WatchConfig{
        Addr: "localhost:8500",
    }
    cancel, err := zwatch.GlobalLevel("logger/level", wc)
    if err != nil {
        glog.Fatal().Msg("could not subscribe on key in Consul")
    }
    defer cancel()
    
    glog.Debug().
        Str("Name", "Tom").
        Send()
    
    // Update the key, e.g. set the level to "fatal"
    glog.Debug().Str("Name", "Bob").Send()
}

// Output: {"level":"debug","Name":"Tom","time":1562212768}
```

Обратите внимание, что нельзя установить уровень логирования ниже уровня, с которым этот логер был инициализирован:

```go
// Create a new logger and set the minimum level to Error.
log := zlog.New(os.Stdout).Level(zlog.ErrorLevel)
// Update the global level.
zlog.SetGlobalLevel(zlog.DebugLevel)
// This will not produce any output because the logger minimum level is Error.
log.Debug().Msg("debug message")
// This works fine.
log.Error().Msg("error message")
// Bump the global level to Fatal.
zlog.SetGlobalLevel(zlog.FatalLevel)
// The global level is Fatal so this message will not be printed.
log.Error().Msg("error message") 
```
