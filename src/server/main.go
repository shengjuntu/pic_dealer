// hello.go
package main

import (
	"io"
	"os"
	"fmt"
	"net/http"
	"io/ioutil"
	"log"
)

const (
	UPLOAD_DIR = "./upload"
	REDIS_ADDRESS = ""
)

/*
type ImageAttribute struct {
	Length int32,
	Name string,
	Format string,
	TimeStamp string,
	Status int,
}
*/



var cur_img_id int64 = 0


func check(e error) {
    if e != nil {
        panic(e)
    }
}

func helloHandler(w http.ResponseWriter, req *http.Request) {
	dat,err := ioutil.ReadFile("hello.html")
	check(err)
	io.WriteString(w, string(dat))
}

func viewHandler(w http.ResponseWriter, req *http.Request) {
	dat,err := ioutil.ReadFile("view.html")
	check(err)
	io.WriteString(w, string(dat))
}

func get_cur_img_id()(string) {
	return fmt.Sprintf("raw_img_%d", cur_img_id)	
}

func get_img_id(id string)(string) {
	return "raw_img_" + id
}


func getImageHandler(w http.ResponseWriter, req *http.Request) {
	str_id := req.URL.Query().Get("id")
	log.Println("--------" + req.Method)
	log.Println("str_id:" + str_id)
		 
	if str_id == "" {
		str_id = get_cur_img_id()
	}
	dat, err := db_get_image(str_id)
	check(err)
	w.Write(dat)
}

func listHandler(w http.ResponseWriter, req *http.Request) {
	keys, err := db_list_image()
	check(err)
	
	io.WriteString(w, "<html><body>")
	for _,key := range keys {
		io.WriteString(w, "<li>")
		io.WriteString(w, key)
		img_elem := fmt.Sprintf("<img src=\"images?id=%s\"/>",key)
		io.WriteString(w, img_elem)
		io.WriteString(w, "</li>")
	}
	io.WriteString(w, "</body></html>")

}

func uploadHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method == "GET" {
		dat,err := ioutil.ReadFile("upload.html")
		check(err)
		io.WriteString(w, string(dat))
		return
	} 
	if req.Method == "POST" {
		f,h,err := req.FormFile("image")
		if err != nil {
			http.Error(w, err.Error(),
				http.StatusInternalServerError)
			return
		}
		filename := h.Filename
		defer f.Close()
		
		bytes,err := ioutil.ReadAll(f)
		
		cur_img_id ++
		db_insert_image(bytes, get_cur_img_id())
		
		
		t,err := os.Create(UPLOAD_DIR + "/" + filename)
		if err != nil {
			http.Error(w, err.Error(),
				http.StatusInternalServerError)
		}
		defer t.Close()
		
		if _, err := io.Copy(t, f); err != nil {
			http.Error(w, err.Error(),
				http.StatusInternalServerError)
			return
		}
		
		io.WriteString(w, "upload OK")
	}
}

func main() {
		
	init_pool()
	_id, err := db_get_image_id()
	check(err)
	
	cur_img_id = _id
	log.Printf("INIT OK cur_img_id:%d", cur_img_id)
	
	http.HandleFunc("/hello", helloHandler)
	http.HandleFunc("/upload", uploadHandler)
	http.HandleFunc("/images", getImageHandler)
	http.HandleFunc("/view", viewHandler)
	http.HandleFunc("/list", listHandler)

	err = http.ListenAndServe(":8082", nil)
	if err != nil {
		log.Fatal("ListenAndServer: ", err.Error())
	}
	
	log.Println("exiting")
}
