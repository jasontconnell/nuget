package protocol

type resourcesResponse struct {
	Resources []resource `json:"resources"`
}

type resource struct {
	Id      string `json:"@id"`
	Type    string `json:"@type"`
	Comment string `json:"comment"`
}

type queryResult struct {
	TotalHits int         `json:"totalHits"`
	Data      []v3Package `json:"data"`
}

type v3Package struct {
	Id            string      `json:"id"`
	LatestVersion string      `json:"version"`
	Versions      []v3Version `json:"versions"`
}

type v3Version struct {
	Version         string `json:"version"`
	RegistrationUrl string `json:"@id"`
}

type v3Registration struct {
	Id          string `json:"@id"`
	DownloadUrl string `json:"packageContent"`
}
