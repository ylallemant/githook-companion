package dependency

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

func moveFile(decompressedFilename, binaryFilename, sourceDirectory, targetDirectory string) error {
	sourcePath := filepath.Join(sourceDirectory, decompressedFilename)
	targetPath := filepath.Join(targetDirectory, binaryFilename)

	content, err := os.ReadFile(sourcePath)
	if err != nil {
		return errors.Wrapf(err, "failed to read source file %s", sourcePath)
	}

	err = os.Remove(targetPath)
	if err != nil && !os.IsNotExist(err) {
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
