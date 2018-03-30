package nuget

import (
    "fmt"
    "net/http"
    "encoding/json"
    "io"
    "path"
    "path/filepath"
    "os"
    "archive/zip"
    "strings"
    "strconv"
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

func DownloadAndExtract(svcUrl, id, version, downloadFolder, extractFolder string) (string, error) {
    fn, err := Download(svcUrl, id, version, downloadFolder)

    if err != nil {
        return "", err
    }

    zr,zrerr := zip.OpenReader(filepath.Join(downloadFolder, fn))

    if zrerr != nil {
        return "", zrerr
    }

    defer zr.Close()

    dir := filepath.Join(extractFolder, id, version)
    merr := os.MkdirAll(dir, os.ModePerm)

    if merr != nil {
        return "", merr
    }

    for _, f := range zr.File {
        fpath,n := path.Split(f.Name)
        extPath := filepath.Join(append([]string{dir}, strings.Split(fpath, "/")...)...)
        meerr := os.MkdirAll(extPath, os.ModePerm)
        if meerr != nil {
            fmt.Println("couldn't create folder", extPath)
            continue
        }
        fout,ferr := os.Create(filepath.Join(extPath, n))
        docopy := true

        if ferr != nil {
            fmt.Println(ferr)
            docopy = false
        }

        rc, rcerr := f.Open()
        if rcerr != nil {
            fmt.Println(rcerr)
            docopy = false
        }

        if docopy {
            _, cerr := io.Copy(fout, rc)
            if cerr != nil {
                fmt.Println(cerr)
            }
        }

        rc.Close()
    }

    return dir, nil
}

func GetLatestVersion(svcUrl, id string) (string, error) {
    var err error
    svc := NewService(svcUrl)
    err = svc.GetResources()

    if err != nil {
        return "", fmt.Errorf("Couldn't get resources, %v", err)
    }

    pkgdata,err := svc.GetPackageData(id)
    version := getHighestVersion(pkgdata.Versions)

    return version, nil
}

func getHighestVersion(versions []Version) string {
    vs := ""
    high := int64(0)
    for _, v := range versions {
        pts := strings.Split(v.Version, ".")
        m := int64(1)
        cv := int64(0)
        for i := len(pts) - 1; i >= 0; i-- {
            pt := pts[i]
            if len(pt) > 2 {
                pt = string(pt[:2])
            }
            parsed, _ := strconv.ParseInt(pt, 10, 64)
            cv = cv + (parsed * m)
            m = m * int64(100)
        }

        if cv > high {
            high = cv
            vs = v.Version
        }
    }
    return vs
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

    url := fmt.Sprintf(`%s?q=@Id:"%s"&prerelease=false`, q.Id, id)
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

    if err != nil {
        return "", err
    }

    resp, gerr := http.Get(vd.PackageUrl)
    if gerr != nil {
        return "", gerr
    }

    defer resp.Body.Close()
    defer f.Close()

    
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