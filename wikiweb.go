
// basic Go webapp tutorial at: https://golang.org/doc/articles/wiki/

package main

import(
    // "fmt"
    "io/ioutil"
	"net/http"
	"html/template"
)

var templates = template.Must(template.ParseFiles("edit.html", "view.html"))

func main() {
	// p1 := &Page{Title: "Test Page", Body: []byte("This is an example Page.")}
	// p1.save()

	// p2, _ := loadPage("Test Page")
	// fmt.Println(string(p2.Body))

	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/edit/", editHandler)
	http.HandleFunc("/save/", saveHandler)
	http.ListenAndServe(":8080", nil)
}

type Page struct {
	Title string
	Body []byte
}

func (p *Page) save() error  {
	filename := p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	body, err := ioutil.ReadFile(filename)

	if err != nil {
		return nil, err
	}

	return &Page{Title: title, Body: body}, nil
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/view/"):]
	p, err := loadPage(title)
	// fmt.Fprintf(w, "<h1>%s</h1><div>%s</div>", p.Title, p.Body)

	if err != nil {
		http.Redirect(w, r, "/edit/" + title, http.StatusFound)
		return
	}

	// t, _ := template.ParseFiles("view.htm")
	// t.Execute(w, p)

	renderTemplate(w, "view", p)
}

func editHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/edit/"):]
	p, err := loadPage(title)

	if err != nil {
		p = &Page{Title: title}
	}

	// t, _ := template.ParseFiles("edit.htm")
	// t.Execute(w, p)

	renderTemplate(w, "edit", p)

	// fmt.Fprintf(w, "<h1>Editing %s</h1>"+
    //     "<form action=\"/save/%s\" method=\"POST\">"+
    //     "<textarea name=\"body\">%s</textarea><br>"+
    //     "<input type=\"submit\" value=\"Save\">"+
    //     "</form>",
    //     p.Title, p.Title, p.Body)
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/save/"):]
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	err := p.save()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/view/" + title, http.StatusFound)
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
    err := templates.ExecuteTemplate(w, tmpl+".html", p)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}