package testutils

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/valyala/fasthttp/fasthttputil"
	"gitlab.com/vedhabhavanam/smarthome/dwarka/pkg/api"
	"gitlab.com/vedhabhavanam/smarthome/dwarka/pkg/store"
	"io/ioutil"
	"net"
	"net/http"
)

// ServeHTTPRequest serves http request using provided fasthttp handler
func ServeHTTPRequest(store store.Store, req *http.Request) (*http.Response, error) {
	ln := fasthttputil.NewInmemoryListener()
	defer func() {
		_ = ln.Close()
	}()

	go func() {
		httpServer := api.NewServer("", "", store)
		err := httpServer.Serve(ln)
		if err != nil {
			panic(fmt.Errorf("failed to ServeHTTPRequest: %v", err))
		}
	}()

	client := http.Client{
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				return ln.Dial()
			},
		},
	}

	return client.Do(req)
}

// Read unmarshal response.Body into the type provided as input
func Read(response *http.Response, in interface{}) error {
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, in)
}

// ReadError unmarshal response.Body into the type provided as input
func ReadError(response *http.Response) (string, error) {
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	res := map[string]string{}
	err = json.Unmarshal(data, &res)
	if err != nil {
		return "", err
	}

	value, ok := res["error"]
	if !ok {
		return "", fmt.Errorf("no error res response")
	}
	return value, nil
}
