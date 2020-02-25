package main

import (
	"debug/pe"
	"fmt"
	"log"
)

func main() {
	exe, err := pe.Open("rufus-3.8.exe")
	if err != nil {
		log.Fatalln("Error opening executable: ", err)
	}
	defer exe.Close()
	fmt.Println(exe)

	imports, err := exe.ImportedLibraries()
	fmt.Println("Imports: ", imports)

}
