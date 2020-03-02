package main

import (
	"fmt"
	"log"
	"os"

	"github.com/masonj188/binanalysis/ganalyze/pinfo"
)

func main() {
	f, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatalln("Error opening executable: ", err)
	}
	defer f.Close()

	props := pinfo.NewProps(f)
	err = props.ExportHTML()
	if err != nil {
		fmt.Printf("Error exporting html doc for %s\n", props.Name)
	}

	//fmt.Printf("%#v", props)

}
