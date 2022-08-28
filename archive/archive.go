package archive

import (
	"archive/tar"
	"io"
	"io/fs"
	"path/filepath"
	"os"
	"strings"

	"github.com/facebookgo/symwalk"
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
	if err := symwalk.Walk(input_dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		intake_path := strings.Replace(path, input_dir, "", 1)
		tar_mode := info.Mode()
		tar_size := info.Size()

		// To read symbolic links
		if info.Mode() & fs.ModeSymlink != 0 {
			path_sym, err := filepath.EvalSymlinks(path)
			if err != nil {
				return err
			}
			path = path_sym
			stat, err := os.Stat(path_sym)
			if err != nil {
				return err
			}
			tar_mode = stat.Mode()
			tar_size = stat.Size()
		}

		// write header
		if err := tw.WriteHeader(&tar.Header{
			Name:    intake_path,
			Mode:    int64(tar_mode),
			ModTime: info.ModTime(),
			Size:    tar_size,
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
