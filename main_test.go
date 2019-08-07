package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"sync"
	"syscall"
	"testing"
	"time"

	"github.com/phayes/freeport"
)

func TestHelloWorld(t *testing.T) {
	port, err := freeport.GetFreePort()
	if err != nil {
		t.Fatal(err)
	}

	sigs := make(chan os.Signal, 1)
	var closeWG sync.WaitGroup
	go Run([]string{"hello-world", "--port", strconv.Itoa(port)}, sigs, &closeWG)

	cleanup := func() {
		sigs <- syscall.SIGTERM
		closeWG.Wait()
	}
	defer cleanup()

	// Wait for http server
	time.Sleep(time.Millisecond * 200)

	resp, err := http.Get(fmt.Sprintf("http://127.0.0.1:%d", port))
	if err != nil {
		t.Fatal(err)
	}

	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(b, []byte("assalamualaikum, dunia")) {
		t.Fatalf("expecting response 'assalamualaikum, dunia', got '%s' instead", string(b))
	}
}
