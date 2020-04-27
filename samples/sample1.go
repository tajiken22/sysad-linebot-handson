package main

// 必要なライブラリを使う宣言
import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/line/line-bot-sdk-go/linebot"
)

const verifyToken = "00000000000000000000000000000000"

// 最初に実行される関数
func main() {
	// lineのAPIを利用する設定
	bot, err := linebot.New(
		os.Getenv("CHANNEL_SECRET"),
		os.Getenv("CHANNEL_ACCESS_TOKEN"),
	)
	if err != nil {
		log.Fatal(err)
	}

	// LINEサーバからのリクエストを受け取ったときの処理
	http.HandleFunc("/callback", func(w http.ResponseWriter, req *http.Request) {
		fmt.Print("Accessed\n")

		// リクエストを扱いやすい形に変換する
		events, err := bot.ParseRequest(req)
		// 変換に失敗したとき
		if err != nil {
			fmt.Println("ParseRequest error:", err)
			if err == linebot.ErrInvalidSignature {
				w.WriteHeader(400)
			} else {
				w.WriteHeader(500)
			}
			return
		}

		// LINEサーバから来たメッセージによってやる処理を変える
		for _, event := range events {
			// LINEサーバのverify時は何もしない
			if event.ReplyToken == verifyToken {
				return
			}

			// メッセージが来たとき
			if event.Type == linebot.EventTypeMessage {
				// 返信を生成する
				replyMessage := getReplyMessage(event)
				// 生成した返信を送信する
				if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(replyMessage)).Do(); err != nil {
					log.Print(err)
				}
			}
		}
	})

	// LINEサーバからのリクエストを受け取る
	if err := http.ListenAndServe(":"+os.Getenv("PORT"), nil); err != nil {
		log.Fatal(err)
	}
}

const helpMessage = `使い方
テキストメッセージ:
	やまびこを返すよ！
スタンプ:
	スタンプの情報を答えるよ！
それ以外:
	それ以外にはまだ対応してないよ！ごめんね...`

// 返信を生成する
func getReplyMessage(event *linebot.Event) (replyMessage string) {
	// 来たメッセージの種類によって分岐する
	switch message := event.Message.(type) {
	// テキストメッセージが来たとき
	case *linebot.TextMessage:
		return message.Text

	// スタンプが来たとき
	case *linebot.StickerMessage:
		replyMessage := fmt.Sprintf("sticker id is %s, stickerResourceType is %s", message.StickerID, message.StickerResourceType)
		return replyMessage

	// どっちでもないとき
	default:
		return helpMessage
	}
}
