package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/acsellers/sassy"
)

func main() {
	fmt.Println("Starting parse")
	fs := &sassy.FileSet{}
	_, err := fs.Parse("thing.scss", ".name {\ncolor: white;\n}")
	if err != nil {
		panic(err)
	}
	fmt.Println("Starting Server")
	http.Handle("/", fs)
	log.Fatal(http.ListenAndServe(":8989", nil))
}
