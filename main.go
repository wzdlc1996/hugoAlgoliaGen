package main

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	tarDir := os.Args[1]
	infos, err := CreateIndex(tarDir)
	if err != nil {
		fmt.Println("Error, exit")
		os.Exit(-1)
	}

	data, _ := json.MarshalIndent(infos, "", "\t")

	os.WriteFile("./out.json", data, 0644)

}

func CreateIndex(workDir string) ([]PostInfo, error) {
	indexList := []PostInfo{}
	err := filepath.Walk(workDir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("Failed at path:\t%q: %v\n", path, err)
			return err
		}
		if info.IsDir() && info.Name() == "skip" {
			return filepath.SkipDir
		}
		// Indexing all markdown file
		if filepath.Ext(path) == ".md" {
			fmt.Println("[Complete]:", path)
			relPath := strings.SplitN(path, "content/", 2)[1]
			relPathItem := strings.Split(relPath, string(os.PathSeparator))
			relPathItem = relPathItem[:len(relPathItem)-1]

			if bsname := filepath.Base(path); bsname != "index.md" {
				relPathItem = append(relPathItem, bsname[:len(bsname)-3])
			}

			mdfile, err := os.ReadFile(path)
			if err != nil {
				fmt.Printf("Failed to read %q: %v\n", path, err)
				return err
			}

			postData := PostParser(string(mdfile))
			postData.Uri = strings.Join(relPathItem, "/")
			postData.ObjectID = strings.Join(relPathItem, "/")
			indexList = append(indexList, postData)
		}
		return nil
	})
	if err != nil {
		fmt.Println("Failed to indexing. Exit with empty")
		return []PostInfo{}, err
	}
	return indexList, nil
}
