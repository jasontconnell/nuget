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
	TotalHits int          `json:"totalHits"`
	Data      []Package `json:"data"`
}
