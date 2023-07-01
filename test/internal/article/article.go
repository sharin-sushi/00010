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
	Unique_id int    `db:"unique_id"`
	Movie     string `db:"movie"`
	Url       string `db:"url"`
	SingStart string `db:"singStart"`
	Song      string `db:"song"`
}

// db:"カラム名"という形式のタグを使用して、
// 構造体のフィールドとデータベースのカラムを対応付け
//https://leben.mobi/go/post-370/practice/#Query_8211_SELECT

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

//↓gormを使って省略可　selected.Nextとか
func Index(w http.ResponseWriter, r *http.Request) {
	selected, err := utility.Db.Query("SELECT * FROM karaokelist")
	fmt.Printf("selected=%v\n, err=%v\n", selected, err) //確認用
	if err != nil {
		panic(err.Error())
	}
	//変数(構造体？)selectedにkaraokelistテーブルの全てが代入された

	tabeleData := []karaokelist{} //懐かしのスライス
	for selected.Next() {
		//selected.Nect データベースのクエリ結果セット内で次の行が存在するかどうかを確認する
		//forとセットで、次の行が有る限りループする
		//bool型だからっぽい　良そうだけど、次の行が無くなるとfalseになる

		kList := karaokelist{}
		fmt.Printf("Scan前"+"%v\n", kList)
		err = selected.Scan(&kList.Unique_id, &kList.Movie, &kList.Url, &kList.SingStart, &kList.Song)
		//カラムが頭文字が大文字ダメ説がある
		fmt.Printf("Scan後"+"%v\n", kList)
		//errにレコードをぶち込んでる=関数結果そのもの？　＆kListにぶちこんでるんだよな？
		//&はポインターみたいな
		if err != nil {
			panic(err.Error())
		}
		tabeleData = append(tabeleData, kList)
		//tableDataスライスにkList(要素)を組み込んでる
	}

	selected.Close()
	fmt.Printf("selected.Closeの1行目")
	// ↑ここまでは処理してる

	tmpl, err := template.ParseFiles("index.html")
	if err != nil {
		fmt.Printf("Failed to parse template: %s", err)
		log.Fatal(err)
	}

	if err := tmpl.Execute(w, tabeleData); err != nil {
		fmt.Printf("Failed to execute template: %s", err)
		log.Fatal(err)
	}
	// ーーー↑↓どっちかーーーー
	// エラーハンドリングが違う。個々で見るかまとめて見るか。

	// でも↓だと二ルポでる
	// if err := tmpl.ExecuteTemplate(w, "index.html", tabeleData); err != nil {
	// 	//"index.html”テンプレートはプロジェクト内ならどのディレクトリでも良い

	// 	//↑ここまででエラー発生、because ↓これ処理されてないつまり nil == err
	// 	fmt.Printf("74行目")
	// 	log.Fatal(err)
	// }

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
