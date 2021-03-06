package helpers

import (
	"bytes"
	"io/ioutil"
	"net/http"
)

type Headers map[string]string

func Post(endpoint string, headers Headers, body []byte) (int, []byte, error) {
	return Request(http.MethodPost, endpoint, headers, body)

}

func Get(endpoint string, headers Headers) (int, []byte, error) {
	return Request(http.MethodGet, endpoint, headers, []byte{})
}

func Request(method string, endpoint string, headers Headers, body []byte) (int, []byte, error) {
	request, err := http.NewRequest(method, endpoint, bytes.NewReader(body))

	if err != nil {
		return 0, nil, err
	}

	for header, value := range headers {
		request.Header.Set(header, value)
	}

	response, err := http.DefaultClient.Do(request)

	if err != nil {
		return 0, nil, err
	}

	return extractCodeAndBodyFromResponse(response)
}

func extractCodeAndBodyFromResponse(response *http.Response) (int, []byte, error) {
	defer response.Body.Close()

	responseBody, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return 0, nil, err
	}

	return response.StatusCode, responseBody, nil
}
