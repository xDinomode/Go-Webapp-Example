package main

import (
	"bytes"
	"fmt"
	"github.com/russross/blackfriday"
	"gopkg.in/yaml.v2"
	"html/template"
	"io/ioutil"
	"net/http"
)

//Parse all files inside /templates into var templates
var templates = template.Must(template.ParseGlob("templates/*"))

//Handle web posts (/posts)
type webPost struct {
	Title       string
	Author      string
	Description string
	Body        template.HTML
}

//Grabs yaml from the markdown and inserts into a type webPost (Title, Author, etc..)
func (post *webPost) Parse(markdownYAML []byte) error {

	//Inserts yaml into type webPost and returns nil if no error
	return yaml.Unmarshal(markdownYAML, post)
}

//Handle main webPages (index.html, about.html, contact.html, etc..)
type webPage struct {
	WelcomeMsg   string
	Introduction string
}

//Parse markdown and serve to client from templates
func handlePost(res http.ResponseWriter, req *http.Request) {
	//Read in markdown
	input, _ := ioutil.ReadFile("test.md")

	//Split the []byte input into two.
	//One part is yaml and the other is markdown
	inputSplit := bytes.Split(input, []byte("\n\n\n\n"))

	//Convert markdown into html
	html := blackfriday.MarkdownCommon(inputSplit[1])

	//Create new post type webPost
	var post webPost
	//Send the first []byte to be parsed by yaml.Unmarshal
	err := post.Parse(inputSplit[0])
	if err != nil {
		panic("failed parsing yaml")
	}

	//Assign html to post.Body and serve to the client
	post.Body = template.HTML(html)
	templates.ExecuteTemplate(res, "Post", post)
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

	//Use this function to handle Posts
	http.HandleFunc("/post/", handlePost)

	//Use the function handleHomePage for any requests to "/"
	http.HandleFunc("/", handleHomePage)

	//Start the server on Port 8080
	http.ListenAndServe(":8080", nil)

	//Or https (requires a cert.pem and a key.pem file)
	//http://www.kaihag.com/https-and-go/
	//http.ListenAndServeTLS(":8080", "cert.pem", "key.pem", nil)
}
