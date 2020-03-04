package main

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"

	"github.com/masonj188/binanalysis/ganalyze/pinfo"
)

type Pathlink struct {
	LinkName []struct {
		Link string
		Name string
	}
}

func main() {
	links := Pathlink{}
	err := filepath.Walk(os.Args[1], func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if ext := filepath.Ext(path); ext != ".exe" && ext != ".dll" {
			return nil
		}
		f, err := os.Open(path)
		if err != nil {
			fmt.Println("Unable to open file:", err)
		}
		defer f.Close()
		props := pinfo.NewProps(f, true)
		outpath := strings.Join([]string{"report/", path, ".html"}, "")
		err = props.ExportHTML(outpath)
		if err != nil {
			fmt.Println("Error exporting html for", props.Name)
		}
		links.LinkName = append(links.LinkName, struct {
			Link string
			Name string
		}{strings.TrimPrefix(outpath, "report/"), filepath.Base(path)})
		return nil
	})
	if err != nil {
		fmt.Println("Error walking filepath", err)
	}

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
