package protocol

import (
	"testing"
)

func TestNewV2Service(t *testing.T) {
	url := "https://nuget.episerver.com/feed/packages.svc/"

	svc := NewV2Client(url)

	pkg, err := svc.GetPackageData("EPiServer.CMS.Core")

	if err != nil {
		t.Error(err)
	}

	v := pkg.Versions[0]

	nuspec, err := svc.GetNuspec(pkg, v)

	if err != nil {
		t.Error(err)
	}

	t.Log(nuspec)

}
