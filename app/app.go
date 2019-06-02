package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"text/template"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"google.golang.org/appengine"
)

var db *sql.DB

func main() {
	db = DB()

	http.HandleFunc("/", indexHandle)
	http.HandleFunc("/send", sendHandle)
	http.HandleFunc("/messages", messagesHandle)
	appengine.Main()
}

/* Controllers */
// TitleData : タイトル用の構造体
type TitleData struct {
	Title string
}

// MessageData : メッセージ用の構造体
type MessageData struct {
	Name    string
	Content string
}

// ResultData : メッセージ + タイトルの構造体
type ResultData struct {
	Title   string
	Message MessageData
}

// MessagesData : 複数メッセージ + タイトルの構造体
type MessagesData struct {
	Title    string
	Messages []MessageData
}

func indexHandle(w http.ResponseWriter, r *http.Request) {
	// index.htmlをパース
	tmpl := template.Must(template.ParseFiles("views/index.html"))

	// タイトル挿入
	title := "ボトルメール"

	// テンプレートを実行して出力
	data := TitleData{title}

	// 実行
	if err := tmpl.ExecuteTemplate(w, "base", data); err != nil {
		log.Fatal(err)
		fmt.Fprintln(w, "Error has occured")
	}
}

func sendHandle(w http.ResponseWriter, r *http.Request) {
	// Postかチェック
	var name, content string
	if r.Method == "POST" {
		name = r.FormValue("name")
		content = r.FormValue("content")
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}

	// 自分以外の全メッセージ取得
	//messages := GetExceptMeDB(name)
	messages := GetAllDB()

	// 自分のメッセージを格納
	InsertDB(name, content)

	// 他人のメッセージをランダムに選ぶ
	rand.Seed(time.Now().UnixNano())
	length := len(messages)
	num := 0
	if length > 0 {
		num = rand.Intn(length)
	}

	// 構造体にデータを整形
	title := "ボトルメール"
	var data ResultData

	//	はじめての人だったらば, Welcomeメッセージを表示
	if length == 0 {
		data = ResultData{title, MessageData{"you are first person", "welcome!"}}
	} else {
		data = ResultData{title, messages[num]}
	}

	tmpl := template.Must(template.ParseFiles("views/result.html"))
	if err := tmpl.ExecuteTemplate(w, "result", data); err != nil {
		log.Fatal(err)
		fmt.Fprintln(w, "Error has occured")
	}
}

func messagesHandle(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("views/messages.html"))

	// 全メッセージを取得
	var posts []MessageData
	posts = GetAllDB()

	// 構造体にデータを整形して入れる
	title := "ボトルメール"
	data := MessagesData{title, posts}
	if err := tmpl.ExecuteTemplate(w, "messages", data); err != nil {
		log.Fatal(err)
		fmt.Fprintln(w, "Error has occured")
	}
}

/* Models */
var (
	connectionName = mustGetenv("CLOUDSQL_CONNECTION_NAME")
	user           = mustGetenv("CLOUDSQL_USER")
	password       = os.Getenv("CLOUDSQL_PASSWORD") // NOTE: password may be empty
	host_address   = mustGetenv("CLOUDSQL_HOST_ADDRESS")
	database       = mustGetenv("CLOUDSQL_DBNAME")
	socket         = os.Getenv("CLOUDSQL_SOCKET_PREFIX")
)

func mustGetenv(k string) string {
	v := os.Getenv(k)
	if v == "" {
		log.Panicf("%s environment variable not set.", k)
	}
	return v
}

// DB : Databaseに接続
func DB() *sql.DB {
	if socket == "" {
		socket = "/cloudsql"
	}
	var dbURI string
	if appengine.IsDevAppServer() {
		dbURI = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, password, "127.0.0.1", "3306", database)
	} else {
		dbURI = fmt.Sprintf("%s:%s@unix(%s/%s)/%s", user, password, socket, connectionName, database)
	}
	conn, err := sql.Open("mysql", dbURI)

	if err != nil {
		panic(fmt.Sprintf("DB: %v", err))
	}

	return conn
}

// GetExceptMeDB : me文字列以外の名前を持つ全メッセージを返却, 空の場合は空配列を返却
func GetExceptMeDB(me string) []MessageData {
	var messages []MessageData
	rows, err := db.Query("SELECT * FROM content where name != ?", me)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var name, content string
		/*
			if err := rows.Scan(&name, &content); err != nil {
				log.Fatal(err)
			}
		*/
		rows.Scan(&name, &content)
		messages = append(messages, MessageData{name, content})
	}
	return messages
}

// GetAllDB : 全メッセージを返却
func GetAllDB() []MessageData {
	var messages []MessageData
	rows, _ := db.Query("SELECT * FROM content")
	/*
		if err != nil {
			log.Fatal(err)
		}
	*/
	defer rows.Close()
	for rows.Next() {
		var name, content string
		/*
			if err := rows.Scan(&name, &content); err != nil {
				log.Fatal(err)
			}
		*/
		rows.Scan(&name, &content)
		messages = append(messages, MessageData{name, content})
	}
	return messages
}

// InsertDB : メッセージ(名前:内容)を格納
func InsertDB(name string, content string) {
	var err error
	_, err = db.Exec("INSERT INTO content VALUES (?, ?)", name, content)
	if err != nil {
		log.Fatal(err)
	}
}
