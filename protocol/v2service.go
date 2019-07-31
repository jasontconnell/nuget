package protocol

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type v2Service struct {
	resourceUrl string
}

func NewV2Client(url string) v2Service {
	return v2Service{resourceUrl: url}
}

func (svc v2Service) IsValid() bool {
	return !strings.HasSuffix(svc.resourceUrl, ".json")
}

func (svc v2Service) GetServiceVersion() int {
	return 2
}

func (svc v2Service) GetPackageData(id string) (Package, error) {
	var pkg Package
	url := fmt.Sprintf(svc.getSearchUrlFormat(), id)
	var feed v2feed

	err := xmlRequest(url, &feed)

	if err != nil {
		return pkg, err
	}

	for _, entry := range feed.Entries {
		if pkg.Id == "" {
			pkg.Id = entry.Properties.Id
		}
		pkg.Versions = append(pkg.Versions, Version{Version: entry.Properties.Version, DownloadUrl: entry.Content.DownloadUrl})
	}
	pkg.VersionMap = makeVersionMap(pkg.Versions)

	return pkg, nil
}

func (svc v2Service) DownloadPackage(version Version) (io.Reader, error) {
	return downloadData(version.DownloadUrl)
}

func (svc v2Service) GetVersion(id, version string) (Version, error) {
	return Version{}, nil
}

func (svc v2Service) GetNuspec(pkg Package, version Version) (*Nuspec, error) {
	r, err := svc.DownloadPackage(version)
	if err != nil {
		return nil, err
	}

	return getNuspec(pkg.Id, version.Version, r)
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
