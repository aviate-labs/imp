package main

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/aviate-labs/imp/internal/cmd"
)

var get = cmd.Command{
	Name: "get",
	Args: []string{"module-name", "tag"},
	Method: func(args []string, _ map[string]string) error {
		if _, err := os.Stat(pwd + "/mo.mod"); err != nil {
			fmt.Println("need to initalize first")
			return nil
		}

		if _, err := os.Stat(pwd + "/.mod"); err != nil {
			if err = os.Mkdir(pwd+"/.mod", 0755); err != nil {
				return err
			}
		}

		resp, err := http.Get(fmt.Sprintf("https://%s/archive/refs/tags/%s.tar.gz", args[0], args[1]))
		if err != nil {
			fmt.Println(err)
			return nil
		}
		gzr, err := gzip.NewReader(resp.Body)
		if err != nil {
			panic(err)
		}
		defer gzr.Close()

		r := tar.NewReader(gzr)
		for {
			h, err := r.Next()
			if err != nil {
				if err == io.EOF {
					break
				}
				return err
			}
			if h == nil {
				continue
			}
			target := filepath.Join(pwd, ".mod", h.Name)
			switch h.Typeflag {
			case tar.TypeDir:
				if _, err := os.Stat(target); err != nil {
					if err := os.MkdirAll(target, 0755); err != nil {
						return err
					}
				}
			case tar.TypeReg:
				f, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(h.Mode))
				if err != nil {
					return err
				}
				if _, err := io.Copy(f, r); err != nil {
					return err
				}
				f.Close()
			}
		}

		f, err := os.OpenFile(pwd+"/mo.mod", os.O_APPEND|os.O_WRONLY, 0600)
		if err != nil {
			return err
		}
		defer f.Close()
		if _, err := f.WriteString(fmt.Sprintf("require %s %s", args[0], args[1])); err != nil {
			return err
		}
		return nil
	},
}
