package binary

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/google/go-github/github"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/ylallemant/githook-companion/pkg/git"
)

func ListReleases() ([]*github.RepositoryRelease, error) {
	var releases []*github.RepositoryRelease

	uri, err := git.Repository()
	if err != nil {
		return releases, err
	}

	owner, repo, err := git.OwnerAndRepositoryFromUri(uri)
	if err != nil {
		return releases, err
	}

	log.Debug().Msgf("list releases for repo %s/%s\n", owner, repo)

	client := github.NewClient(nil)

	releases, _, err = client.Repositories.ListReleases(context.Background(), owner, repo, nil)

	return releases, err
}

func Latest(releases []*github.RepositoryRelease, allowPrerelease bool) *github.RepositoryRelease {
	for _, release := range releases {
		if release.GetPrerelease() && allowPrerelease {
			return release
		}

		if !release.GetPrerelease() {
			return release
		}
	}

	return releases[0]
}

func VersionsInSync() (bool, error) {
	releases, err := ListReleases()
	if err != nil {
		return false, errors.Wrapf(err, "failed to list binery releases")
	}

	latest := Latest(releases, false)
	localVersion := Semver()

	if localVersion == defaultVersion {
		// go run is being used
		return false, nil
	}

	return latest.GetTagName() == localVersion, nil
}

func Upgrade(currentLocation, tempDir string, wanted *github.RepositoryRelease) error {
	localVersion := Semver()

	uri, err := git.Repository()
	if err != nil {
		return err
	}

	binaryName, err := git.RepositoryFromUri(uri)
	if err != nil {
		return err
	}

	binaryAsset, found := matchingBinary(wanted)
	if !found {
		return errors.Errorf("no matching binary found for %s/%s in release %s", runtime.GOARCH, runtime.GOOS, wanted.GetTagName())
	} else {
		log.Debug().Msgf("matching binary \"%s\" found for %s/%s at %s\n", wanted.GetTagName(), runtime.GOARCH, runtime.GOOS, binaryAsset.GetBrowserDownloadURL())
	}

	binary, err := downloadArchive(binaryAsset.GetBrowserDownloadURL())
	if err != nil {
		return errors.Wrapf(err, "failed to download compressed binary")
	}
	defer binary.Close()

	err = saveFile(fullPath(tempDir, binaryAsset.GetName()), binary)
	if err != nil {
		return errors.Wrapf(err, "failed to write compressed binary locally")
	}
	log.Debug().Msgf("downloaded compressed binary at %s", fullPath(tempDir, binaryAsset.GetName()))

	binaryMD5sum, err := calculateMD5(fullPath(tempDir, binaryAsset.GetName()))
	if err != nil {
		return errors.Wrapf(err, "failed to calculate binary checksum locally")
	}
	log.Debug().Msgf("binary checksum %s", binaryMD5sum)

	binaryChecksum, checksumFound := matchingChecksum(wanted)
	if !checksumFound {
		log.Debug().Msgf("no matching checksum found for %s/%s in release %s\n", runtime.GOARCH, runtime.GOOS, wanted.GetTagName())
	} else {
		checksum, err := downloadArchive(binaryChecksum.GetBrowserDownloadURL())
		if err != nil {
			return errors.Wrapf(err, "failed to download checksum")
		}
		defer checksum.Close()

		err = saveFile(fullPath(tempDir, binaryChecksum.GetName()), checksum)
		if err != nil {
			return errors.Wrapf(err, "failed to write checksum locally")
		}

		log.Debug().Msgf("downloaded checksum at %s", fullPath(tempDir, binaryChecksum.GetName()))
		checksumValue, err := readAsString(fullPath(tempDir, binaryChecksum.GetName()))
		if err != nil {
			return errors.Wrapf(err, "failed to get md5 checksum value")
		}

		if checksumValue != binaryMD5sum {
			return errors.Errorf("downloaded compressed binary is corrupted, md5sum mismatch %s != %s (awaited)", binaryMD5sum, checksumValue)
		}

		log.Debug().Msgf("downloaded compressed binary is pristine, md5sum match (%s)\n", checksumValue)
	}

	err = decompress(tempDir, binaryAsset.GetName())
	if err != nil {
		return errors.Wrap(err, "failed to decompress binary")
	}

	log.Debug().Msgf("decompressed binary at %s", fullPath(tempDir, binaryName))

	err = MoveFile(fullPath(tempDir, binaryName), currentLocation)
	if err != nil {
		return errors.Wrapf(err, "moving binary failed")
	}

	log.Debug().Msgf("binary updated from \"%s\" to \"%s\"\n", localVersion, wanted.GetTagName())

	return nil
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

func downloadArchive(uri string) (io.ReadCloser, error) {
	downloader := http.Client{}

	request, err := downloader.Get(uri)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to download compressed binary at %s", uri)
	}

	if request.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed with status code %s", strconv.Itoa(request.StatusCode))
	}

	return request.Body, nil
}

func matchingBinary(release *github.RepositoryRelease) (*github.ReleaseAsset, bool) {
	for _, asset := range release.Assets {
		if checkForMatchingAsset(asset.GetName(), false) {
			return &asset, true
		}
	}

	return nil, false
}

func matchingChecksum(release *github.RepositoryRelease) (*github.ReleaseAsset, bool) {
	for _, asset := range release.Assets {
		if checkForMatchingAsset(asset.GetName(), true) {
			return &asset, true
		}
	}

	return nil, false
}

func checkForMatchingAsset(name string, wantMD5 bool) bool {
	return strings.Contains(name, runtime.GOARCH) &&
		strings.Contains(name, runtime.GOOS) &&
		isMD5Asset(name) == wantMD5
}

func isMD5Asset(name string) bool {
	return strings.Contains(name, ".md5")
}

func decompress(tempDir, fileName string) error {
	// Open the tar.gz file
	file, err := os.Open(fullPath(tempDir, fileName))
	if err != nil {
		return err
	}
	defer file.Close()

	// Create a gzip reader
	gr, err := gzip.NewReader(file)
	if err != nil {
		return err
	}
	defer gr.Close()

	// Create a tar reader
	tr := tar.NewReader(gr)

	// Extract files
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break // End of archive
		}
		if err != nil {
			return err
		}

		// Construct the output path
		outputPath := filepath.Join(tempDir, header.Name)

		// Create directories if needed
		if header.Typeflag == tar.TypeDir {
			if err := os.MkdirAll(outputPath, 0755); err != nil {
				return err
			}
			continue
		}

		// Create the output file
		outFile, err := os.Create(outputPath)
		if err != nil {
			return err
		}
		defer outFile.Close()

		// Copy the file content
		if _, err := io.Copy(outFile, tr); err != nil {
			return err
		}

		fmt.Printf("decompressed archive at %s\n", outputPath)
	}

	return nil
}

func fullPath(path, filename string) string {
	return filepath.Join(path, filename)
}

func Location() (string, error) {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	output := filepath.Dir(ex)
	filename := filepath.Base(ex)

	return fullPath(output, filename), nil
}

func MoveFile(sourcePath, targetPath string) error {
	content, err := os.ReadFile(sourcePath)
	if err != nil {
		return errors.Wrapf(err, "failed to read source file %s", sourcePath)
	}

	err = os.Remove(targetPath)
	if err != nil {
		return errors.Wrapf(err, "failed to delete old target file version %s", targetPath)
	}

	err = os.WriteFile(targetPath, content, 0777)
	if err != nil {
		return errors.Wrapf(err, "failed to write to target file %s", targetPath)
	}

	err = os.Chmod(targetPath, 0777)
	if err != nil {
		return errors.Wrapf(err, "failed to change target file's mode to 0777 %s", targetPath)
	}

	// The copy was successful, so now delete the original file
	err = os.Remove(sourcePath)
	if err != nil {
		return errors.Wrapf(err, "failed to delete source file %s", sourcePath)
	}
	return nil
}

func calculateMD5(filePath string) (string, error) {
	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Initialize the MD5 hash
	hash := md5.New()

	// Read the file in chunks and update the hash
	buffer := make([]byte, 4096)
	for {
		bytesRead, err := file.Read(buffer)
		if err != nil {
			break
		}
		hash.Write(buffer[:bytesRead])
	}

	// Get the hexadecimal representation of the hash
	md5Hash := hex.EncodeToString(hash.Sum(nil))
	return md5Hash, nil
}

func readAsString(filepath string) (string, error) {
	content, err := os.ReadFile(filepath)
	if err != nil {
		return "", errors.Wrapf(err, "failed to read file %s", filepath)
	}

	return strings.TrimSpace(string(content)), nil
}
