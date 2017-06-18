/**
 * Created by Francis Yang. User: izayacity, Email: izayacity@gmail.com, Date: 2016/06/17
*/
package main
import (
	"io/ioutil"
	"net/http"
	"html/template"
	"regexp"
	"errors"
)

type Page struct {
	Title string
	Body []byte		// a byte slice
}

// global variable that renders HTML templates for handlers for caching, called in renderTemplate function
var templates = template.Must(template.ParseFiles("edit.html", "view.html"))
// global variable to store our validation expression and validate the title to prevent a user can supply an arbitrary path to be read/written on the server.
var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")

func main() {
	http.HandleFunc ("/view/", makeHandler(viewHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))
	http.ListenAndServe (":8080", nil)
}

// Function getTitle: security helper that uses the validPath expression to validate path and extract the page title.
func getTitle (w http.ResponseWriter, r *http.Request) (string, error) {
	m := validPath.FindStringSubmatch(r.URL.Path)
	if m == nil {
		http.NotFound(w, r)
		return "", errors.New("Invalid Page Title")
	}
	// The title is the second subexpression.
	return m[2], nil
}

//  a wrapper function that takes a function of the above type, and returns a closure function of type http.HandlerFunc
// handles security that uses the validPath expression to validate path and extract the page title, like function getTitle
// variable fn is enclosed by the closure to make the handler. The variable fn will be one of our save, edit, or view handlers.
func makeHandler (fn func (http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		// Here we will extract the page title from the Request
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		// and call the provided handler 'fn'
		fn (w, r, m[2])
	}
}

// Function save: This is a method named save that takes as its receiver p, a pointer to Page .
// It takes no parameters, and returns a value of type error.
// Example:
// p1 := &Page{Title: "TestPage", Body: []byte ("This is a sample Page.")}
// p1.save()
func (p *Page) save () error {
	filename := p.Title + ".txt"
	// perm 0600: the file should be created with read-write permissions for the current user only.
	return ioutil.WriteFile(filename, p.Body, 0600)
}

// Function loadPage: constructs the file name from the title parameter, reads the file's contents into a new variable body,
// and returns a pointer to a Page literal constructed with the proper title and body values.
// Example:
// p2, _ := loadPage("TestPage")
// fmt.Println(string(p2.Body))
func loadPage (title string) (*Page, error) {
	filename := title + ".txt"
	body, err := ioutil.ReadFile (filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

// Function renderTemplate: executes the templates global variable
func renderTemplate (w http.ResponseWriter, tmplt string, p *Page) {
	err := templates.ExecuteTemplate(w, tmplt + ".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Function viewHandler: The http.Redirect function adds an HTTP status code of http.StatusFound (302) and a Location header to the HTTP response.
func viewHandler (w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		http.Redirect(w, r, "/edit/" + title, http.StatusFound)
		return
	}
	renderTemplate(w, "view", p)
}

// Function editHandler: extracts the page title from r.URL.Path, loads the page data,
// formats the page with a string of simple HTML, and writes it to http.ResponseWriter
func editHandler (w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	// if it doesn't exist, create an empty Page struct
	if err != nil {
		p = &Page{Title: title}
	}
	renderTemplate(w, "edit", p)
}

// The page title (provided in the URL) and the form's only field, Body, are stored in a new Page.
func saveHandler (w http.ResponseWriter, r *http.Request, title string) {
	body := r.FormValue("body")
	// The value returned by FormValue is of type string. We must convert that value to []byte before it will fit into the Page struct.
	p := &Page {Title: title, Body: []byte (body)}
	// write the data to a file, and the client is redirected to the /view/ page.
	err := p.save()
	if err != nil {
		http.Error (w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/view/" + title, http.StatusFound)
}