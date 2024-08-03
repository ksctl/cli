## here is the kt for the profiling of the code cli

```go
package main

import (
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"sync"

	"github.com/ksctl/cli/cli/cmd"
)

func main() {
	cc := make(chan os.Signal, 1)

	signal.Notify(cc, os.Interrupt)
	signal.Notify(cc, os.Kill)

	wg := new(sync.WaitGroup)

	go func() {
		wg.Add(1)
		fmt.Println(http.ListenAndServe("localhost:8080", nil))
	}()

	cmd.Execute()

	select {
	case sig := <-cc:
		log.Println("Received terminate, graceful shutdown", sig)
		os.Exit(1)
	}

	wg.Wait()
}

```

```shell
go run . create local -n demo --yes
```

```shell
go tool pprof \
  -raw -output=cpu.txt \
  'http://localhost:8080/debug/pprof/profile?seconds=120'
```

```shell
./stackcollapse-go.pl cpu.txt | flamegraph.pl > flame.svg
```

<https://github.com/brendangregg/FlameGraph>
