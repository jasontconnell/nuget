package protocol

import (
	"testing"
)

func TestV3Service(t *testing.T) {
	url := "https://api.nuget.org/v3/index.json"

	svc := NewV3Client(url)

	pkg, err := svc.GetPackageData("Newtonsoft.Json")

	if err != nil {
		t.Error(err)
	}

	if len(pkg.Versions) == 0 {
		t.Error("NO versions")
		t.Fail()
	}

	v := pkg.VersionMap["12.0.1"]
	nuspec, err := svc.GetNuspec(pkg, v)

	if err != nil {
		t.Error(err)
	}

	t.Log(nuspec)
}
