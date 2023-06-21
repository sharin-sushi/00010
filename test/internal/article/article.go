package article

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/sharin-sushi/0010/test/internal/utility"
)

type karaokelist struct {
	unique_id int
	movie     string
	url       string
	singStart string
	song      string
}

const tableklist = "karaokelist" //Mysql側のtable名

var tmpl *template.Template

func init() {
	funcMap := template.FuncMap{
		"nl2br": func(text string) template.HTML {
			return template.HTML(strings.Replace(template.HTMLEscapeString(text), "\n", "<br />", -1))
		},
	}

	tmpl, _ = template.New("article").Funcs(funcMap).ParseGlob("web/template/*")
}

//↓gormを使って省略する　selected.Nextとかダルい
func Index(w http.ResponseWriter, r *http.Request) {
	selected, err := utility.Db.Query("SELECT * FROM " + tableklist)
	fmt.Printf("selected=%v\n, err=%v\n", selected, err) //確認用
	if err != nil {
		panic(err.Error())
	}

	data := []karaokelist{} //懐かしのスライス
	for selected.Next() {
		kList := karaokelist{}
		fmt.Printf("Scan前")
		err = selected.Scan(&kList.unique_id, &kList.movie, &kList.url, &kList.singStart, &kList.song)
		//idとか、頭文字が大文字(型から変える)でないとダメかも
		fmt.Printf("Scan後")
		//errにレコードをぶち込んでる
		//&はポインターみたいな
		if err != nil {
			panic(err.Error())
		}
		data = append(data, kList)
	}
	selected.Close()
	if err := tmpl.ExecuteTemplate(w, "index.html", data); err != nil {
		//↑ここまででエラー発生、because ↓これ処理されてない
		fmt.Printf("%s", err)
		log.Fatal(err)
	}
}

// func Show(w http.ResponseWriter, r *http.Request) {
// 	id := r.URL.Query().Get("id")
// 	selected, err := utility.Db.Query("SELECT * FROM KaraokeListe WHERE id=?", id)
// 	if err != nil {
// 		panic(err.Error())
// 	}
// 	var kList karaokeList
// 	for selected.Next() {
// 		err = selected.Scan(&kList.id, &kList.song, &kList.movie, &kList.url, &kList.singStart)
// 		if err != nil {
// 			panic(err.Error())
// 		}
// 	}
// 	selected.Close()
// 	tmpl.ExecuteTemplate(w, "show.html", kList)
// }

// func Create(w http.ResponseWriter, r *http.Request) {
// 	if r.Method == "GET" {
// 		tmpl.ExecuteTemplate(w, "create.html", nil)
// 	} else if r.Method == "POST" {
// 		title := r.FormValue("title")
// 		body := r.FormValue("body")
// 		insert, err := utility.Db.Prepare("INSERT INTO KaraokeList(title, body) VALUES(?,?)")
// 		if err != nil {
// 			panic(err.Error())
// 		}
// 		insert.Exec(title, body)
// 		http.Redirect(w, r, "/", 301)
// 	}
// }

// func Edit(w http.ResponseWriter, r *http.Request) {
// 	if r.Method == "GET" {
// 		id := r.URL.Query().Get("id")
// 		selected, err := utility.Db.Query("SELECT * FROM KaraokeList WHERE id=?", id)
// 		if err != nil {
// 			panic(err.Error())
// 		}
//		↓変更する
// 		article := karaokeList{}
// 		for selected.Next() {
// 			err = selected.Scan(&karaokeList.id, &karaokeList.song, &karaokeList.movie, &karaokeList.url, &karaokeList.singStart)
// 			if err != nil {
// 				panic(err.Error())
// 			}
// 		}
// 		selected.Close()
// 		tmpl.ExecuteTemplate(w, "edit.html", article)
// 	} else if r.Method == "POST" {
// 		title := r.FormValue("title")
// 		body := r.FormValue("body")
// 		id := r.FormValue("id")
// 		insert, err := utility.Db.Prepare("UPDATE KaraokeList SET title=?, body=? WHERE id=?")
// 		if err != nil {
// 			panic(err.Error())
// 		}
// 		insert.Exec(title, body, id)
// 		http.Redirect(w, r, "/", 301)
// 	}
// }

// func Delete(w http.ResponseWriter, r *http.Request) {
// 	if r.Method == "GET" {
// 		id := r.URL.Query().Get("id")
// 		selected, err := utility.Db.Query("SELECT * FROM KaraokeList WHERE id=?", id)
// 		if err != nil {
// 			panic(err.Error())
// 		}
// 		article := karaokeList{}
// 		for selected.Next() {
// 			err = selected.Scan(&karaokeList.id, &karaokeList.song, &karaokeList.movie, &karaokeList.url, &karaokeList.singStart)
// 			if err != nil {
// 				panic(err.Error())
// 			}
// 		}
// 		selected.Close()
// 		tmpl.ExecuteTemplate(w, "delete.html", article)
// 	} else if r.Method == "POST" {
// 		id := r.FormValue("id")
// 		insert, err := utility.Db.Prepare("DELETE FROM KaraokeList WHERE id=?")
// 		if err != nil {
// 			panic(err.Error())
// 		}
// 		insert.Exec(id)
// 		http.Redirect(w, r, "/", 301)
// 	}
// }
