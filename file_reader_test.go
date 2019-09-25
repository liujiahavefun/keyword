package main

import (
    //"regexp"
    "testing"
    //"sort"
    //"fmt"
    "fmt"
)

func TestPDF(t *testing.T) {
    content, err := readPDFFile("liujia.pdf")
    if err != nil {
        t.Error("read pdf file failed, err: ", err)
    }

    fmt.Println(content)
}

/*
func TestWord(t *testing.T) {
    content, err := ReadWordFile("liujia.docx")
    if err != nil {
        t.Error("read word file failed, err: ", err)
    }

    fmt.Println(content)
}

func TestExcel(t *testing.T) {
    content, err := ReadExcelFile("liujia.xlsx")
    if err != nil {
        t.Error("read excel file failed, err: ", err)
    }

    fmt.Println(content)
}

func TestPPT(t *testing.T) {
    content, err := ReadPPTFile("liujia.pptx")
    if err != nil {
        t.Error("read ppt file failed, err: ", err)
    }

    fmt.Println(content)
}
*/