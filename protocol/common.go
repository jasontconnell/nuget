package protocol

import (
	"net/http"
	"io"
	"io/ioutil"
	"bytes"
	"strings"
	"strconv"
)

func DownloadData(url string) (io.Reader, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return bytes.NewBuffer(b), nil
}

func GetHighestVersion(versions []Version) string {
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

func VersionExists(c Client, id, version string) bool {
	pkg, err := c.GetPackageData(id)
	if err != nil {
		return false
	}

	for _, v := range pkg.Versions {
		if v.Version == version {
			return true
		}
	}

	return false
}