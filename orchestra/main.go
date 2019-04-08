package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/dhowden/tag"
)

func main() {
	wd, _ := os.Getwd()

	filepath.Walk(".",
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if info.IsDir() {
				return nil
			}

			file, err := os.Open(path)

			if err != nil {
				fmt.Println(err)
			}

			defer file.Close()

			m, err := tag.ReadFrom(file)

			if err != nil {
				fmt.Println(err)
				return nil
			}

			fmt.Println(m.Title())

			return nil
		})

	fmt.Println("Hello World! from", wd)
}
