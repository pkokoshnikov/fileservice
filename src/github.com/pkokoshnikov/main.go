package main

import (
    "fmt"
    "net/http"
	"log"
	"github.com/pkokoshnikov/fs"

)

func uploadServiceHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("UploadServiceHandler is called")
	
	defer func(){
    	if rec := recover(); rec != nil {
	        log.Print("Recovered in uploadServiceHandler, message = ", rec)
	        
	        switch rec.(type) {
	        case string:
				err := rec.(string)
				http.Error(w, err, http.StatusInternalServerError)	            
	        case error:
				err := rec.(error)
				http.Error(w, err.Error(), http.StatusInternalServerError)	            
	        default:
	            http.Error(w, "Unknown error", http.StatusInternalServerError)	            
	        }
    	}
	}()
	
	switch r.Method {		
		// Show all files
		case "GET":
			log.Print("GET is called")
			fmt.Fprintf(w, "Files:\n")
			dao := fs.NewDAO()
			l := dao.ShowAllFiles()
			for e := l.Front(); e != nil; e = e.Next() {
				fmt.Fprintf(w, e.Value.(fs.File).Filename + "\n")
			}
			
		// Upload a new file.    
		case "POST":
			log.Print("POST is called")
			dao := fs.NewDAO()
			err := dao.UploadFile(r)
			if err != nil {
				log.Print(err.Error())
			}
			fmt.Fprintf(w, "New file was posted\n")
			
		// Remove a file.
		case "DELETE":
			log.Print("DELETE is called")			
			dao := fs.NewDAO()
			filename := r.FormValue("fileName")
			dao.DeleteFile(filename)
			
			fmt.Fprintf(w, "Delete a file\n")
			
		default:
		   	http.Error(w, "Unknown operation", http.StatusNotFound)	             
	}
}

func viewHandler(w http.ResponseWriter, r *http.Request) {	
	log.Print("ViewHandler is called")
	title := r.URL.Path[len("/view/"):]
	if len(title) == 0 {
		title = "index"
	}

	http.ServeFile(w,r, title + ".html");
}

func main() {
    http.HandleFunc("/upload/service", uploadServiceHandler)
	http.HandleFunc("/view/", viewHandler)
	http.ListenAndServe(":8080", nil)	
}
