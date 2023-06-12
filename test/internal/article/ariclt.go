package article

import (
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/sharin-sushi/0010/test/internal/utility"
)

type Article struct {
	Id    int
	Title string
	Body  string
}

var tmpl *template.Template

func init() {
	funcMap := template.FuncMap{
		"nl2br": func(text string) template.HTML {
			return template.HTML(strings.Replace(template.HTMLEscapeString(text), "\n", "<br />", -1))
		},
	}

	tmpl, _ = template.New("article").Funcs(funcMap).ParseGlob("web/template/*")
}

func Index(w http.ResponseWriter, r *http.Request) {
	selected, err := utility.Db.Query("SELECT * FROM Article")
	if err != nil {
		panic(err.Error())
	}
	data := []Article{}
	for selected.Next() {
		article := Article{}
		err = selected.Scan(&article.Id, &article.Title, &article.Body)
		if err != nil {
			panic(err.Error())
		}
		data = append(data, article)
	}
	selected.Close()
	if err := tmpl.ExecuteTemplate(w, "index.html", data); err != nil {
		log.Fatal(err)
	}
}

func Show(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	selected, err := utility.Db.Query("SELECT * FROM Article WHERE id=?", id)
	if err != nil {
		panic(err.Error())
	}
	article := Article{}
	for selected.Next() {
		err = selected.Scan(&article.Id, &article.Title, &article.Body)
		if err != nil {
			panic(err.Error())
		}
	}
	selected.Close()
	tmpl.ExecuteTemplate(w, "show.html", article)
}

func Create(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		tmpl.ExecuteTemplate(w, "create.html", nil)
	} else if r.Method == "POST" {
		title := r.FormValue("title")
		body := r.FormValue("body")
		insert, err := utility.Db.Prepare("INSERT INTO Article(title, body) VALUES(?,?)")
		if err != nil {
			panic(err.Error())
		}
		insert.Exec(title, body)
		http.Redirect(w, r, "/", 301)
	}
}

func Edit(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		id := r.URL.Query().Get("id")
		selected, err := utility.Db.Query("SELECT * FROM Article WHERE id=?", id)
		if err != nil {
			panic(err.Error())
		}
		article := Article{}
		for selected.Next() {
			err = selected.Scan(&article.Id, &article.Title, &article.Body)
			if err != nil {
				panic(err.Error())
			}
		}
		selected.Close()
		tmpl.ExecuteTemplate(w, "edit.html", article)
	} else if r.Method == "POST" {
		title := r.FormValue("title")
		body := r.FormValue("body")
		id := r.FormValue("id")
		insert, err := utility.Db.Prepare("UPDATE Article SET title=?, body=? WHERE id=?")
		if err != nil {
			panic(err.Error())
		}
		insert.Exec(title, body, id)
		http.Redirect(w, r, "/", 301)
	}
}

func Delete(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		id := r.URL.Query().Get("id")
		selected, err := utility.Db.Query("SELECT * FROM Article WHERE id=?", id)
		if err != nil {
			panic(err.Error())
		}
		article := Article{}
		for selected.Next() {
			err = selected.Scan(&article.Id, &article.Title, &article.Body)
			if err != nil {
				panic(err.Error())
			}
		}
		selected.Close()
		tmpl.ExecuteTemplate(w, "delete.html", article)
	} else if r.Method == "POST" {
		id := r.FormValue("id")
		insert, err := utility.Db.Prepare("DELETE FROM Article WHERE id=?")
		if err != nil {
			panic(err.Error())
		}
		insert.Exec(id)
		http.Redirect(w, r, "/", 301)
	}
}
