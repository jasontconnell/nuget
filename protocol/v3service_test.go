package protocol

import (
	"testing"
)

func TestV3Service(t *testing.T) {
	url := "https://api.nuget.org/v3/index.json"

	svc := NewV3Client(url)

	if !svc.IsValid() {
		t.Fatal("not a valid v3 service")
	}

	pkg, err := svc.GetPackageData("Newtonsoft.Json")

	if err != nil {
		t.Fatal(err)
	}

	t.Log(pkg.Id, pkg.Versions)
}