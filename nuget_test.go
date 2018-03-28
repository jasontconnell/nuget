package nuget

import (
    "testing"
)

var serviceUrl = "https://api.nuget.org/v3/index.json"

func GetService(t *testing.T) *NugetService {
    t.Helper()

    n := NewService(serviceUrl)
    return n
}

func TestCreate(t *testing.T){
    n := GetService(t)

    if n.DefinitionUrl != serviceUrl {
        t.Fail()
    }
}

func TestGetResources(t *testing.T){
    n := GetService(t)

    err := n.GetResources()
    if err != nil {
        t.Errorf("Couldn't get resources, %v", err)
        t.Fail()
    }

    for _, r := range n.resources.Resources {
        t.Log(r.Type)
    }
}

func TestGetPackage(t *testing.T){
    n := GetService(t)

    pkgs := []struct {
        id string
        fail bool
    }{
        { id: "Newtonsoft.Json", fail: true },
        { id: "Sitecore.Kernel.NoReference", fail: false },
    }

    for _, p := range pkgs {
        t.Run("Package " + p.id, func (t *testing.T){
            pkg, err := n.GetPackageData(p.id)
            if err != nil && p.fail {
                t.Log(err)
                t.Fail()
                return
            }

            if len(pkg.Versions) == 0 {
                t.Log("No versions")
            }
        })
    }
}

func TestDownloadPackage(t *testing.T){
    n := GetService(t)

    pkgs := []struct {
        id string
        version string
    }{
        { id: "Newtonsoft.Json", version: "8.0.1" },
        { id: "Glass.Mapper.Sc" , version: "4.2.1.188" },
        { id: "bootstrap" , version: "4.0.0" },
    }

    for _, p := range pkgs {
        t.Run("Package " + p.id + "-" + p.version, func (t *testing.T){
            pkg, err := n.GetPackageData(p.id)
            if err != nil {
                t.Log(err)
                t.Fail()
                return
            }

            v := n.GetVersion(pkg, p.version)

            s,err := n.DownloadVersion(v, ".")

            if err != nil {
                t.Log(err)
                t.Fail()
            }

            t.Log(s)
        })
    }
}

func TestDownloadAndExtractPackage(t *testing.T){

    pkgs := []struct {
        id string
        version string
    }{
        { id: "Newtonsoft.Json", version: "8.0.1" },
        { id: "Glass.Mapper.Sc" , version: "4.2.1.188" },
        { id: "bootstrap" , version: "4.0.0" },
    }

    for _, p := range pkgs {
        t.Run("Package " + p.id + "-" + p.version, func (t *testing.T){
            fn, err := DownloadAndExtract(serviceUrl, p.id, p.version, `c:\test\test`, `c:\test\extract`)
            if err != nil {
                t.Log(err)
                t.Fail()
                return
            }

            t.Log(fn)
        })
    }
}

func TestGetLatestVersion(t *testing.T){
    pkgs := []struct {
        id string
        version string
        result bool // should equal or should not equal
    }{
        { id: "Newtonsoft.Json", version: "8.0.1", result: false },
        { id: "Glass.Mapper.Sc" , version: "4.2.1.188", result: false },
        { id: "bootstrap" , version: "4.0.0", result: true },
    }

        for _, p := range pkgs {
        t.Run("Package " + p.id + "-" + p.version, func (t *testing.T){
            vstr,err := GetLatestVersion(serviceUrl, p.id)
            if err != nil {
                t.Logf("Got error, expecting version, %v", err)
                t.Fail()
            }

            if vstr != p.version && p.result {
                t.Logf("Versions not equal. %v and got latest from service: %v", p.version, vstr)
                t.Fail()
            } else {
                t.Logf("Got latest version from service for %v, latest is %v. The result is expected.", p.id, vstr)
            }
        })
    }
}