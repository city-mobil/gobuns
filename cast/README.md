# Cast

Пакет для работы с преобразованием типов.

## Список доступных функций

- ParseInt
- ParseUint*
- ParseInt8
- ParseUint8*
- ParseInt16
- ParseUint16*
- ParseInt32
- ParseUint32*
- ParseInt64
- ParseUint64*
- ParseBool
- ParseString
- ParseDuration
- ParseStringSlice
- ParseIntSlice
- ParseBoolSlice

В случае отрицательного значения возвращается ошибка. Не учитывается переполнение значения.

## Scan

Данная функция копирует значения из `tuple` в `dest`. Длина среза `tuple` и количество значений `dest` должно быть
равным.

### Пример

```go
var (
    str string
    num int
)

err := cast.Scan([]interface{}{"string_value", 12}, &str, &num)
if err != nil {  
    panic(err)  
}
```