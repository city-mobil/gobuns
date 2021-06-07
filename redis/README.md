# Redis

Примеры вы можете найти в директориях `examples/cluster` и `examples/standalone`

# Конфигурация

```yaml
redis:
  // Адрес хоста редиса
  addr: '127.0.0.1:6379'
  // Используйте перечисление нод если вы используете подключение к кластеру
  #addr: '127.0.0.1:6379;127.0.0.2:6379'

  // Имя пользователя
  username: 'user'

  // Пароль пользователя
  password: 'secret'

  // Таймаут на подключение
  dial-timeout: 300ms

  // Таймаут на чтение
  read-timeout: 300ms

  // Таймаут на запись
  write-timeout: 300ms

  // Максимальное количество попыток
  max-retries: 3
  min-retry-backoff: 8ms
  max-retry-backoff: 512ms

  pool-size: 250
  min-idle-conns: 0
  pool-timeout: 100ms
  idle-timeout: 0
  idle-check-frequency: 0
  max-conn-age: 0

  // Настройки для трейсера
  tracer:
    // Добавлять ли hook трейсера для команд
    with-hook: true

  // Настройки для circuit breaker
  barber:
    max_fails: 50
    threshold: 42

  // Специфичные для кластера параметры
  cluster:
    // Максимальное число редиректов (см. MOVED/ASK)
    max-redirects: 8

    // Включить режим чтения на слейвах
    readonly: true

    // Позволяет направлять команды чтения к ближайшей ноде мастера или слейва
    route-by-latency: false

    // Позволяет направлять команды чтения к случайной ноде мастера или слейва
    route-randomly: true

  // Специфичные для стандалона параметры
  standalone:
    // Список реплик если таковые имеются
    replicas:
      - '127.0.0.2:6379'
      - '127.0.0.3:6379'
      - '127.0.0.4:6379'

    // Номер базы данных к которой нужно подключиться
    db: 0
```
