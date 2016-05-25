package main

import (
	"html/template"
	"net/http"
)

//Cache all templates inside templates directory
var templates = template.Must(template.ParseGlob("templates/*"))

//Struct that holds your web pages variables to be used in templates
type webPage struct {
	WelcomeMsg   string
	Introduction string
}

//Handles http requests for url path "/"
func handleHomePage(res http.ResponseWriter, req *http.Request) {

	//Create struct and assign it variables
	homePage := new(webPage)
	homePage.WelcomeMsg = "Welcome to a website!"
	homePage.Introduction = "Golang is the best programming language ever!"

	//res: http.ResponseWriter to write our data back to the client
	//{{ define "landingPage" }} is inside templates/index.html and tells the method which template to use
	//*homePage: Passes the values stored inside our homePage struct to be combined with the template
	templates.ExecuteTemplate(res, "landingPage", *homePage)
}

func main() {

	//serves static css/js ..etc files
	http.HandleFunc("/public/", func(res http.ResponseWriter, req *http.Request) {
		http.ServeFile(res, req, req.URL.Path[1:])
	})

	//Use the function handleHomePage for any requests to "/"
	http.HandleFunc("/", handleHomePage)

	//Start the server on Port 8080
	http.ListenAndServe(":8080", nil)
}
