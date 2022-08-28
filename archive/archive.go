package archive

import (
	"archive/tar"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func Create(input_dir, output_file string) (err error) {

	dist, err := os.Create(output_file)
	if err != nil {
		panic(err)
	}
	defer dist.Close()

	tw := tar.NewWriter(dist)
	defer tw.Close()

	// 再帰的にファイルを取得する
	if err := filepath.Walk(input_dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		intake_path := strings.Replace(path, input_dir, "", 1)
		// write header
		if err := tw.WriteHeader(&tar.Header{
			Name:    intake_path,
			Mode:    int64(info.Mode()),
			ModTime: info.ModTime(),
			Size:    info.Size(),
		}); err != nil {
			return err
		}

		// Writing to tar file
		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()

		if _, err := io.Copy(tw, f); err != nil {
			return err
		}
		return nil

	}); err != nil {
		panic(err)
	}

	return nil
}
