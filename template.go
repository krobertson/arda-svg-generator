package main

import (
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"text/template"
)

// funcMap is a list od functions to expose to the templates
var funcMap = template.FuncMap{
	"embed":      tplEmbed,
	"blankOrNum": tplBlankOrNum,
	"offset":     tplOffset,
}

// docFuncMap is a list od functions to expose to the templates
var docFuncMap = template.FuncMap{
	"hasEntry": tplHasEntry,
}

// generateEntry is used to generate the specified tmplFile to the specified outFile
// with the item record.
func generateEntry(item *Item, tmplFile, outFile string) error {
	f, err := os.Create(outFile)
	if err != nil {
		return err
	}
	defer f.Close()

	tname := filepath.Base(tmplFile)
	tpl, err := template.New("").Funcs(funcMap).ParseFiles(tmplFile)
	if err != nil {
		return err
	}

	return tpl.ExecuteTemplate(f, tname, item)
}

// tplEmbed is used to embed the specified image in the data URI format.
// Currently, it assumes image/png format.
func tplEmbed(filename string) string {
	b, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Embed error: %q\n", err)
	}
	return "data:image/png;base64," + base64.StdEncoding.EncodeToString(b)
}

// tplBlankOrNum is used to return a blank string when the specified number is
// zero, or return a formatted string when the number is set.
func tplBlankOrNum(num *int) string {
	if num == nil {
		return ""
	} else {
		return strconv.Itoa(*num)
	}
}

// tplOffset can be used to apply an offset for a location when multiple
// elements are shown.
func tplOffset(base, offset, pos int) int {
	return base + (offset * pos)
}

// generateDocument is used to generate the printable document containing
// multiple entries.
func generateDocument(entries []string, tmplFile, outFile string) error {
	f, err := os.Create(outFile)
	if err != nil {
		return err
	}
	defer f.Close()

	tname := filepath.Base(tmplFile)
	tpl, err := template.New("").Funcs(docFuncMap).ParseFiles(tmplFile)
	if err != nil {
		return err
	}

	return tpl.ExecuteTemplate(f, tname, entries)
}

// tplHasEntry can be used to check if the index exists in the array.
func tplHasEntry(entries []string, i int) bool {
	return len(entries) > i
}
