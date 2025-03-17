package integrations

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
)

var URL url.URL

func SendRequest(method, path string, body interface{}) (status int, content []byte) {
	URL.RawQuery = ""
	parts := strings.Split(path, "?")

	if len(parts) >= 2 {
		URL.RawQuery = strings.Join(parts[1:], "?")
	}

	URL.Path = parts[0]

	data, err := json.Marshal(body)
	if err != nil {
		smartPanic(err)
	}

	if method == http.MethodGet || method == http.MethodHead {
		data = nil
	}

	request, err := http.NewRequestWithContext(context.Background(), method, URL.String(), bytes.NewReader(data))
	if err != nil {
		smartPanic(err)
	}

	request.Header.Set("Content-Type", "application/json; charset=utf-8")
	request.Header.Set("Accept", "*/*")

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		smartPanic(err)
	}

	defer response.Body.Close()

	return ExtractResponse(response)
}

func ExtractResponse(response *http.Response) (status int, content []byte) {
	status = response.StatusCode
	contentBytes, err := io.ReadAll(response.Body)

	if err != nil {
		smartPanic(err)
	}

	return status, contentBytes
}
