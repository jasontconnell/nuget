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