package dependency

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"hash"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

func verify(archiveFilename, checksumFilename, algorithm, directory string) error {
	var hasher hash.Hash

	switch algorithm {
	case "sha256":
		hasher = sha256.New()
	case "sha1":
		hasher = sha1.New()
	case "md5":
		hasher = md5.New()
	default:
		return errors.Errorf("unknown checksum algorithm \"%s\"", algorithm)
	}

	// calculate archive checksum
	file, err := os.Open(filepath.Join(directory, archiveFilename))
	if err != nil {
		return errors.Wrapf(err, "failed to open archive at %s", filepath.Join(directory, archiveFilename))
	}
	defer file.Close()

	_, err = io.Copy(hasher, file)
	if err != nil {
		return errors.Wrapf(err, "failed to hash archive at %s", filepath.Join(directory, archiveFilename))
	}

	archiveChecksum := hex.EncodeToString(hasher.Sum(nil))

	// read expected checksum
	content, err := os.ReadFile(filepath.Join(directory, checksumFilename))
	if err != nil {
		return errors.Wrapf(err, "failed to read checksum file at %s", filepath.Join(directory, checksumFilename))
	}

	expectedChecksum := strings.TrimSpace(string(content))

	// check validity
	if archiveChecksum != expectedChecksum {
		return errors.Errorf("downloaded archive is corrupted, checksum mismatch %s != %s (expected)", archiveChecksum, expectedChecksum)
	}

	return nil
}
