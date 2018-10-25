
package main

import (
    "Crawler"
    "encoding/json"
    "fmt"
    "log"
    "github.com/gorilla/mux"
    "net/http"
)

type crawler_input struct {
    Url string `json:"url"`
    Traversal string `json:"traversal"`
    Keyword string `json:"keyword,omitempty"`
    Depth string `json:"depth"` //need to be string in order to parse from form data, not int
}


func GetCrawl(w http.ResponseWriter, req *http.Request) {
    fmt.Fprintf(w, "Hello");
}

func PostCrawl(w http.ResponseWriter, req *http.Request) {


    decoder := json.NewDecoder(req.Body)


    var input crawler_input
    err := decoder.Decode(&input)

    if err != nil {
        panic(err)
    }
    
    log.Println("url is: " + input.Url)
    log.Println("traversal is: " + input.Traversal)
    log.Println("keyword is: " + input.Keyword)
    log.Println("depth is:", input.Depth)

    if input.Traversal == "bfs" {
        Crawler.BreadthFirstSearchCrawl(input.Url)
    } else {
        Crawler.DepthFirstSearchCrawl(input.Url)
    }
    

}



func main() {
    router := mux.NewRouter()
   
    router.HandleFunc("/crawl", GetCrawl).Methods("GET")

    router.HandleFunc("/crawl", PostCrawl).Methods("POST")

    log.Println("starting server at 12345..\n")
    log.Fatal(http.ListenAndServe(":12345", router))
}
