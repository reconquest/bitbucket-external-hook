package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/url"

	"github.com/bndr/gopencils"
	"github.com/reconquest/hierr-go"
)

func request(
	method string, resource *gopencils.Resource, options ...interface{},
) (interface{}, error) {
	uri := resource.Url
	if len(options) > 0 {
		if query, ok := options[len(options)-1].(map[string]string); ok {
			queryValues := url.Values{}
			for key, value := range query {
				queryValues.Set(key, value)
			}
			uri += "?" + queryValues.Encode()

			options = options[:len(options)-1]
		}
	}

	var (
		request *gopencils.Resource
		err     error
	)

	// oh... gopencils doesn't have exported method Do()
	switch method {
	case "GET":
		request, err = resource.Get(options...)
	case "HEAD":
		request, err = resource.Head(options...)
	case "POST":
		request, err = resource.Post(options...)
	case "PUT":
		request, err = resource.Put(options...)
	case "DELETE":
		request, err = resource.Delete(options...)
	default:
		panic("unexpected method")
	}

	if err == nil || err == io.EOF {
		err = extractRequestError(request, err)
	}

	if err != nil {
		return nil, hierr.Errorf(err, "can't request %s %s", method, uri)
	}

	return request.Response, nil
}

func extractRequestError(request *gopencils.Resource, err error) error {
	if request.Raw.StatusCode < 400 {
		return nil
	}

	unexpectedStatusCode := "remote server responds with unexpected status " +
		"(" + request.Raw.Status + ")"

	responseBody, err := ioutil.ReadAll(request.Raw.Body)
	if err != nil {
		return hierr.Errorf(
			hierr.Errorf(err, "can't read response body"),
			unexpectedStatusCode,
		)
	}

	var responseError ResponseError

	err = json.Unmarshal(responseBody, &responseError)

	if len(responseError.Errors) > 0 {
		return hierr.Errorf(
			responseError.String(),
			unexpectedStatusCode,
		)
	}

	if err != nil {
		return hierr.Errorf(
			hierr.Errorf(
				string(responseBody),
				"can't decode JSON",
			),
			unexpectedStatusCode,
		)
	}

	// response is JSON, but not ResponseError (or empty message with
	// Java exception), should return as is
	return hierr.Errorf(string(responseBody), unexpectedStatusCode)
}
