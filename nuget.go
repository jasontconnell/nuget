package nuget

import (
	"fmt"

	"github.com/jasontconnell/nuget/protocol"
)

type Service struct {
	client protocol.Client
}

func NewService(url string) Service {
	client := protocol.GetClient(url)
	return Service{client: client}
}

func FindPackage(urls []string, packageId string) (Package, error) {
	var err error
	var pkg Package
	for _, url := range urls {
		ns := NewService(url)
		if ns.client.IsValid() {
			pkg, err = ns.client.GetPackageData(packageId)
			if err != nil {
				continue
			}
			err = nil
		}
		break
	}
	return pkg, err
}

func (svc Service) Download(id, version, folder string) (string, error) {

	vexists := svc.client.VersionExists(id, version)

	if !vexists {
		return "", fmt.Errorf("package version doesn't exist %s - %s", id, version)
	}

	pkgdata, err := svc.client.GetPackageData(id)

	if err != nil {
		return "", fmt.Errorf("Couldn't get package data, %v", err)
	}



	// v := svc.GetVersion(pkgdata, version)

	// if v.Version == "" {
	// 	return "", fmt.Errorf("Version not found, %v", version)
	// }

	// fn, err := svc.DownloadVersion(v, folder)

	// if err != nil {
	// 	return "", fmt.Errorf("Error downloading version, %v.", err)
	// }

	// return fn, nil
}

// func DownloadAndExtract(svcUrl, id, version, downloadFolder, extractFolder string) (string, error) {
// 	return "", nil
// 	// fn, err := Download(svcUrl, id, version, downloadFolder)

// 	// if err != nil {
// 	// 	return "", err
// 	// }

// 	// zr, zrerr := zip.OpenReader(filepath.Join(downloadFolder, fn))

// 	// if zrerr != nil {
// 	// 	return "", zrerr
// 	// }

// 	// defer zr.Close()

// 	// dir := filepath.Join(extractFolder, id, version)
// 	// merr := os.MkdirAll(dir, os.ModePerm)

// 	// if merr != nil {
// 	// 	return "", merr
// 	// }

// 	// for _, f := range zr.File {
// 	// 	fpath, n := path.Split(f.Name)
// 	// 	extPath := filepath.Join(append([]string{dir}, strings.Split(fpath, "/")...)...)
// 	// 	meerr := os.MkdirAll(extPath, os.ModePerm)
// 	// 	if meerr != nil {
// 	// 		fmt.Println("couldn't create folder", extPath)
// 	// 		continue
// 	// 	}
// 	// 	fout, ferr := os.Create(filepath.Join(extPath, n))
// 	// 	docopy := true

// 	// 	if ferr != nil {
// 	// 		fmt.Println(ferr)
// 	// 		docopy = false
// 	// 	}

// 	// 	rc, rcerr := f.Open()
// 	// 	if rcerr != nil {
// 	// 		fmt.Println(rcerr)
// 	// 		docopy = false
// 	// 	}

// 	// 	if docopy {
// 	// 		_, cerr := io.Copy(fout, rc)
// 	// 		if cerr != nil {
// 	// 			fmt.Println(cerr)
// 	// 		}
// 	// 	}

// 	// 	rc.Close()
// 	// }

// 	// return dir, nil
// }

// func GetLatestVersion(svcUrl, id string) (string, error) {
// 	return "", nil
// 	// var err error
// 	// svc := NewService(svcUrl)
// 	// err = svc.GetResources()

// 	// if err != nil {
// 	// 	return "", fmt.Errorf("Couldn't get resources, %v", err)
// 	// }

// 	// pkgdata, err := svc.GetPackageData(id)
// 	// version := getHighestVersion(pkgdata.Versions)

// 	// return version, nil
// }

// func (svc *NugetService) GetPackageData(id string) (ResultData, error) {
// 	return ResultData{}, nil
// 	// pkg := ResultData{}

// 	// q, err := svc.getQueryService()

// 	// if err != nil {
// 	// 	return pkg, err
// 	// }

// 	// url := fmt.Sprintf(`%s?q=@Id:"%s"&prerelease=false`, q.Id, id)
// 	// resp, rerr := http.Get(url)
// 	// if rerr != nil {
// 	// 	return pkg, fmt.Errorf("Failed in call to query, %v", rerr)
// 	// }

// 	// defer resp.Body.Close()

// 	// qr := QueryResult{}
// 	// decodeJson(resp.Body, &qr)

// 	// if qr.TotalHits == 0 {
// 	// 	return pkg, fmt.Errorf("No results, %v", url)
// 	// }

// 	// for _, res := range qr.Data {
// 	// 	if res.Id == id {
// 	// 		pkg = res
// 	// 		break
// 	// 	}
// 	// }

// 	// return pkg, nil
// }

// func (svc *NugetService) GetVersion(res ResultData, version string) Version {
// 	return Version{}
// 	// v := Version{}
// 	// for _, ver := range res.Versions {
// 	// 	if ver.Version == version {
// 	// 		v = ver
// 	// 		break
// 	// 	}
// 	// }

// 	// return v
// }

// func (svc *NugetService) DownloadVersion(version Version, folder string) (string, error) {
// 	return "", nil
// 	// if filepath.IsAbs(folder) {
// 	// 	err := os.MkdirAll(folder, os.ModePerm)
// 	// 	if err != nil {
// 	// 		return "", fmt.Errorf("Couldn't create directories, %v - %v.", folder, err)
// 	// 	}
// 	// }

// 	// vd, vderr := svc.getVersionData(version)
// 	// if vderr != nil {
// 	// 	return "", vderr
// 	// }
// 	// _, fn := path.Split(vd.PackageUrl)
// 	// f, err := os.Create(filepath.Join(folder, fn))

// 	// if err != nil {
// 	// 	return "", err
// 	// }

// 	// resp, gerr := http.Get(vd.PackageUrl)
// 	// if gerr != nil {
// 	// 	return "", gerr
// 	// }

// 	// defer resp.Body.Close()
// 	// defer f.Close()

// 	// _, cerr := io.Copy(f, resp.Body)

// 	// return fn, cerr
// }

// func (svc *NugetService) getVersionData(version Version) (VersionData, error) {

// 	return VersionData{}, nil
// 	// vd := VersionData{}
// 	// resp, err := http.Get(version.Url)

// 	// if err != nil {
// 	// 	return vd, err
// 	// }

// 	// defer resp.Body.Close()

// 	// decodeJson(resp.Body, &vd)

// 	// return vd, nil
// }
