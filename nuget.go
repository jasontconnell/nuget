package nuget

import (
	"github.com/jasontconnell/nuget/protocol"
	"github.com/pkg/errors"
	"io/ioutil"
	"os"
	"path"
)

type Service struct {
	client client
}

func getClient(url string) client {
	v2client := protocol.NewV2Client(url)
	if v2client.IsValid() {
		return v2client
	}

	v3client := protocol.NewV3Client(url)
	if v3client.IsValid() {
		return v3client
	}

	return nil
}

func NewService(url string) Service {
	client := getClient(url)
	return Service{client: client}
}

func FindPackage(urls []string, packageId string) (protocol.Package, error) {
	var err error
	var pkg protocol.Package
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
	v, err := svc.client.GetVersion(id, version)

	if err != nil {
		return "", errors.Wrapf(err, "getting version %s %s", id, version)
	}

	if v.Version == "" {
		return "", errors.Wrapf(err, "version not found %s %s", id, version)
	}

	r, err := svc.client.DownloadPackage(v)

	if err != nil {
		return "", errors.Wrapf(err, "Error downloading version, %s %s.", id, version)
	}

	fn := path.Join(folder, id+version+".nupkg")

	f, err := os.OpenFile(fn, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return "", errors.Wrapf(err, "opening destination file failed, %s", fn)
	}

	b, err := ioutil.ReadAll(r)
	if err != nil {
		return "", errors.Wrapf(err, "reading downloaded file failed, %s %s %s", fn, id, version)
	}

	_, err = f.Write(b)
	if err != nil {
		return "", errors.Wrapf(err, "writing destination file, %s %s %s", fn, id, version)
	}

	return fn, nil
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
