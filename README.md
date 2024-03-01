# cronticker

`cronticker` implements a ticker that works with crontab schedules.

## Install

```
go get github.com/martin-viggiano/cronticker
```

## Usage

Go doc: [https://pkg.go.dev/github.com/martin-viggiano/cronticker]()

```go
import (
	"log"

	"github.com/martin-viggiano/cronticker"
)

func main() {
	ticker, err := cronticker.NewTicker("5 4 * * *")
	if err != nil {
		log.Fatal("failed to create ticker")
	}
	defer ticker.Stop()

	select {
	case <-ticker.C:
		log.Print("it is 04:05")
	}

	err = ticker.Reset("5 4 * * *")
}
```

## Useful links

- [cron Wikipedia page](https://en.wikipedia.org/wiki/Cron)
- [robfig/cron pkg](https://github.com/robfig/cron)