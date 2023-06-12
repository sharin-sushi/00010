package utility

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var Db *sql.DB

func init() {
	user := os.Getenv("MYSQL_USER")
	// if user == "" {
	// 	user = "namae" // デフォルトの値
	// }
	pw := os.Getenv("MYSQL_PASSWORD")
	// if pw == "" {
	// 	pw = "00000" // デフォルトの値
	// }
	db_name := os.Getenv("MYSQL_DATABASE")
	// if db_name == "" {
	// 	db_name = "namae" // デフォルトの値
	// }
	//以上、コメントアウト部は環境変数をプログラム内で指定したい場合の記述方法
	//ハードコーティングなのでやめましょう
	//(テスト時ならokだろうけど忘れて公開されたら
	// コメントアウトしてあっても終わりだからやめた方が良さげ？)

	var path string = fmt.Sprintf("%s:%s@tcp(db:3306)/%s?charset=utf8&parseTime=true", user, pw, db_name)
	var err error
	if Db, err = sql.Open("mysql", path); err != nil {
		fmt.Printf("database.goのinitでエラー発生:%s", err)
		// log.Fatal("Db open error:", err.Error())
	}
	checkConnect(5)

}

func checkConnect(count uint) {
	var err error
	if err = Db.Ping(); err != nil {
		time.Sleep(time.Second * 2)
		count--
		fmt.Printf("retry... count:%v\n", count)
		if count > 0 {
			checkConnect(count)
		} else {
			fmt.Println("Connection retries exhausted")
			return
		}
	} else {
		fmt.Println("db connected!!")
	}
}
