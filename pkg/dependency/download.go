package dependency

import (
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

type archiveContext struct {
	Version string
	Os      string
	Arch    string
	Ext     string
}

func filenameFromUrl(uri string) (string, error) {
	uriObj, err := url.Parse(uri)
	if err != nil {
		return "", errors.Wrapf(err, "failed to parse origin uri %s", uri)
	}

	pathParts := strings.Split(uriObj.Path, "/")

	return pathParts[(len(pathParts) - 1)], nil
}

func download(uri, filename, targetDirectory string) error {
	payload, err := downloadFile(uri)
	if err != nil {
		return errors.Wrapf(err, "failed download payload")
	}
	defer payload.Close()

	err = saveFile(filepath.Join(targetDirectory, filename), payload)
	if err != nil {
		return errors.Wrapf(err, "failed to save payload locally")
	}

	return nil
}

func downloadFile(uri string) (io.ReadCloser, error) {
	downloader := http.Client{}

	request, err := downloader.Get(uri)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to download file at %s", uri)
	}

	if request.StatusCode != http.StatusOK {
		return nil, errors.Errorf("request failed with status code %s", strconv.Itoa(request.StatusCode))
	}

	return request.Body, nil
}

func saveFile(path string, content io.ReadCloser) error {
	file, err := os.Create(path)
	if err != nil {
		return errors.Wrapf(err, "failed to create file at %s", path)
	}

	_, err = io.Copy(file, content)
	if err != nil {
		return errors.Wrapf(err, "failed to write content to file at %s", path)
	}

	return nil
}
