package demo

import (
	"fmt"
	"io"
	"net/http"
)

func HTTPDemo() error {
	const (
		targetUrl string = "https://hyperzsb.io"
	)

	response, err := http.Get(targetUrl)
	if err != nil {
		return err
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("get target %s failed", targetUrl)
	}

	fmt.Println(string(body))

	return nil
}
