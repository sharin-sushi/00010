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

//`db:"unique_id"`は関数外に書いたら意味ないらしい↑意味ない
//  Unique_id int    `db:"unique_id"`
// 	Movie     string `db:"movie"`
// 	Url       string `db:"url"`
// 	SingStart string `db:"singStart"`
// 	Song      string `db:"song"`

// Unique_id, Movie, Url, SingStart, Song
// unique_id, movie, url, singStart, song

// db:"カラム名"という形式のタグを使用して、
// 構造体のフィールドとデータベースのカラムを対応付け
//https://leben.mobi/go/post-370/practice/#Query_8211_SELECT

// 関数utility.Db.Prepare
// https://ota42y.com/blog/2014/10/04/go-mysql/

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
	// fmt.Printf("selected=%v\n, err=%v\n", selected, err) //確認用
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
		// // fmt.Printf("Scan前"+"%v\n", kList)
		err = selected.Scan(&kList.Unique_id, &kList.Movie, &kList.Url, &kList.SingStart, &kList.Song)
		// カラムが頭文字が大文字ダメ説がある
		// // fmt.Printf("Scan後"+"%v\n", kList)
		//errにレコードをぶち込んでる=関数結果そのもの？　＆kListにぶちこんでるんだよな？
		//&はポインターみたいな
		if err != nil {
			panic(err.Error())
		}
		tabeleData = append(tabeleData, kList)
		//tableDataスライスにkList(要素)を組み込んでる…？
	}

	selected.Close()
	// fmt.Printf("selected.Closeの1行目")
	// // ↑ここまでは処理してる

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

func Show(w http.ResponseWriter, r *http.Request) {
	Unique_id := r.URL.Query().Get("Unique_id")
	fmt.Println(r.URL.String())
	fmt.Println(r.URL.Query().Get("Unique_id"))
	// r というhttp.RequestオブジェクトのURLを文字列として表示
	selected, err := utility.Db.Query("SELECT * FROM KaraokeList WHERE Unique_id=?", Unique_id)
	if err != nil {
		panic(err.Error())
	}
	kList := karaokelist{}
	for selected.Next() {
		err = selected.Scan(&kList.Unique_id, &kList.Movie, &kList.Url, &kList.SingStart, &kList.Song)
		if err != nil {
			panic(err.Error())
		}
	}
	fmt.Println(kList)
	selected.Close()
	// tmpl.ExecuteTemplate(w, "show.html", kList) ←元サイトの記述ダメだった
	tmpl, err := template.ParseFiles("show.html")
	if err != nil {
		fmt.Printf("Failed to parse template: %s", err)
		log.Fatal(err)
	}

	if err := tmpl.Execute(w, kList); err != nil {
		fmt.Printf("Failed to execute template: %s", err)
		log.Fatal(err)
	}
}

func Create(w http.ResponseWriter, r *http.Request) {
	kList := karaokelist{}
	if r.Method == "GET" {
		// tmpl.ExecuteTemplate(w, "create.html", nil)
		tmpl, err := template.ParseFiles("create.html")
		if err != nil {
			fmt.Printf("Failed to parse createtemplate: %s", err)
			log.Fatal(err)
		}

		if err := tmpl.Execute(w, kList); err != nil {
			fmt.Printf("Failed to execute createtemplate: %s", err)
			log.Fatal(err)
		}

		fmt.Println("GET受信")
	} else if r.Method == "POST" {
		fmt.Println("POST受信")
		// Unique_id := r.FormValue("Unique_id")
		//""の中の
		Movie := r.FormValue("Movie")
		Url := r.FormValue("url")
		SingStart := r.FormValue("singStart")
		Song := r.FormValue("song")
		fmt.Println("klist:", Movie, Url, SingStart, Song)
		insert, err := utility.Db.Prepare("INSERT INTO KaraokeList(movie, url, singStart, song) VALUES(?,?,?,?)")
		if err != nil {
			panic(err.Error())
		}
		insert.Exec(Movie, Url, SingStart, Song)
		fmt.Println("klist:", insert)

		http.Redirect(w, r, "/", 301)
	}
}

func Edit(w http.ResponseWriter, r *http.Request) {
	// これGETリクエストへのレスポンスって何のために記述してるの？？？↓
	// edit.htmlにはpostしか書いていてない
	// index.htmlのリンクがGET?

	if r.Method == "GET" {
		unique_id := r.URL.Query().Get("Unique_id")
		//urlのuniqueid=の値を取得

		// ↓の手打ち selected, err := utility.Db.Query("SELECT * FROM KaraokeList WHERE unique_id=1")
		selected, err := utility.Db.Query("SELECT * FROM KaraokeList WHERE unique_id=?", unique_id)
		if err != nil {
			panic(err.Error())
		}
		// ↓変更する
		kList := karaokelist{}
		for selected.Next() {
			err = selected.Scan(&kList.Unique_id, &kList.Movie, &kList.Url, &kList.SingStart, &kList.Song)
			if err != nil {
				panic(err.Error())
			}
			fmt.Printf("selectした値: %v", selected)
		}
		selected.Close() //メモリ解放

		tmpl, err := template.ParseFiles("edit.html")
		if err != nil {
			fmt.Printf("Failed to edit'parse template: %s", err)
			log.Fatal(err)
		}

		if err := tmpl.Execute(w, kList); err != nil {
			fmt.Printf("Failed to edit'execute template: %s", err)
			log.Fatal(err)
		}

		// ↑↓どっちか
		// tmpl.ExecuteTemplate(w, "edit.html", kList)
		//"第２引数"内を実行してwに渡し、第３引数オブジェクトを利用してHTMLを生成。

	} else if r.Method == "POST" {

		// if r.Method == "POST" {
		Unique_id := r.FormValue("Unique_id")
		// ↓この１文変更したかったけど、動的にしないと無理か…htmlじゃ無理じゃね？js?
		// Unique_id, err := utility.Db.Query("SELECT * FROM KaraokeList WHERE unique_id=?", Unique_id)
		Movie := r.FormValue("Movie")
		Url := r.FormValue("url")
		SingStart := r.FormValue("singStart")
		Song := r.FormValue("song")
		fmt.Println("klist:", Unique_id, Movie, Url, SingStart, Song)
		database, err := utility.Db.Prepare("UPDATE KaraokeList SET movie=?, url=?, singStart=?, song=? WHERE unique_id=?")
		if err != nil {
			panic(err.Error())
		}
		result, err := database.Exec(Movie, Url, SingStart, Song, Unique_id)
		if err != nil {
			panic(err.Error())

		}
		fmt.Println(result)

		http.Redirect(w, r, "/", 301)
		fmt.Printf("POSTリクエストです")
	} else {
		fmt.Printf("POSTリクエストではありません")
	}
}

func Delete(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		unique_id := r.URL.Query().Get("Unique_id")

		selected, err := utility.Db.Query("SELECT * FROM KaraokeList WHERE unique_id=?", unique_id)
		if err != nil {
			panic(err.Error())
		}
		kList := karaokelist{}
		for selected.Next() {
			err = selected.Scan(&kList.Unique_id, &kList.Movie, &kList.Url, &kList.SingStart, &kList.Song)
			if err != nil {
				panic(err.Error())
			}
		}
		selected.Close()
		tableDate, err := template.ParseFiles("delete.html")
		if err != nil {
			fmt.Printf("Failed to edit'parse template: %s", err)
			log.Fatal(err)
		}

		if err := tableDate.Execute(w, kList); err != nil {
			fmt.Printf("Failed to edit'execute template: %s", err)
			log.Fatal(err)
		}

		// tmpl.ExecuteTemplate(w, "delete.html", kList)
	} else if r.Method == "POST" {
		Unique_id := r.FormValue("Unique_id")
		insert, err := utility.Db.Prepare("DELETE FROM KaraokeList WHERE unique_id=?")
		if err != nil {
			panic(err.Error())
		}
		insert.Exec(Unique_id)
		//検索結果を取得しない場合（create, insert, update, delete）
		http.Redirect(w, r, "/", 301)
	}
}

// SQLのロールバック ROLLBACK;
// なお、コミットcommit; 後は不可能
