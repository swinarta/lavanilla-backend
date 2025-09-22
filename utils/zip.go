package utils

import (
	"archive/zip"
	"bytes"
	"fmt"
	"log"
)

type FileData struct {
	Key  string
	Data []byte
	Err  error
}

func CreateZipArchive(results <-chan FileData) (*bytes.Buffer, error) {
	var zipBuf bytes.Buffer
	zipWriter := zip.NewWriter(&zipBuf)

	for result := range results {
		if result.Err != nil {
			log.Printf("failed to download %s: %v\n", result.Key, result.Err)
			continue
		}

		// fw, err := zipWriter.Create(path.Base(result.Key))
		fw, err := zipWriter.Create(result.Key)
		if err != nil {
			log.Printf("failed to create zip entry: %v", err)
			continue
		}

		_, err = fw.Write(result.Data)
		if err != nil {
			log.Printf("failed to write %s to zip: %v\n", result.Key, err)
			continue
		}
	}

	if err := zipWriter.Close(); err != nil {
		return nil, fmt.Errorf("failed to close zip writer: %w", err)
	}
	return &zipBuf, nil
}
