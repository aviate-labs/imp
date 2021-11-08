package main

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/aviate-labs/imp/internal/cmd"
)

var modCommand = cmd.Command{
	Name:    "mod",
	Summary: "a Motoko package manager",
	Commands: []cmd.Command{
		initializeCommand,
		getCommand,
	},
}

var initializeCommand = cmd.Command{
	Name:    "init",
	Summary: "initialize new module in current directory",
	Args:    []string{"module-name"},
	Method: func(args []string, _ map[string]string) error {
		if _, err := os.Stat(pwd + "/mo.mod"); err == nil {
			fmt.Println("mo.mod already exists")
			return nil
		}
		f, err := os.Create(pwd + "/mo.mod")
		if err != nil {
			panic(err)
		}
		defer f.Close()
		f.WriteString(fmt.Sprintf("module %s\n\n", args[0]))
		return nil
	},
}

var getCommand = cmd.Command{
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
		if resp.StatusCode != 200 {
			fmt.Printf("could not get package: %s@%s\n", args[0], args[1])
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

		mod, err := ioutil.ReadFile("./mo.mod")
		if err != nil {
			return err
		}
		dep := fmt.Sprintf("require %s %s\n", args[0], args[1])
		if !strings.Contains(string(mod), dep) {
			mod = append(mod, []byte(dep)...)
		}
		if err = ioutil.WriteFile("./mo.mod", mod, 0666); err != nil {
			return nil
		}
		return nil
	},
}
