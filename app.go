package main

import (
	"bytes"
	"github.com/russross/blackfriday"
	"gopkg.in/yaml.v2"
	"html/template"
	"io/ioutil"
	"net/http"
	"strings"
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
	Title string
	Body  template.HTML
}

//Grabs yaml from the markdown and inserts into a type webPage (Title, Author, etc..)
func (page *webPage) Parse(markdownYAML []byte) error {

	//Inserts yaml into type webPage and returns nil if no error
	return yaml.Unmarshal(markdownYAML, page)
}

//Parse markdown and serve to client from templates
func handlePost(res http.ResponseWriter, req *http.Request) {

	//Grab the file name from path
	url := req.URL.Path[len("/post/"):]

	//Read markdown requested by client
	input, err := ioutil.ReadFile("posts/" + url + ".md")
	//404 if file not found
	if err != nil {
		http.NotFound(res, req)
		return
	}

	//Split the []byte input into two.
	//One part is yaml and the other is markdown
	inputSplit := bytes.Split(input, []byte("\n\n\n\n"))

	//Convert markdown into html
	html := blackfriday.MarkdownCommon(inputSplit[1])

	//Create new post type webPost
	var post webPost
	//Send the first []byte to be parsed by yaml.Unmarshal
	err = post.Parse(inputSplit[0])
	if err != nil {
		panic("failed parsing yaml")
	}

	//Assign html to post.Body and serve to the client
	post.Body = template.HTML(html)
	templates.ExecuteTemplate(res, "Post", post)
}

//Handles http requests for url path "/"
func handleWebPage(res http.ResponseWriter, req *http.Request) {

	//Grab url path
	url := strings.ToLower(req.URL.Path[1:])

	//Set the url manually for index
	if url == "" || url == "index.html" {
		url = "index"
	}

	//Try to read in the markdown file based on the path
	input, err := ioutil.ReadFile(url + ".md")
	//404 if file not found
	if err != nil {
		http.NotFound(res, req)
		return
	}

	//Split the input into two [][]byte to extract yaml/markdown seperately
	//Your markdown files should contain 4 new lines between the yaml and markdown
	inputSplit := bytes.Split(input, []byte("\n\n\n\n"))

	//Convert markdown into html
	html := blackfriday.MarkdownCommon(inputSplit[1])

	//Creat a type webPage
	var page webPage
	//Insert yaml into webPage variables
	err = page.Parse(inputSplit[0])
	if err != nil {
		panic("failed parsing yaml")
	}

	//Assign html to page.Body and serve to the client
	page.Body = template.HTML(html)
	templates.ExecuteTemplate(res, "Page", page)
}

func main() {

	//serves static css/js ..etc files
	http.HandleFunc("/public/", func(res http.ResponseWriter, req *http.Request) {
		http.ServeFile(res, req, req.URL.Path[1:])
	})

	//Use this function to handle Posts (/post/sample)
	http.HandleFunc("/post/", handlePost)

	//Routes for type webPage (/, /about)
	http.HandleFunc("/", handleWebPage)

	//Start the server on Port 8080
	http.ListenAndServe(":8080", nil)

	//Or https (requires a cert.pem and a key.pem file)
	//http://www.kaihag.com/https-and-go/
	//http.ListenAndServeTLS(":8080", "cert.pem", "key.pem", nil)
}
