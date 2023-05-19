package protocol

import (
	"archive/zip"
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
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

func getNuspec(pkgId, version string, r io.Reader) (*Nuspec, error) {
	nuspec := &Nuspec{}

	b := new(bytes.Buffer)
	n, err := b.ReadFrom(r)
	if err != nil {
		return nil, fmt.Errorf("nuspec: reading from download stream %s %s. %w", pkgId, version, err)
	}

	br := bytes.NewReader(b.Bytes())
	zr, err := zip.NewReader(br, n)
	if err != nil {
		return nil, fmt.Errorf("nuspec: creating reader for nupkg zip %s %s, file len %d. %w", pkgId, version, n, err)
	}

	var nuspecReader io.Reader
	var ferr error
	nfilename := pkgId + ".nuspec"
	for _, f := range zr.File {
		if f.Name == nfilename {
			nuspecReader, ferr = f.Open()
			if ferr != nil {
				return nil, fmt.Errorf("nuspec: opening file inside zip %s. %w", nfilename, ferr)
			}
			break
		}
	}

	if nuspecReader == nil {
		return nil, fmt.Errorf("nuspec: no reader, no .nuspec file in package, %s %s", pkgId, version)
	}

	xr := xml.NewDecoder(nuspecReader)
	xerr := xr.Decode(nuspec)
	if xerr != nil {
		return nil, fmt.Errorf("nuspec: decoding xml %s %s. %w", pkgId, version, xerr)
	}

	return nuspec, nil
}
