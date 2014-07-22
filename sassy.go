package sassy

/*
#cgo LDFLAGS: -lstdc++

#include <stdlib.h>
#include "sass_interface.h"
*/
import "C"

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"unsafe"
)

type FileSet struct {
	IncludeDir   []string
	Style        OutputStyle
	ShowComments bool
	files        map[string]*File
}

func (fs *FileSet) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for name, file := range fs.files {
		if strings.HasSuffix(r.URL.Path, name) {
			w.Header().Add("Content-Type", "text/css; charset=utf-8")
			io.WriteString(w, file.Output)
			return
		}
	}
	w.WriteHeader(404)
}

type OutputStyle int

const (
	NestedStyle OutputStyle = iota
	ExpandedStyle
	CompactStyle
	CompressedStyle
)

func (fs *FileSet) ParseFile(filename string) (*File, error) {
	ctx := C.sass_new_file_context()
	ctx.input_path = C.CString(filename)
	defer C.free(unsafe.Pointer(ctx.input_path))

	ctx.options.output_style = C.int(int(fs.Style))
	if fs.ShowComments {
		ctx.options.source_comments = C.int(1)
	} else {
		ctx.options.source_comments = C.int(0)
	}

	C.sass_compile_file(ctx)

	fr := &File{
		Filename: filename,
		Output:   C.GoString(ctx.output_string),
	}

	es := C.GoString(ctx.error_message)
	if es != "" {
		return nil, fmt.Errorf(es)
	}
	C.sass_free_file_context(ctx)

	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	fr.Source = string(b)

	if fs.files == nil {
		fs.files = make(map[string]*File)
	}
	fs.files[cssize(filename)] = fr
	return fr, nil
}

func (fs *FileSet) Parse(filename, content string) (*File, error) {
	ctx := C.sass_new_context()
	ctx.source_string = C.CString(content)
	defer C.free(unsafe.Pointer(ctx.source_string))

	ctx.options.output_style = C.int(int(fs.Style))
	if fs.ShowComments {
		ctx.options.source_comments = C.int(1)
	} else {
		ctx.options.source_comments = C.int(0)
	}

	C.sass_compile(ctx)

	fr := &File{
		Filename: filename,
		Output:   C.GoString(ctx.output_string),
	}

	es := C.GoString(ctx.error_message)
	if es != "" {
		return nil, fmt.Errorf(es)
	}
	C.sass_free_context(ctx)

	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	fr.Source = string(b)

	if fs.files == nil {
		fs.files = make(map[string]*File)
	}
	fs.files[cssize(filename)] = fr

	return fr, nil
}

type File struct {
	Source, Output string
	Filename       string
}

func cssize(s string) string {
	return strings.Replace(s, ".scss", ".css", 1)
}
