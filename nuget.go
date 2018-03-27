package nuget

import (
    "fmt"
    "net/http"
    "encoding/json"
    "io"
    "path"
    "path/filepath"
    "os"
)

func decodeJson(r io.Reader, obj interface{}){
    dec := json.NewDecoder(r)
    dec.Decode(obj)
}

type NugetService struct {
    DefinitionUrl string
    resources ResourceList
}


func NewService(url string) *NugetService {
    n := new(NugetService)
    n.DefinitionUrl = url
    return n
}

func Download(svcUrl, id, version, folder string) (string, error) {
    var err error
    svc := NewService(svcUrl)
    err = svc.GetResources()

    if err != nil {
        return "", fmt.Errorf("Couldn't get resources, %v", err)
    }

    pkgdata,err := svc.GetPackageData(id)

    if err != nil {
        return "", fmt.Errorf("Couldn't get package data, %v", err)
    }

    v := svc.GetVersion(pkgdata, version)

    if v.Version == "" {
        return "", fmt.Errorf("Version not found, %v", version)
    }

    fn, err := svc.DownloadVersion(v, folder)

    if err != nil {
        return "", fmt.Errorf("Error downloading version, %v.", err)
    }

    return fn, nil
}

func (svc *NugetService) GetResources() error {
    if len(svc.resources.Resources) > 0 {
        return nil
    }

    resp, err := http.Get(svc.DefinitionUrl)
    if err != nil {
        return fmt.Errorf("Couldn't get resources, %v", err)
    }

    defer resp.Body.Close()

    decodeJson(resp.Body, &svc.resources)

    if len(svc.resources.Resources) == 0 {
        return fmt.Errorf("Couldn't get resources, no resources in body")
    }

    return nil
}

func (svc *NugetService) getQueryService() (Resource, error) {
    r := Resource{}
    err := svc.GetResources()
    if err != nil {
        return r, fmt.Errorf("Couldn't get resources, %v", err)
    }

    // just get the first query service
    for _, res := range svc.resources.Resources {
        if res.Type == "SearchQueryService" {
            r = res
            break
        }
    }

    if r.Id == "" {
        return r, fmt.Errorf("Couldn't find query service, %v", svc.DefinitionUrl)
    }

    return r, nil
}

func (svc *NugetService) GetPackageData(id string) (ResultData, error) {
    pkg := ResultData{}

    q, err := svc.getQueryService()

    if err != nil {
        return pkg, err
    }

    url := fmt.Sprintf(`%s?q=@Id:"%s"`, q.Id, id)
    resp, rerr := http.Get(url)
    if rerr != nil {
        return pkg, fmt.Errorf("Failed in call to query, %v", rerr)
    }

    defer resp.Body.Close()

    qr := QueryResult{}
    decodeJson(resp.Body, &qr)

    if qr.TotalHits == 0 {
        return pkg, fmt.Errorf("No results, %v", url)
    }

    for _, res := range qr.Data {
        if res.Id == id {
            pkg = res
            break
        }
    }

    return pkg, nil
}

func (svc *NugetService) GetVersion(res ResultData, version string) Version {
    v := Version{}
    for _, ver := range res.Versions {
        if ver.Version == version {
            v = ver
            break
        }
    }

    return v
}

func (svc *NugetService) DownloadVersion(version Version, folder string) (string, error) {

    if filepath.IsAbs(folder){
        err := os.MkdirAll(folder, os.ModePerm)
        if err != nil {
            return "", fmt.Errorf("Couldn't create directories, %v - %v.", folder, err)
        }
    }

    vd,vderr := svc.getVersionData(version)
    if vderr != nil {
        return "", vderr
    }
    _, fn := path.Split(vd.PackageUrl)
    f,err := os.Create(filepath.Join(folder, fn))

    resp, gerr := http.Get(vd.PackageUrl)
    if gerr != nil {
        return "", gerr
    }

    defer resp.Body.Close()
    defer f.Close()

    if err != nil {
        return "", err
    }

    _, cerr := io.Copy(f, resp.Body)

    return fn, cerr
}

func (svc *NugetService) getVersionData(version Version) (VersionData, error) {
    vd := VersionData{}
    resp, err := http.Get(version.Url)

    if err != nil {
        return vd, err
    }

    defer resp.Body.Close()

    decodeJson(resp.Body, &vd)

    return vd, nil
}