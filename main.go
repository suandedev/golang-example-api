package main

import (
	"fmt"
	"log"
	"io/ioutil"
	"net/http"
	"encoding/json"
	"github.com/gorilla/mux"
)

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Homepage endpoint hit")
	log.Println("enpoint hit: home page")
}

func handleRequests() {
	// creates a new instance of a mux router
	myRouter := mux.NewRouter().StrictSlash(true)
	// replace http.HandleFunc with myRouter.HandleFunc
	myRouter.HandleFunc("/", homePage)
	// add our articles route and map it to our 
    // returnAllArticles function like so
	myRouter.HandleFunc("/articles", returnAllArticle)
	myRouter.HandleFunc("/articles/{id}", returnSingleArticle)
	myRouter.HandleFunc("/articles/create", createNewArticle).Methods("POST")
	myRouter.HandleFunc("/article/{id}", deleteArticle).Methods("DELETE")
	myRouter.HandleFunc("/article/{id}", updateArticle).Methods("PUT")
	 // finally, instead of passing in nil, we want
    // to pass in our newly created router as the second
    // argument
	log.Fatal(http.ListenAndServe(":8081", myRouter))
}

type Article struct {
	Id 	string `json:"Id"`
	Title string `json:"Title"`
	Desc string `json:"desc"`
	Content string `json:"content"`
}


// let's declare a global Articles array
// that we can then populate in our main function
// to simulate a database
var Articles []Article

func returnAllArticle(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: returnAllArticles")
	json.NewEncoder(w).Encode(Articles)
}

func returnSingleArticle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["id"]

	for _, article := range Articles {
		if article.Id == key {
			json.NewEncoder(w).Encode(article)
		}
	}
}

func createNewArticle(w http.ResponseWriter, r *http.Request) {
    // get the body of our POST request
    // unmarshal this into a new Article struct
    // append this to our Articles array.    
	fmt.Println("Endpoint Hit: createNewArticle")
    reqBody, _ := ioutil.ReadAll(r.Body)
    var article Article 
    json.Unmarshal(reqBody, &article)
    // update our global Articles array to include
    // our new Article
    Articles = append(Articles, article)

    json.NewEncoder(w).Encode(article)
}

func deleteArticle(w http.ResponseWriter, r *http.Request) {
	// once again, we will need to parse the path parameters
	vars := mux.Vars(r)
	// we will need to extract the `id` of the article we
    // wish to delete
	id := vars["id"]

	// we then need to loop through all our articles
	for index, article := range Articles {
		 // if our id path parameter matches one of our
        // articles
		if article.Id == id {
			// updates our Articles array to remove the 
            // article
			Articles = append(Articles[:index], Articles[index+1:]...)
		}
	}
}

func updateArticle(w http.ResponseWriter, r *http.Request) {
	// get the body of our POST request
	// unmarshal this into a new Article struct
	// append this to our Articles array.
	reqBody, _ := ioutil.ReadAll(r.Body)
	var article Article
	json.Unmarshal(reqBody, &article)
	// once again, we will need to parse the path parameters
	vars := mux.Vars(r)
	// we will need to extract the `id` of the article we
    // wish to delete
	id := vars["id"]

	// we then need to loop through all our articles
	for index, article := range Articles {
		// if our id path parameter matches one of our
		// articles
		if article.Id == id {
			// updates our Articles array to remove the
			// article
			Articles = append(Articles[:index], Articles[index+1:]...)
		}
	}
}

func main() {
	Articles = []Article{
		Article{Id: "1", Title: "Test Title", Desc: "Test Desc", Content: "Hello world"},
		Article{Id: "2", Title: "Test Title 2", Desc: "Test Desc", Content: "Hello world 2"},
	}
	handleRequests()
}