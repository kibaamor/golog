# golog

A simple logger for go.

## Install

```bash
go get github.com/kibaamor/golog
```

## Example

[example.go](example/example.go)

```go
func main() {
    // basic logger
    logger := golog.NewTermLogger(os.Stderr, true)

    // got: `[INFO] 1:1 k1:v1 k2:[1 1]`
    _ = logger.Log(context.Background(), golog.LevelInfo, 1, 1, "k1", "v1", "k2", []int{1, 1})

    // combine multiple logger
    // Discard is logger with discard everything
    logger = golog.MultiLogger(logger, golog.Discard)

    // got: `[INFO] 1:1 k1:v1 k2:[1 1]`
    _ = logger.Log(context.Background(), golog.LevelInfo, 1, 1, "k1", "v1", "k2", []int{1, 1})

    // filter with log level
    logger = golog.WithFilter(logger, golog.FilterLevel(golog.LevelWarn))
    // got: ``
    _ = logger.Log(context.Background(), golog.LevelInfo, 1, 1)

    // auto add timestamp and caller information
    logger = golog.WithHandler(logger, golog.HandlerDefaultTimestamp, golog.HandlerDefaultCaller)

    // got:`[2021-12-10 12:33:26.968][example.go:24][WARN] 1:1`
    _ = logger.Log(context.Background(), golog.LevelWarn, 1, 1)

    // Helper provides useful apis, such as Info, Infow.
    helper := golog.NewHelper(logger)

    // got: `[2021-12-10 12:37:52.699][helper.go:76][ERROR] golog: hi`
    helper.Errorf("golog: %v", "hi")
}
```
