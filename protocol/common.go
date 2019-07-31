package protocol

import (
	"archive/zip"
	"bytes"
	"encoding/xml"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

func makeVersionMap(versions []Version) map[string]Version {
	m := make(map[string]Version)
	for _, v := range versions {
		m[v.Version] = v
	}
	return m
}

func downloadData(url string) (io.Reader, error) {
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

func getHighestVersion(versions []Version) string {
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

func getNuspec(pkgId, version string, r io.Reader) (*Nuspec, error) {
	nuspec := &Nuspec{}

	b := new(bytes.Buffer)
	n, err := b.ReadFrom(r)
	if err != nil {
		return nil, errors.Wrapf(err, "nuspec: reading from download stream %s %s", pkgId, version)
	}

	br := bytes.NewReader(b.Bytes())
	zr, err := zip.NewReader(br, n)
	if err != nil {
		return nil, errors.Wrapf(err, "nuspec: creating reader for nupkg zip %s %s, file len %d", pkgId, version, n)
	}

	var nuspecReader io.Reader
	var ferr error
	nfilename := pkgId + ".nuspec"
	for _, f := range zr.File {
		if f.Name == nfilename {
			nuspecReader, ferr = f.Open()
			if ferr != nil {
				return nil, errors.Wrapf(ferr, "nuspec: opening file inside zip %s", nfilename)
			}
			break
		}
	}

	if nuspecReader == nil {
		return nil, errors.Errorf("nuspec: no reader, no .nuspec file in package, %s %s", pkgId, version)
	}

	xr := xml.NewDecoder(nuspecReader)
	xerr := xr.Decode(nuspec)
	if xerr != nil {
		return nil, errors.Wrapf(xerr, "nuspec: decoding xml %s %s", pkgId, version)
	}

	return nuspec, nil
}
