package main

import (
	"flag"
	"fmt"
	"github.com/babelrpc/swagger2"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

func main() {
	force := flag.String("force", "", "Forces all files to be interpreted as the given type (json or yaml)")
	flag.Parse()
	if *force != "" && *force != "yaml" && *force != "json" {
		fmt.Fprintln(os.Stderr, "The -force option must be json or yaml")
		return
	}
	args := flag.Args()
	if args != nil {
		for _, pattern := range args {
			files, err := filepath.Glob(pattern)
			if err != nil {
				log.Fatal(err)
			}
			for _, f := range files {
				fi, err := os.Stat(f)
				if err != nil {
					log.Fatal(err)
				}
				if !fi.IsDir() {
					var ext string
					if *force != "" {
						ext = "." + *force
					} else {
						ext = filepath.Ext(f)
					}
					if ext == ".yaml" || ext == ".json" {
						fmt.Printf("%s:\n", f)
						b, err := ioutil.ReadFile(f)
						if err != nil {
							fmt.Println("\t", err)
						} else {
							var swag *swagger2.Swagger
							if ext == ".yaml" {
								swag, err = swagger2.LoadYaml(b)
							} else {
								swag, err = swagger2.LoadJson(b)
							}
							if err != nil {
								fmt.Println("\t", err)
							} else {
								errs := swag.Validate()
								if len(errs) > 0 {
									fmt.Println(swagger2.ErrorList(errs).Indent("\t"))
								}
							}
						}
					}
				}
			}
		}
	}
}
