package main

//参考　https://www.yoheim.net/blog.php?q=20170403
// https://zenn.dev/ajapa/articles/03dcf8fd12d086

import (
	"log"
	"net/http"

	"github.com/sharin-sushi/0010/test/internal/article"
)

//importするときは”ディレクトリ".関数名
//ファイル名じゃないから注意

func main() {
	// http.HandleFunc("/secret", utility.Oimo)
	// http.HandleFunc("/", handler)
	http.HandleFunc("/", article.Index)
	http.HandleFunc("/show", article.Show)
	http.HandleFunc("/create", article.Create)
	http.HandleFunc("/edit", article.Edit)
	http.HandleFunc("/delete", article.Delete)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

//fmt.Fprint()の第２引数はstring型である必要があり、引用元の関数を
// func Oimo() stringとする
