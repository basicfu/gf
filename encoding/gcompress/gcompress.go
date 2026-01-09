package gcompress

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func UnZipContent(data []byte, path string) error {
	r, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return fmt.Errorf("zip reader error: %w", err)
	}
	if err := os.MkdirAll(path, 0755); err != nil {
		return err
	}
	for _, f := range r.File {
		targetPath := filepath.Join(path, f.Name)
		if !strings.HasPrefix(
			filepath.Clean(targetPath),
			filepath.Clean(path)+string(os.PathSeparator),
		) {
			return fmt.Errorf("illegal file path: %s", f.Name)
		}
		if f.FileInfo().IsDir() {
			if err := os.MkdirAll(targetPath, f.Mode()); err != nil {
				return err
			}
			continue
		}
		if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
			return err
		}
		rc, err := f.Open()
		if err != nil {
			return err
		}
		out, err := os.OpenFile(targetPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, f.Mode())
		if err != nil {
			rc.Close()
			return err
		}
		_, err = io.Copy(out, rc)
		rc.Close()
		out.Close()
		if err != nil {
			return err
		}
	}
	return nil
}
