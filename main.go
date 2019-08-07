package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/go-chi/chi"
)

func main() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	os.Exit(Run(os.Args, sigs, nil))
}

func Run(args []string, sigs <-chan os.Signal, closeWG *sync.WaitGroup) int {
	if closeWG != nil {
		closeWG.Add(1)
		defer closeWG.Done()
	}

	if len(args) < 1 {
		fmt.Println("usage: helloworld [--port <port>]")
		return 1
	}

	var port int
	fs := flag.NewFlagSet("helloworld", flag.ExitOnError)
	fs.IntVar(&port, "port", 3000, "http port")
	if err := fs.Parse(args[1:]); err != nil {
		fmt.Println(err)
		return 1
	}

	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		fmt.Println(err)
		return 1
	}

	defer ln.Close()

	r := chi.NewRouter()
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("assalamualaikum, dunia"))
	})

	fmt.Printf("listening on http://127.0.0.1:%d\n", port)

	server := &http.Server{Handler: r}
	go func() {
		if err := server.Serve(ln); err != nil && err != http.ErrServerClosed {
			fmt.Println(err)
		}
	}()

	sig := <-sigs
	fmt.Println("received signal", sig)

	clsCtx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()
	server.Shutdown(clsCtx)

	return 0
}
