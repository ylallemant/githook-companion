package dependency

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

func decompress(archiveFilename, directory string) (string, error) {
	path := filepath.Join(directory, archiveFilename)
	ext := filepath.Ext(archiveFilename)

	var err error
	decompressedFilename := ""

	switch ext {
	case ".zip":
		decompressedFilename, err = Unzip(path, directory)
		if err != nil {
			return "", errors.Wrapf(err, "failed decompress Zip archive at %s", path)
		}
	case ".gz":
		decompressedFilename, err = unTarGz(path, directory)
		if err != nil {
			return "", errors.Wrapf(err, "failed decompress tar.gz archive at %s", path)
		}
	default:
		return "", errors.Errorf("unknown archive extention \"%s\"", ext)
	}

	return decompressedFilename, nil
}

func unTarGz(src, dest string) (string, error) {
	// Open the archive
	file, err := os.Open(src)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Create a gzip reader
	gr, err := gzip.NewReader(file)
	if err != nil {
		return "", err
	}
	defer gr.Close()

	// Create a tar reader
	tr := tar.NewReader(gr)

	decompressedFilename := ""

	// Extract files
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break // End of archive
		}
		if err != nil {
			return "", err
		}
		if strings.Contains(header.Name, "..") {
			return "", errors.Errorf("probable arbitrary file access tentative with: \"%s\"", header.Name)
		}

		decompressedFilename = header.Name

		// Construct the output path
		outputPath := filepath.Join(dest, header.Name)

		// Create directories if needed
		if header.Typeflag == tar.TypeDir {
			if err := os.MkdirAll(outputPath, 0755); err != nil {
				return "", err
			}
			continue
		}

		// Create the output file
		outFile, err := os.Create(outputPath)
		if err != nil {
			return "", err
		}
		defer outFile.Close()

		// Copy the file content
		if _, err := io.Copy(outFile, tr); err != nil {
			return "", err
		}
	}

	return decompressedFilename, nil
}

func Unzip(src, dest string) (string, error) {
	r, err := zip.OpenReader(src)
	if err != nil {
		return "", err
	}
	defer func() {
		if err := r.Close(); err != nil {
			panic(err)
		}
	}()

	os.MkdirAll(dest, 0755)

	decompressedFilename := ""

	// Closure to address file descriptors issue with all the deferred .Close() methods
	extractAndWriteFile := func(f *zip.File) error {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer func() {
			if err := rc.Close(); err != nil {
				panic(err)
			}
		}()

		path := filepath.Join(dest, f.Name)

		// Check for ZipSlip (Directory traversal)
		if !strings.HasPrefix(path, filepath.Clean(dest)+string(os.PathSeparator)) {
			return errors.Errorf("illegal file path: %s", path)
		}

		if f.FileInfo().IsDir() {
			os.MkdirAll(path, f.Mode())
		} else {
			os.MkdirAll(filepath.Dir(path), f.Mode())
			f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer func() {
				if err := f.Close(); err != nil {
					panic(err)
				}
			}()

			_, err = io.Copy(f, rc)
			if err != nil {
				return err
			}

			decompressedFilename = filepath.Base(f.Name())
		}
		return nil
	}

	for _, f := range r.File {
		err := extractAndWriteFile(f)
		if err != nil {
			return "", err
		}
	}

	return decompressedFilename, nil
}
