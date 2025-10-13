package dependency

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/ylallemant/githook-companion/pkg/api"
)

func targetUriTemplate(dependency *api.Dependency) (string, error) {
	if dependency.ArchiveUrlTemplate != "" && dependency.BinaryUrlTemplate != "" {
		return "", errors.Errorf("only an archive OR a binary target url template are allowed, found both for dependency \"%s\"", dependency.Name)
	}

	if dependency.ArchiveUrlTemplate != "" {
		return dependency.ArchiveUrlTemplate, nil
	}

	if dependency.BinaryUrlTemplate != "" {
		return dependency.BinaryUrlTemplate, nil
	}

	return "", errors.Errorf("no target url template found for dependency \"%s\"", dependency.Name)
}

func isArchive(dependency *api.Dependency) bool {
	return dependency.ArchiveUrlTemplate != ""
}

func Available(dependency *api.Dependency, directory string) (bool, error) {
	path := filepath.Join(directory, dependency.Name)

	_, err := os.Lstat(path)
	if os.IsNotExist(err) {
		return false, nil
	}

	if err != nil {
		return false, errors.Wrapf(err, "failed to check if binary \"%s\" exists at %s", dependency.Name, directory)
	}

	return true, nil
}

func Delete(dependency *api.Dependency, directory string) error {
	path := filepath.Join(directory, dependency.Name)

	err := os.Remove(path)
	if err != nil {
		return errors.Wrapf(err, "failed to delete old binary \"%s\" at %s", dependency.Name, directory)
	}

	return nil
}

func Install(dependency *api.Dependency, directory string) error {
	templateData := templateDataFromDependency(dependency)
	targetUriTemplate, err := targetUriTemplate(dependency)
	if err != nil {
		return errors.Wrapf(err, "failed to find \"%s\" target uri template", dependency.Name)
	}

	targetUri, err := renderUriTempate(targetUriTemplate, templateData)
	if err != nil {
		return errors.Wrapf(err, "failed to render \"%s\" archive uri", dependency.Name)
	}

	tempDirectory, err := os.MkdirTemp(os.TempDir(), dependency.Name)
	if err != nil {
		return errors.Wrapf(err, "failed to create temporary directory")
	}
	defer os.RemoveAll(tempDirectory)

	targetFilename, err := filenameFromUrl(targetUri)
	if err != nil {
		return errors.Wrapf(err, "failed to get \"%s\" archive filename", dependency.Name)
	}

	err = download(targetUri, targetFilename, tempDirectory)
	if err != nil {
		return errors.Wrapf(err, "failed to download \"%s\" archive", dependency.Name)
	}

	if dependency.ChecksumUrlTemplate != "" {
		checksumUri, err := renderUriTempate(dependency.ChecksumUrlTemplate, templateData)
		if err != nil {
			return errors.Wrapf(err, "faile to render \"%s\" checksum uri", dependency.Name)
		}

		checksumFilename, err := filenameFromUrl(targetUri)
		if err != nil {
			return errors.Wrapf(err, "failed to get \"%s\" checksum filename", dependency.Name)
		}

		err = download(checksumUri, checksumFilename, tempDirectory)
		if err != nil {
			return errors.Wrapf(err, "failed to download \"%s\" checksum", dependency.Name)
		}

		err = verify(targetFilename, checksumFilename, dependency.ChecksumAlgorithm, tempDirectory)
		if err != nil {
			return errors.Wrapf(err, "failed to verify \"%s\" archive", dependency.Name)
		}
	}

	if isArchive(dependency) {
		targetFilename, err = decompress(targetFilename, tempDirectory)
		if err != nil {
			return err
		}
	}

	err = moveFile(targetFilename, dependency.Name, tempDirectory, directory)
	if err != nil {
		return errors.Wrapf(err, "failed to move \"%s\" binary", dependency.Name)
	}

	return nil
}
