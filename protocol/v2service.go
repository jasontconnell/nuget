package protocol

import (
	"net/http"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"strings"
)

type v2Service struct {
	resourceUrl string
}

func NewV2Client(url string) Client {
	return v2Service{ resourceUrl: url }
}

func (svc v2Service) IsValid() bool {
	return !strings.HasSuffix(svc.resourceUrl, ".json")
}

func (svc v2Service) GetServiceVersion() int {
	return 2
}

func (svc v2Service) GetPackageData(id string) (Package, error) {
	var result Package
	url := fmt.Sprintf(svc.getSearchUrlFormat(), id)
	var feed v2feed

	err := xmlRequest(url, &feed)

	if err != nil {
		return result, err
	}

	for _, entry := range feed.Entries {
		if result.Id == "" {
			result.Id = entry.Properties.Id
		}
		result.Versions = append(result.Versions, Version { Url: entry.Id, Version: entry.Properties.Version, DownloadUrl: entry.Content.DownloadUrl })
	}
	return result, nil
}

func (svc v2Service) VersionExists(id, version string) bool {
	return VersionExists(svc, id, version)
}

func (svc v2Service) DownloadPackage(version Version) (io.Reader, error) {
	
	return DownloadData(version.DownloadUrl)
	
}

// private helpers
func (svc v2Service) getSearchUrlFormat() string {
	return svc.resourceUrl + "FindPackagesById()?id='%s'"
}

func (svc v2Service) getDownloadUrl(id, version string) string {
	return ""
}




func getResourceUrl(url string) (string, error) {
	var service v2serviceResponse
	err := xmlRequest(url, &service)
	if err != nil {
		return "", err
	}
	if service.Workspace.Collection.Href == "" {
		return "", errors.New("invalid response, probably not a v2 service")
	}

	return url + service.Workspace.Collection.Href, nil

}

func xmlRequest(url string, out interface{}) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	dec := xml.NewDecoder(resp.Body)
	err = dec.Decode(&out)
	if err != nil {
		return err
	}

	return nil
}