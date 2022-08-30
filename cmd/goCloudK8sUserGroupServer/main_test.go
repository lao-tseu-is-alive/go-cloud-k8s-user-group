package main

import (
	"fmt"
	"github.com/lao-tseu-is-alive/go-cloud-k8s-user-group/pkg/gohttpclient"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"os"
	"sync"
	"testing"
	"time"
)

func TestMainExecution(t *testing.T) {
	listenAddr := fmt.Sprintf("http://localhost:%d/", defaultPort)
	err := os.Setenv("PORT", fmt.Sprintf("%d", defaultPort))
	if err != nil {
		t.Errorf("Unable to set env variable PORT")
		return
	}
	// starting main in his own go routine
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		main()
	}()
	gohttpclient.WaitForHttpServer(listenAddr, 1*time.Second, 10)

	resp, err := http.Get(listenAddr)
	if err != nil {
		t.Fatalf("Cannot make http get: %v\n", err)
	}
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode, "Should return an http status ok")

	receivedHtml, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Error reading response body: %v\n", err)
	}

	// check that receivedJson contains the specified tt.wantBody substring . https://pkg.go.dev/github.com/stretchr/testify/assert#Contains
	assert.Contains(t, string(receivedHtml), "<html", "Response should contain the html tag.")
	//assert.Contains(t, string(receivedHtml), "\"request_id\":", "Response should contain the request_id field.")

}
