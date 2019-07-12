package protocol

import (
	"io"
)

type Package struct {
	Id       string    `json:"id"`
	Versions []Version `json:"versions"`
}

type Version struct {
	Version     string `json:"version"`
	Url         string `json:"@id"`
	DownloadUrl string `json:"downloadUrl"`
}

type Client interface {
	GetServiceVersion() int
	GetPackageData(id string) (Package, error)
	VersionExists(id, version string) bool
	DownloadPackage(version Version) (io.Reader, error)
	IsValid() bool
}

