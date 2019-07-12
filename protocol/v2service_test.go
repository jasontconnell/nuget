package protocol

import (
	"testing"
	"fmt"
)

func TestNewV2Service(t *testing.T){
	url := "https://nuget.episerver.com/feed/packages.svc/"

	svc := NewV2Client(url)

	fmt.Println(svc.GetPackageData("EPiServer.CMS.Core"))

}