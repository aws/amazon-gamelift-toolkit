package tools

import (
	"archive/zip"
	"context"
	"fmt"
	"strings"

	"github.com/aws/amazon-gamelift-toolkit/fast-build-update-tool/internal/gamelift"
)

type ZipValidator struct {
	buildZipPath string
}

func NewZipValidator(buildZipPath string) *ZipValidator {
	return &ZipValidator{buildZipPath: buildZipPath}
}

// ValidateZip will validate that the zip file provided is valid for the given fleet
func (z *ZipValidator) ValidateZip(ctx context.Context, fleet *gamelift.Fleet) error {
	zipReader, zipErr := zip.OpenReader(z.buildZipPath)
	if zipErr != nil {
		return fmt.Errorf("error opening zip file %w", zipErr)
	}
	defer zipReader.Close()

	for _, executablePath := range fleet.ExecutablePaths {
		normalizedFileToFind := strings.ReplaceAll(executablePath, "C:\\game\\", "")
		normalizedFileToFind = strings.ReplaceAll(normalizedFileToFind, "/local/game/", "")
		normalizedFileToFind = strings.ReplaceAll(normalizedFileToFind, "\\", "/")

		if !isFileInZip(zipReader, normalizedFileToFind) {
			return fmt.Errorf("zip file does not contain executable %s", normalizedFileToFind)
		}
	}

	return nil
}

func isFileInZip(zipReader *zip.ReadCloser, fileToFind string) bool {
	for _, file := range zipReader.File {
		if file.Name == fileToFind {
			return true
		}
	}
	return false
}
