package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func serveFile(mux *http.ServeMux, path string) {
	prefix := "/files/"
	mux.Handle(prefix, http.StripPrefix(prefix, http.FileServer(http.Dir(path))))
}

func GetCurrentDirectory() *string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	absPath := strings.Replace(dir, "\\", "/", -1)
	return &absPath
}

func execSystemComnand(comm string) string{
	cmd := exec.Command(comm)
	var stdout io.ReadCloser
	var err error
	if stdout, err = cmd.StdoutPipe();err != nil{
		log.Fatal(err)
	}
	defer stdout.Close()
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
	var opBytes []byte
	if opBytes, err = ioutil.ReadAll(stdout); err != nil{
		log.Fatal(err)
	}
	res := string(opBytes)
	return res[0:len(res)-1]
}

func main() {
	port := flag.Int("p", 6371, "your server port, default 6371")
	flag.Parse()
	httpMux := http.NewServeMux()
	currP := execSystemComnand("pwd")
	//fmt.Printf("ss%sss", currP)
	log.Printf("serve files in current path: %s\nsee 127.0.0.1:%d/files\n", currP, *port)
	serveFile(httpMux, currP)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", *port), httpMux); err != nil {
		panic(err)
	}
}
