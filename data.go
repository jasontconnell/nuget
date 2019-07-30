package nuget

import (
	"github.com/jasontconnell/nuget/protocol"
	"io"
)

type client interface {
	GetServiceVersion() int
	GetPackageData(id string) (protocol.Package, error)
	GetVersion(id, version string) (protocol.Version, error)
	DownloadPackage(version protocol.Version) (io.Reader, error)
	GetNuspec(id, version string) (*protocol.Nuspec, error) 
	IsValid() bool
}
