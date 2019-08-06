package sh

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/craiggwilson/goke/task"
)

// DownloadHTTP issues a GET request against the provided url and downloads the contents to the toPath.
func DownloadHTTP(ctx *task.Context, url string, toPath string) error {
	ctx.Logf("download: %s -> %s\n", url, toPath)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed creating GET request: %v", err)
	}
	req.Header.Add("cache-control", "no-cache")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed issuing GET request: %v", err)
	}
	if res.Body == nil {
		return errors.New("no body from the GET request")
	}
	defer res.Body.Close()

	return copyTo(url, res.Body, toPath, 0666)
}
