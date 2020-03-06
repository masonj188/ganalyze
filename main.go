package main

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/masonj188/binanalysis/ganalyze/pinfo"
)

type LinkName struct {
	Link string
	Name string
}
type Pathlink struct {
	LinkNames []LinkName
}

func processFile(path string, outQueue chan<- LinkName) error {
	f, err := os.Open(path)
	if err != nil {
		fmt.Println("Unable to open file: ", err)
		return err
	}
	defer f.Close()
	props := pinfo.NewProps(f, true)
	outpath := strings.Join([]string{"report/", path, ".html"}, "")
	err = props.ExportHTML(outpath)
	if err != nil {
		fmt.Println("Error exporting html for", props.Name)
		return err
	}
	outQueue <- LinkName{strings.TrimPrefix(outpath, "report/"), filepath.Base(path)}
	return nil
}

func main() {
	links := Pathlink{}
	fileQueue := make([]string, 0)
	linkchan := make(chan LinkName)
	var wg sync.WaitGroup

	err := filepath.Walk(os.Args[1], func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if ext := filepath.Ext(path); ext != ".exe" && ext != ".dll" {
			return nil
		}
		fileQueue = append(fileQueue, path)
		return nil
	})

	if err != nil {
		fmt.Println("Error walking filepath", err)
	}

	for _, file := range fileQueue {
		wg.Add(1)
		go processFile(file, linkchan)
	}

	go func() {
		for result := range linkchan {
			links.LinkNames = append(links.LinkNames, result)
			wg.Done()
		}
	}()
	wg.Wait()

	t, err := template.ParseFiles("index.html.template")
	if err != nil {
		fmt.Println("Error parsing template", err)
	}
	os.Chdir("report")
	f, err := os.Create("index.html")
	if err != nil {
		fmt.Println("Error creating file", err)
		os.Exit(1)
	}
	defer f.Close()
	err = t.ExecuteTemplate(f, "index.html.template", links)
	if err != nil {
		fmt.Println("Error executing index template", err)
	}

}
