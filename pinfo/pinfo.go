package pinfo

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"debug/pe"
	"encoding/hex"
	"fmt"
	"html/template"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
)

// Consts for Magic numbers and Machine Codes
const (
	BIT64 = 0x8664
	BIT32 = 0x14c
	PE32  = 0x10B
	PE32P = 0x20B
)

// BasicProps contains basic data about a PE file
type BasicProps struct {
	Name   string
	MD5    string
	SHA1   string
	SHA256 string
	//Imphash      string
	//SSDEEP       string
	FileType   string
	Magic      string
	FSize      string
	Libraries  []string
	Symbols    []string
	Sections   []pe.Section
	UsingModel bool
	ModelRes   bool
}

// NewProps returns a pointer to a basicProps struct
func NewProps(file *os.File, useModel bool) (*BasicProps, error) {
	props := BasicProps{}
	props.Name = filepath.Base(file.Name())
	err := props.fillHashes(file)
	if err != nil {
		return nil, err
	}
	err = props.fillFileType(file)
	if err != nil {
		return nil, err
	}
	err = props.fillMagic(file)
	if err != nil {
		return nil, err
	}
	err = props.fillFileSize(file)
	if err != nil {
		return nil, err
	}
	err = props.fillLibraries(file)
	if err != nil {
		return nil, err
	}
	err = props.fillSymbols(file)
	if err != nil {
		return nil, err
	}
	err = props.fillSections(file)
	if err != nil {
		return nil, err
	}
	/*
		props.UsingModel = useModel
		if useModel {
			props.fillFromModel(file)
		}
	*/

	return &props, nil
}

func (p *BasicProps) fillHashes(file *os.File) error {
	mh := md5.New()
	_, err := io.Copy(mh, file)
	if err != nil {
		return err
	}
	mhbytes := mh.Sum(nil)
	p.MD5 = hex.EncodeToString(mhbytes[:])

	s1h := sha1.New()
	file.Seek(0, 0)
	_, err = io.Copy(s1h, file)
	if err != nil {
		log.Fatalln("error copying file into sha1 hash:", err)
		return fmt.Errorf("Error copying file into sha1 hash: %v", err)
	}
	s1hbytes := s1h.Sum(nil)
	p.SHA1 = hex.EncodeToString(s1hbytes)

	s2h := sha256.New()
	file.Seek(0, 0)
	_, err = io.Copy(s2h, file)
	if err != nil {
		return fmt.Errorf("Error copying file into sha256 hash: %v", err)
	}
	s2hbytes := s2h.Sum(nil)
	p.SHA256 = hex.EncodeToString(s2hbytes)

	file.Seek(0, 0)
	return nil
}

func (p *BasicProps) fillFileType(file *os.File) error {
	exe, err := pe.NewFile(file)
	if err != nil {
		return fmt.Errorf("Error converting %s to PE: %v", p.Name, err)
	}
	if exe.Machine == BIT32 {
		p.FileType = "Win32 Exe"
	} else if exe.Machine == BIT64 {
		p.FileType = "Win64 Exe"
	} else {
		p.FileType = "Unknown"
	}
	return nil
}

func (p *BasicProps) fillMagic(file *os.File) error {
	exe, err := pe.NewFile(file)
	if err != nil {
		return fmt.Errorf("Error converting %s to PE: %v", p.Name, err)
	}

	switch exe.OptionalHeader.(type) {
	case *pe.OptionalHeader32:
		magic := exe.OptionalHeader.(*pe.OptionalHeader32).Magic
		if magic == PE32 {
			p.Magic = "PE32"
		} else if magic == PE32P {
			p.Magic = "PE32P"
		} else {
			p.Magic = "Unknown"
		}
	case *pe.OptionalHeader64:
		magic := exe.OptionalHeader.(*pe.OptionalHeader64).Magic
		if magic == PE32 {
			p.Magic = "PE32"
		} else if magic == PE32P {
			p.Magic = "PE32P"
		} else {
			p.Magic = "Unknown"
		}
	default:
		p.Magic = "Unknown"
	}
	return nil
}

func (p *BasicProps) fillFileSize(f *os.File) error {
	info, err := f.Stat()
	if err != nil {
		return fmt.Errorf("Error calling stat on %s: %v", p.Name, err)
	}
	p.FSize = strconv.FormatInt(info.Size(), 10)
	return nil
}

func (p *BasicProps) fillSymbols(f *os.File) error {
	exe, err := pe.NewFile(f)
	if err != nil {
		return fmt.Errorf("Error converting %s to PE in fillSymbols: %v", p.Name, err)
	}

	p.Symbols, err = exe.ImportedSymbols()
	if err != nil {
		return fmt.Errorf("Error getting imported symbols for %s: %v", p.Name, err)
	}
	return nil
}

func (p *BasicProps) fillLibraries(f *os.File) error {
	exe, err := pe.NewFile(f)
	if err != nil {
		return fmt.Errorf("Error converting %s to PE in fillImports: %v", p.Name, err)
	}

	p.Libraries, err = exe.ImportedLibraries()
	if err != nil {
		return fmt.Errorf("Error getting imported libraries for %s: %v", p.Name, err)
	}
	return nil
}

func (p *BasicProps) fillSections(f *os.File) error {
	exe, err := pe.NewFile(f)
	if err != nil {
		return fmt.Errorf("Error converting %s to PE in fillSections: %v", p.Name, err)
	}
	for _, val := range exe.Sections {
		p.Sections = append(p.Sections, *val)
	}
	return nil
}

func (p *BasicProps) fillFromModel(f *os.File) {
	fullPath, err := filepath.Abs(f.Name())
	pmodel := exec.Command("python3", "prediction.py", fullPath)
	pmodel.Dir = "python"
	out, err := pmodel.Output()
	if err != nil {
		fmt.Println("Error running model", err)
	}
	res, err := strconv.Atoi(string(out[:1]))
	if err != nil {
		fmt.Println("Error converting stdout to an int", err)
	}
	switch res {
	case 0:
		p.ModelRes = false
	case 1:
		p.ModelRes = true
	case -1:
		fmt.Println("Error from the model")
	}
}

func (p *BasicProps) String() string {
	return fmt.Sprintf("---Basic Info---\n%-15s%s\n%-15s%s\n%-15s%s\n%-15s%s\n%-15s%s\n%-15s%s", "MD5 Hash: ", p.MD5, "SHA1 Hash: ", p.SHA1, "SHA256 Hash: ", p.SHA256, "File Type: ", p.FileType, "Magic: ", p.Magic, "File Size: ", p.FSize)
}

func (p *BasicProps) ExportHTML(outfilePath string) error {
	t, err := template.New("infopage").Parse(bindata)
	//t, err := template.ParseFiles("binpage.html.template")
	if err != nil {
		fmt.Println("Error parsing binpage template: ", err)
		return err
	}
	path := filepath.Dir(outfilePath)
	os.MkdirAll(path, 0755)
	if err != nil {
		fmt.Println("Error getting working directory")
	}
	f, err := os.Create(outfilePath)
	if err != nil {
		fmt.Println("Error creating file", err)
		return fmt.Errorf("Error creating %s: %v", p.Name, err)
	}
	defer f.Close()
	//err = t.ExecuteTemplate(f, "binpage.html.template", *p)
	err = t.ExecuteTemplate(f, "infopage", *p)
	if err != nil {
		return fmt.Errorf("Error executing template for %s: %v", p.Name, err)
	}
	return nil
}
