package bintray

import (
	"encoding/json"
	"fmt"
)

type Package struct {
	Name          string
	Repo          string
	Owner         string
	LatestVersion string `json:"latest_version"`
}

func GetPackage(subject string, repo string, pkg string) (*Package, error) {
	path := fmt.Sprintf("/packages/%v/%v/%v", subject, repo, pkg)
	bytes := get(path)

	var result Package
	err := json.Unmarshal(bytes, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
