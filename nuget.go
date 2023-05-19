package nuget

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/jasontconnell/nuget/protocol"
)

type Service struct {
	client
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

	if pkg.Id != packageId {
		return pkg, fmt.Errorf("find package: couldn't find package in any nuget service, %s  %v", packageId, urls)
	}

	return pkg, err
}

func (svc Service) Download(id, version, folder string) (string, error) {
	v, err := svc.client.GetVersion(id, version)

	if err != nil {
		return "", fmt.Errorf("getting version %s %s. %w", id, version, err)
	}

	if v.Version == "" {
		return "", fmt.Errorf("version not found %s %s", id, version)
	}

	r, err := svc.client.DownloadPackage(v)

	if err != nil {
		return "", fmt.Errorf("error downloading version, %s %s. %w", id, version, err)
	}

	fn := path.Join(folder, id+version+".nupkg")

	f, err := os.OpenFile(fn, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return "", fmt.Errorf("opening destination file failed, %s. %w", fn, err)
	}

	b, err := ioutil.ReadAll(r)
	if err != nil {
		return "", fmt.Errorf("reading downloaded file failed, %s %s %s. %w", fn, id, version, err)
	}

	_, err = f.Write(b)
	if err != nil {
		return "", fmt.Errorf("writing destination file, %s %s %s. %w", fn, id, version, err)
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

func GetLatestVersion(urls []string, pkgId string) (string, error) {
	pkg, err := FindPackage(urls, pkgId)
	if err != nil {
		return "", fmt.Errorf("unable to find package %v - %s. %w", urls, pkgId, err)
	}

	hv := getHighestVersion(pkg.Versions)
	return hv, nil
}

func getHighestVersion(versions []protocol.Version) string {
	vs := ""
	high := int64(0)
	for _, v := range versions {
		pts := strings.Split(v.Version, ".")
		if len(pts) > 3 {
			pts = pts[:3]
		}
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
