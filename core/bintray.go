package core

import (
	"fmt"

	"github.com/tobscher/kiss/bintray"
)

var (
	subject = "tobscher"
	repo    = "generic"
	prefix  = "kiss"
)

func bintrayDownloadUrl(pkg string, os string, arch string) string {
	p := fmt.Sprintf("%v-%v", prefix, pkg)
	bintrayPackage, err := bintray.GetPackage(subject, repo, p)
	if err != nil {
		return ""
	}

	version := bintrayPackage.LatestVersion
	ext := ".tar.gz"

	path := fmt.Sprintf("https://dl.bintray.com/%v/%v/%v_%v_%v_%v%v", subject, repo, pkg, version, os, arch, ext)

	return path
}
