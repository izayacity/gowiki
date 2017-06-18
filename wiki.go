package main

import (
	"fmt"
	"io/ioutil"
)

type Page struct {
	Title string
	Body []byte		// a byte slice
}

func main() {
	p1 := &Page{Title: "TestPage", Body: []byte ("This is a sample Page.")}
	p1.save()
	p2, _ := loadPage("TestPage")
	fmt.Println(string(p2.Body))
}

/*	Function save: This is a method named save that takes as its receiver p, a pointer to Page .
	It takes no parameters, and returns a value of type error.
	perm 0600: the file should be created with read-write permissions for the current user only.
*/
func (p *Page) save () error {
	filename := p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600)
}

/*	function loadPage constructs the file name from the title parameter, reads the file's contents into a new variable body,
	and returns a pointer to a Page literal constructed with the proper title and body values.
*/
func loadPage (title string) (*Page, error) {
	filename := title + ".txt"
	body, err := ioutil.ReadFile (filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

