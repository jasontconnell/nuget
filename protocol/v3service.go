package protocol

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type v3Service struct {
	resourceUrl string
	searchUrl   string
}

func NewV3Client(url string) v3Service {
	var searchUrl string
	r, err := getQueryService(url)
	if err == nil {
		searchUrl = r.Id
	} else {
		fmt.Println("error getting query service", err)
	}
	return v3Service{resourceUrl: url, searchUrl: searchUrl}
}

func (svc v3Service) GetServiceVersion() int {
	return 3
}

func (svc v3Service) IsValid() bool {
	return svc.searchUrl != ""
}

func (svc v3Service) GetPackageData(id string) (Package, error) {
	pkg := Package{}
	url := fmt.Sprintf(`%s?q=@Id:"%s"&prerelease=false`, svc.searchUrl, id)
	resp, err := http.Get(url)
	if err != nil {
		return pkg, fmt.Errorf("Failed in call to query, %v", err)
	}

	defer resp.Body.Close()

	qr := queryResult{}

	decodeJson(resp.Body, &qr)

	if qr.TotalHits == 0 {
		return pkg, fmt.Errorf("No results, %v", url)
	}

	for _, res := range qr.Data {
		if res.Id == id {
			pkg = Package{Id: res.Id, Versions: []Version{}}
			for _, v := range res.Versions {
				pkg.Versions = append(pkg.Versions, Version{Version: v.Version, RegistrationUrl: v.RegistrationUrl})
			}
			pkg.VersionMap = makeVersionMap(pkg.Versions)
			break
		}
	}

	return pkg, nil
}

func (svc v3Service) DownloadPackage(version Version) (io.Reader, error) {
	downloadUrl := version.DownloadUrl
	if downloadUrl == "" {
		reg, err := getRegistration(version.RegistrationUrl)
		if err != nil {
			return nil, err
		}
		downloadUrl = reg.DownloadUrl
	}
	return downloadData(downloadUrl)
}

func (svc v3Service) GetVersion(id, version string) (Version, error) {
	return Version{}, nil
}

func (svc v3Service) GetNuspec(pkg Package, version Version) (*Nuspec, error) {
	r, err := svc.DownloadPackage(version)
	if err != nil {
		return nil, err
	}

	return getNuspec(pkg.Id, version.Version, r)
}

// private helpers
func decodeJson(r io.Reader, obj interface{}) {
	dec := json.NewDecoder(r)
	dec.Decode(obj)
}

func getRegistration(url string) (v3Registration, error) {
	reg := v3Registration{}
	resp, err := http.Get(url)
	if err != nil {
		return reg, fmt.Errorf("Couldn't get registration, %v", err)
	}

	defer resp.Body.Close()

	decodeJson(resp.Body, &reg)

	return reg, nil
}

func getResources(url string) ([]resource, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("Couldn't get resources, %v", err)
	}

	defer resp.Body.Close()

	var resresp resourcesResponse
	decodeJson(resp.Body, &resresp)

	if len(resresp.Resources) == 0 {
		return nil, fmt.Errorf("Couldn't get resources, no resources in body")
	}

	return resresp.Resources, nil
}

func getQueryService(url string) (resource, error) {
	r := resource{}
	resources, err := getResources(url)
	if err != nil {
		return r, fmt.Errorf("Couldn't get resources, %v", err)
	}

	// just get the first query service
	for _, res := range resources {
		if res.Type == "SearchQueryService" {
			r = res
			break
		}
	}

	if r.Id == "" {
		return r, fmt.Errorf("Couldn't find query service, %v", url)
	}

	return r, nil
}
