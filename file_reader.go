package main

import (
    "os"
    "fmt"
    "bytes"
    "strings"
    "errors"
    "path"
    "github.com/unidoc/unipdf/extractor"
    pdf "github.com/unidoc/unipdf/model"
    "github.com/unidoc/unioffice/document"
    "github.com/unidoc/unioffice/spreadsheet"
    "github.com/unidoc/unioffice/presentation"
    "github.com/deckarep/golang-set"
)

var (
    supported_files mapset.Set
)

func init() {
    supported_files = mapset.NewSet()
    supported_files.Add("pdf")
    supported_files.Add("docx")
    supported_files.Add("doc")
    supported_files.Add("xlsx")
    supported_files.Add("xls")
    supported_files.Add("pptx")
    supported_files.Add("ppt")
}

type FileReader interface {
    Support() bool
    Read() (data string, err error)
}

type FileReaderImpl struct {
    FileReader
    file string
}

func NewFileReader(file string) FileReader {
    return &FileReaderImpl{
        file : file,
    }
}

func (self *FileReaderImpl) Support() bool {
    fileext := strings.ToLower(strings.Replace(path.Ext(self.file), ".", "", -1))
    return supported_files.Contains(fileext)
}

func (self *FileReaderImpl) Read() (data string, err error) {
    fileext := strings.ToLower(strings.Replace(path.Ext(self.file), ".", "", -1))
    switch fileext {
    case "pdf":
        return readPDFFile(self.file)
    case "doc", "docx":
        return ReadWordFile(self.file)
    case "xls", "xlsx":
        return ReadExcelFile(self.file)
    case "ppt", "pptx":
        return ReadPPTFile(self.file)
    }
    return data, errors.New("do not support encrypted pdf file")
}

func readPDFFile(inputPath string) (data string, err error) {
    f, err := os.Open(inputPath)
    if err != nil {
        return data, err
    }

    defer f.Close()

    pdfReader, err := pdf.NewPdfReader(f)
    if err != nil {
        return data, err
    }

    isEncrypted, err := pdfReader.IsEncrypted()
    if err != nil {
        return data, err
    }

    if isEncrypted {
        /*
        _, err = pdfReader.Decrypt([]byte("it should be replaced by real password..."))
        if err != nil {
            return data, err
        }
        */
        return data, errors.New("do not support encrypted pdf file")
    }

    numPages, err := pdfReader.GetNumPages()
    if err != nil {
        return data, err
    }

    var buffer bytes.Buffer

    for i := 0; i < numPages; i++ {
        pageNum := i + 1
        page, err := pdfReader.GetPage(pageNum)
        if err != nil {
            continue
        }

        ex, err := extractor.New(page)
        if err != nil {
            continue
        }

        text, err := ex.ExtractText()
        if err != nil {
            continue
        }

        /*
        fmt.Println("------------------------------")
        fmt.Printf("Page %d:\n", pageNum)
        fmt.Printf("\"%s\"\n", text)
        fmt.Println("------------------------------")
        */

        buffer.WriteString(text)
    }

    return buffer.String(), nil
}

func ReadWordFile(file string) (data string, err error) {
    doc, err := document.Open(file)
    if err != nil {
        return data, err
    }
    var buffer bytes.Buffer
    for _, para := range doc.Paragraphs() {
        for _, run := range para.Runs() {
            buffer.WriteString(run.Text() + "\t")
        }
    }
    return buffer.String(), nil
}

func ReadExcelFile(file string) (data string, err error) {
    excel, err := spreadsheet.Open(file)
    if err != nil {
        return data, err
    }
    var buffer bytes.Buffer
    for _, sheet := range excel.Sheets() {
        fmt.Println("sheet: ", sheet.Name())
        for _, row := range sheet.Rows() {
            for _, cell := range row.Cells() {
                buffer.WriteString(strings.TrimSpace(cell.GetString()) + "\t")
            }
        }
    }
    return buffer.String(), nil
}

func ReadPPTFile(file string) (data string, err error) {
    ppt, err := presentation.Open(file)
    if err != nil {
        return data, err
    }
    var buffer bytes.Buffer
    for _, slide := range ppt.Slides() {
        for _, placeholder := range slide.PlaceHolders() {
            for _, paragraph := range placeholder.Paragraphs() {
                for _, textrun := range paragraph.X().EG_TextRun {
                    /*
                    if textrun.Fld != nil {
                        fmt.Println(textrun.Fld.T)
                    }*/
                    if textrun.R != nil {
                        buffer.WriteString(textrun.R.T + "\t")
                    }
                    /*
                    if textrun.Br != nil {
                        fmt.Println(textrun.Br.RPr)
                    }*/
                }
            }
        }
    }
    return buffer.String(), nil
}