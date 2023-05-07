package main

import (
	"fmt"
	"log"
	"tgbot/bot"
	"time"
)

func main() {
	bot := bot.New()

	_, err := bot.GetMe()
	if err != nil {
		fmt.Println("Cannot ping Telegram", err)
		log.Fatal(err)
	}

	for {
		time.Sleep(time.Second * 1)

		updates, err := bot.GetUpdates()
		if err != nil {
			fmt.Println("Cannot get updates", err)
		}
		if len(updates.Result) == 0 {
			continue
		}
		for _, upd := range updates.Result {
			fmt.Println(upd)

			_, err = bot.SendMessage(map[string]interface{}{
				"chat_id": upd.Message.Chat.Id,
				"text":    "Привет! Ты написал: " + upd.Message.Text,
				"reply_markup": map[string]interface{}{
					"inline_keyboard": [][]interface{}{
						[]interface{}{
							map[string]interface{}{
								"text": "Кнопарь",
								"url":  "https://google.com",
							},
							map[string]interface{}{
								"text": "Кнопарь2",
								"url":  "https://google.com",
							},
						},
					},
				},
			})

			if err != nil {
				fmt.Errorf("Cannot send message: %s", err)
			}
		}
	}
}
