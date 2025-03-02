// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package flyscrape_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/philippta/flyscrape"
	"github.com/stretchr/testify/require"
)

func TestJSLibParse(t *testing.T) {
	script := `
    import { parse } from "flyscrape"

    const doc = parse('<div class=foo>Hello world</div>')
    export const text = doc.find(".foo").text()
    `

	client := &http.Client{
		Transport: flyscrape.MockTransport(200, html),
	}

	exports, err := flyscrape.Compile(script, flyscrape.NewJSLibrary(client))
	require.NoError(t, err)

	h, ok := exports["text"].(string)
	require.True(t, ok)
	require.Equal(t, "Hello world", h)
}

func TestJSLibHTTPGet(t *testing.T) {
	script := `
    import http from "flyscrape/http"

    const res = http.get("https://example.com")

    export const body = res.body;
    export const status = res.status;
    export const error = res.error;
    export const headers = res.headers;
    `

	client := &http.Client{
		Transport: flyscrape.MockTransport(200, html),
	}

	exports, err := flyscrape.Compile(script, flyscrape.NewJSLibrary(client))
	require.NoError(t, err)

	body, ok := exports["body"].(string)
	require.True(t, ok)
	require.Equal(t, html, body)

	status, ok := exports["status"].(int64)
	require.True(t, ok)
	require.Equal(t, int64(200), status)

	error, ok := exports["error"].(string)
	require.True(t, ok)
	require.Equal(t, "", error)

	headers, ok := exports["headers"].(map[string]any)
	require.True(t, ok)
	require.NotEmpty(t, headers)
}

func TestJSLibHTTPPostForm(t *testing.T) {
	script := `
    import http from "flyscrape/http"

    const res = http.postForm("https://example.com", {
        username: "foo",
        password: "bar",
        arr: [1,2,3],
    })

    export const body = res.body;
    export const status = res.status;
    export const error = res.error;
    export const headers = res.headers;
    `

	client := &http.Client{
		Transport: flyscrape.RoundTripFunc(func(r *http.Request) (*http.Response, error) {
			require.Equal(t, "POST", r.Method)
			require.Equal(t, "application/x-www-form-urlencoded", r.Header.Get("Content-Type"))
			require.Equal(t, "foo", r.FormValue("username"))
			require.Equal(t, "bar", r.FormValue("password"))
			require.Len(t, r.Form["arr"], 3)

			return flyscrape.MockResponse(400, "Bad Request")
		}),
	}

	exports, err := flyscrape.Compile(script, flyscrape.NewJSLibrary(client))
	require.NoError(t, err)

	body, ok := exports["body"].(string)
	require.True(t, ok)
	require.Equal(t, "Bad Request", body)

	status, ok := exports["status"].(int64)
	require.True(t, ok)
	require.Equal(t, int64(400), status)

	error, ok := exports["error"].(string)
	require.True(t, ok)
	require.Equal(t, "", error)

	headers, ok := exports["headers"].(map[string]any)
	require.True(t, ok)
	require.NotEmpty(t, headers)
}

func TestJSLibHTTPPostJSON(t *testing.T) {
	script := `
    import http from "flyscrape/http"

    const res = http.postJSON("https://example.com", {
        username: "foo",
        password: "bar",
    })

    export const body = res.body;
    export const status = res.status;
    export const error = res.error;
    export const headers = res.headers;
    `

	client := &http.Client{
		Transport: flyscrape.RoundTripFunc(func(r *http.Request) (*http.Response, error) {
			require.Equal(t, "POST", r.Method)
			require.Equal(t, "application/json", r.Header.Get("Content-Type"))

			m := map[string]any{}
			json.NewDecoder(r.Body).Decode(&m)
			require.Equal(t, "foo", m["username"])
			require.Equal(t, "bar", m["password"])

			return flyscrape.MockResponse(400, "Bad Request")
		}),
	}

	exports, err := flyscrape.Compile(script, flyscrape.NewJSLibrary(client))
	require.NoError(t, err)

	body, ok := exports["body"].(string)
	require.True(t, ok)
	require.Equal(t, "Bad Request", body)

	status, ok := exports["status"].(int64)
	require.True(t, ok)
	require.Equal(t, int64(400), status)

	error, ok := exports["error"].(string)
	require.True(t, ok)
	require.Equal(t, "", error)

	headers, ok := exports["headers"].(map[string]any)
	require.True(t, ok)
	require.NotEmpty(t, headers)
}
