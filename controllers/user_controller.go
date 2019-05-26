package controllers

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"bottlemail/models"

	"github.com/gin-gonic/gin"
)

//Index : メッセージ入力画面, ランディング
func Index(ctx *gin.Context) {
	models.DBInit()
	ctx.HTML(200, "index.html", gin.H{})
}

//Send : ボトルメールを送り, 他人のボトルメールを受け取る
func Send(ctx *gin.Context) {
	name := ctx.PostForm("name")
	content := ctx.PostForm("content")

	//データベースに自分のメッセージを挿入
	models.DBInsert(name, content)

	//データベースから全メッセージ取得
	//Posts := DBGetAll()

	//データベースから自分以外のメッセージ取得
	Posts := models.DBGetExceptMe(name)

	//他のメッセージからランダムに表示するものを選択
	rand.Seed(time.Now().UnixNano())
	fmt.Println(len(Posts))

	length := len(Posts)
	if length > 0 {
		num := rand.Intn(length)

		ctx.HTML(http.StatusOK, "result.html", gin.H{
			"name":    Posts[num].Name,
			"content": Posts[num].Content,
		})
	} else {
		ctx.HTML(http.StatusOK, "result.html", gin.H{
			"name":    "You are first person",
			"content": "Enjoy bottlemailing!",
		})
	}
	//ctx.Redirect(302, "/result.html")
}

//Messages : メッセージ入力画面, ランディング
func Messages(ctx *gin.Context) {
	models.DBInit()

	Posts := models.DBGetAll()

	ctx.HTML(200, "messages.html", gin.H{
		"posts": Posts,
	})
}
