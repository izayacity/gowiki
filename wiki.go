/**
 * Created by Francis Yang. User: izayacity, Email: izayacity@gmail.com, Date: 2016/06/17
*/

package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

type Page struct {
	Title string
	Body []byte		// a byte slice
}

func main() {
	http.HandleFunc ("/view/", viewHandler)
	http.ListenAndServe (":8080", nil)
	//p1 := &Page{Title: "TestPage", Body: []byte ("This is a sample Page.")}
	//p1.save()
	//p2, _ := loadPage("TestPage")
	//fmt.Println(string(p2.Body))
}

// Function save: This is a method named save that takes as its receiver p, a pointer to Page .
// It takes no parameters, and returns a value of type error.
// perm 0600: the file should be created with read-write permissions for the current user only.
func (p *Page) save () error {
	filename := p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600)
}

// Function loadPage: constructs the file name from the title parameter, reads the file's contents into a new variable body,
// and returns a pointer to a Page literal constructed with the proper title and body values.
func loadPage (title string) (*Page, error) {
	filename := title + ".txt"
	body, err := ioutil.ReadFile (filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

// Function viewHandler: extracts the page title from r.URL.Path, loads the page data,
// formats the page with a string of simple HTML, and writes it to http.ResponseWriter
func viewHandler (w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path [len ("/view/"):]
	p, _ := loadPage(title)
	fmt.Fprintf(w, "<h1>%s</h1><div>%s</div>", p.Title, p.Body)
}