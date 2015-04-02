package bintray

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

func get(path string) []byte {
	client := &http.Client{}

	paths := []string{apiEndpoint, path}

	req, err := http.NewRequest("GET", strings.Join(paths, "/"), nil)
	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	}

	response, err := client.Do(req)
	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	}

	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	}

	return contents
}
