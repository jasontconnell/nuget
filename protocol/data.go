package protocol

import (
	"encoding/xml"
)

type Package struct {
	Id         string
	Versions   []Version
	VersionMap map[string]Version
}

type Version struct {
	Version         string
	DownloadUrl     string
	RegistrationUrl string
}

type Nuspec struct {
	XMLName  xml.Name       `xml:"package"`
	MetaData NuspecMetadata `xml:"metadata"`
}

type NuspecMetadata struct {
	XMLName          xml.Name `xml:"metadata"`
	MinClientVersion string   `xml:"minClientVersion,attr"`

	PackageId string `xml:"id"`

	Dependencies     []NuspecDependency      `xml:"dependencies>dependency"`
	DependencyGroups []NuspecDependencyGroup `xml:"dependencies>group"`
}

type NuspecDependencyGroup struct {
	XMLName         xml.Name           `xml:"group"`
	TargetFramework string             `xml:"targetFramework,attr"`
	Dependencies    []NuspecDependency `xml:"dependency"`
}

type NuspecDependency struct {
	XMLName xml.Name `xml:"dependency"`
	Name    string   `xml:"id,attr"`
	Version string   `xml:"version,attr"`
}
