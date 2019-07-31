package protocol

import (
	"encoding/xml"
)

type v2serviceResponse struct {
	XMLName   xml.Name            `xml:"service"`
	Workspace v2workspaceResponse `xml:"workspace"`
	Xmlns     string              `xml:"xmlns,attr"`
}

type v2workspaceResponse struct {
	XMLName    xml.Name             `xml:"workspace"`
	Collection v2collectionResponse `xml:"collection"`
}

type v2collectionResponse struct {
	XMLName xml.Name `xml:"collection"`
	Href    string   `xml:"href,attr"`
}

type v2feed struct {
	XMLName xml.Name  `xml:"feed"`
	Entries []v2entry `xml:"entry"`
}

type v2entry struct {
	XMLName    xml.Name     `xml:"entry"`
	Id         string       `xml:"id"`
	Content    v2content    `xml:"content"`
	Properties v2properties `xml:"http://schemas.microsoft.com/ado/2007/08/dataservices/metadata properties"`
}

type v2content struct {
	DownloadUrl string `xml:"src,attr"`
}

type v2properties struct {
	XMLName xml.Name `xml:"properties"`
	// Id string `xml:"properties>Id"`
	Id              string `xml:"http://schemas.microsoft.com/ado/2007/08/dataservices Id"`
	Version         string `xml:"http://schemas.microsoft.com/ado/2007/08/dataservices Version"`
	IsLatestVersion bool   `xml:"http://schemas.microsoft.com/ado/2007/08/dataservices IsLatestVersion"`
}
