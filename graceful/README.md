# Graceful

Пакет для работы с [graceful-shutdown](https://whatis.techtarget.com/definition/graceful-shutdown-and-hard-shutdown) (
изящным завершением)

## Пример использования

```go
srv := &http.Server{}

// Настраиваем свой обработчик ошибок для graceful shutdown.
// Будет вызываться каждый раз, когда callback возвращает ошибку.
graceful.ExecOnError(func(err error) {
	glog.Err(err)
})
// Добавляем обработчик, который будет 
// вызван в процессе завершения программы.
graceful.AddCallback(srv.Close)
// Мы можем добавить столько обработчиков, сколько потребуется.
// Обработчики вызываются последовательно в порядке обратном их добавлению.

go func() {
    err := srv.ListenAndServe()
    if err != nil {
    	// Принудительно завершаем выполнение программы
    	// при получении фатальной ошибки.
        graceful.ShutdownNow()	
    }
}()

// Ожидание сигнала ОС для завершения программы. 
graceful.WaitShutdown()
```
