// Ganalyze creates html reports about PE files
package main

import (
	"context"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/masonj188/binanalysis/ganalyze/pinfo"
	"golang.org/x/sync/semaphore"
)

type LinkName struct {
	Link string
	Name string
}
type Pathlink struct {
	LinkNames []LinkName
}

func processFile(path string, outQueue chan<- LinkName, sem *semaphore.Weighted) error {
	sem.Acquire(context.Background(), 1)
	defer sem.Release(1)
	f, err := os.Open(path)
	if err != nil {
		fmt.Println("Unable to open file: ", err)
		outQueue <- (LinkName{})
		return err
	}
	defer f.Close()
	props, err := pinfo.NewProps(f, true)
	if err != nil {
		outQueue <- (LinkName{})
		return err
	}
	outpath := strings.Join([]string{path, ".html"}, "")
	err = props.ExportHTML(outpath)
	if err != nil {
		fmt.Println("Error exporting html for", props.Name)
		outQueue <- (LinkName{})
		return err
	}
	//absPath, err := filepath.Abs(outpath)
	outQueue <- LinkName{outpath, filepath.Base(path)}
	return nil
}

func main() {
	links := Pathlink{}
	fileQueue := make([]string, 0)
	linkchan := make(chan LinkName)
	sem := semaphore.NewWeighted(100)
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

	go func() {
		counter := 0
		for result := range linkchan {
			if result == (LinkName{}) {
				wg.Done()
			} else {
				links.LinkNames = append(links.LinkNames, result)
				wg.Done()
			}
			counter++
			fmt.Printf("\r%d/%d Files Processed", counter, len(fileQueue))
		}
	}()

	for _, path := range fileQueue {
		wg.Add(1)
		go processFile(path, linkchan, sem)
	}

	wg.Wait()
	fmt.Println("")
	t, err := template.New("index").Parse(pinfo.Mainpage)
	//t, err := template.ParseFiles("index.html.template")
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
	err = t.ExecuteTemplate(f, "index", links)
	if err != nil {
		fmt.Println("Error executing index template", err)
	}

}
