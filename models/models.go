package models

import "github.com/jinzhu/gorm"

//Post : ボトルメールの構造体
type Post struct {
	gorm.Model
	Name    string `gorm:"size:255"`
	Content string `gorm:"size:255"`
}

//DBInit : 初期化
func DBInit() {
	//db, err := gorm.Open("sqlite3", "post.sqlite3")
	db, err := gorm.Open("postgres",  "user=tikeda password=root dbname=bottlemail sslmode=disable")
	if err != nil {
		panic("[dbInit] failed to open db")
	}
	db.AutoMigrate(&Post{})
	defer db.Close()
}

//DBInsert : データを挿入
func DBInsert(name string, content string) {
	//db, err := gorm.Open("sqlite3", "post.sqlite3")
        db, err := gorm.Open("postgres",  "user=tikeda password=root dbname=bottlemail sslmode=disable")
	if err != nil {
		panic("[dbInsert] failed to open db")
	}
	db.Create(&Post{Name: name, Content: content})
	defer db.Close()
}

//DBGetAll : 全データを取得
func DBGetAll() []Post {
	//db, err := gorm.Open("sqlite3", "post.sqlite3")
	db, err := gorm.Open("postgres",  "user=tikeda password=root dbname=bottlemail sslmode=disable")
	if err != nil {
		panic("[dbGetAll] failed to open db")
	}
	var posts []Post
	db.Order("created_at desc").Find(&posts)
	db.Close()
	return posts
}

//DBGetExceptMe : 対象以外のデータ取得
func DBGetExceptMe(name string) []Post {
	//db, err := gorm.Open("sqlite3", "post.sqlite3")
        db, err := gorm.Open("postgres",  "user=tikeda password=root dbname=bottlemail sslmode=disable")
	if err != nil {
		panic("[dbGetAll] failed to open db")
	}
	var posts []Post
	db.Not("name = ?", name).First(&posts)
	db.Close()
	return posts
}
