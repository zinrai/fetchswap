package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

// http.DefaultClient, not a custom client: it already follows redirects and
// honors HTTP_PROXY, HTTPS_PROXY and NO_PROXY.
func fetch(ctx context.Context, url string) (io.ReadCloser, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		resp.Body.Close()
		return nil, fmt.Errorf("%s: unexpected status %s", url, resp.Status)
	}
	return resp.Body, nil
}
