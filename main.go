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

var pwd string

func init() {
	pwd, _ = os.Getwd()
}

var imp = cmd.Command{
	Name:    "imp",
	Summary: "experimental command line tool for the Internet Computer",
	Commands: []cmd.Command{
		version,
		stats,
		initialize,
		get,
	},
}

var version = cmd.Command{
	Name:    "version",
	Aliases: []string{"v"},
	Summary: "print Imp version",
	Method: func(args []string, _ map[string]string) error {
		if len(args) != 0 {
			return fmt.Errorf("too long")
		}
		fmt.Println("v0.1.0")
		return nil
	},
}

var initialize = cmd.Command{
	Name:    "init",
	Summary: "initialize new module in current directory",
	Args:    []string{"module-name"},
	Method: func(args []string, _ map[string]string) error {
		if len(args) != 1 {
			fmt.Println("expected 1 argument <module-name>")
			return nil
		}
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

var get = cmd.Command{
	Name: "get",
	Args: []string{"module-name", "tag"},
	Method: func(args []string, _ map[string]string) error {
		if len(args) != 2 {
			fmt.Println("expected 2 arguments <module-name> <tag>")
			return nil
		}

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

func main() {
	if len(os.Args) == 1 {
		imp.Help()
		return
	}
	if err := imp.Call(os.Args[1:]...); err != nil {
		imp.Help()
	}
}
