package main

import (
	"github.com/tlinden/up/upd/api"
	"log"
	"os"
)

func main() {
	if len(os.Args) > 2 {
		dir, err := os.Getwd()
		if err != nil {
		}
		if err := api.ZipSource(dir+"/"+os.Args[1], os.Args[2]); err != nil {
			log.Fatal(err)
		}
	}
}
