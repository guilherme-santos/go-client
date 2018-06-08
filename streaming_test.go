package ldclient

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestStreamProcessor_DoNotBlockInCase401(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer ts.Close()

	cfg := Config{
		StreamUri: ts.URL,
		Logger:    log.New(ioutil.Discard, "", 0),
	}

	sp := newStreamProcessor("key", cfg, nil)

	chanErr := make(chan error)
	go sp.subscribe(chanErr)

	select {
	case <-chanErr:
	case <-time.After(time.Second):
		t.Error("it was not expected to block")
	}
}
