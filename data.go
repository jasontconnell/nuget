package nuget

type Resource struct {
	Id      string `json:"@id"`
	Type    string `json:"@type"`
	Comment string `json:"comment"`
}

type ResourceList struct {
	Version   string     `json:"version"`
	Resources []Resource `json:"resources"`
}

type QueryResult struct {
	TotalHits int          `json:"totalHits"`
	Data      []ResultData `json:"data"`
}

type ResultData struct {
	Id       string    `json:"id"`
	Versions []Version `json:"versions"`
}

type Version struct {
	Version string `json:"version"`
	Url     string `json:"@id"`
}

type VersionData struct {
	PackageUrl string `json:"packageContent"`
}
