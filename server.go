package main

import (
	"fmt"
	"net/http"
	"encoding/json"
	"html/template"
	"io"
	"log"
	"os"
)

type FileData struct {
  Name string `json:"filename"`
  Size int64  `json:"size"`
}

func upload(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t, _ := template.ParseFiles("upload.html")
		t.Execute(w, "")
	} else {
		r.ParseMultipartForm(1 << 20)
		file, handler, err := r.FormFile("uploadfile")

		if err != nil {
			fmt.Println("Parse error: ", err)
			return
		}

		defer file.Close()
		f, err := os.Create("./test/" + handler.Filename)

		if err != nil {
			fmt.Println("File create error: ", err)
			return
		}

		io.Copy(f, file)
		size, err := f.Stat()
		fmt.Println(size.Size())

		fileData := FileData{handler.Filename, size.Size()}
		res, err := json.Marshal(fileData)

    if err != nil {
      fmt.Println("JSON marshal error: ", err)
      return
    }

		w.Header().Set("Content-Type", "application/json")
		w.Write(res)

		f.Close()
		os.Remove("./test/" + handler.Filename)
	}
}

func faviconHandler(w http.ResponseWriter, r *http.Request)  {
	http.ServeFile(w, r, "./favicon.ico")
}

func main() {
	http.HandleFunc("/", upload)
	http.HandleFunc("/favicon.ico", faviconHandler)
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		log.Fatal("ListenAndServe Error: ", err)
	}
}
