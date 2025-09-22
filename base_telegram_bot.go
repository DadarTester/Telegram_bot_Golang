package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func readTokenFromFile(filename string) (string, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(data)), nil
}

func main() {
	var botToken string
	var err error

	// Получаем токен бота
	tokenPtr := flag.String("token", "", "Telegram Bot Token")
	flag.Parse()

	if *tokenPtr != "" {
		botToken = *tokenPtr
	} else {
		botToken = os.Getenv("TELEGRAM_BOT_TOKEN")

		if botToken == "" {
			botToken, err = readTokenFromFile(".token")
			if err != nil {
				log.Fatal("Токен бота не найден. Укажите его через -token, переменную окружения TELEGRAM_BOT_TOKEN или файл .token")
			}
		}
	}

	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Panic("Ошибка создания бота: ", err)
	}

	bot.Debug = true
	log.Printf("Авторизован как %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		log.Printf("[%s] Получено сообщение", update.Message.From.UserName)

		switch {
		// Текстовые сообщения
		case update.Message.Text != "":
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
			msg.ReplyToMessageID = update.Message.MessageID
			_, err := bot.Send(msg)
			if err != nil {
				log.Printf("Ошибка отправки сообщения: %v", err)
			}

		// Стикеры
		case update.Message.Sticker != nil:
			// Пересылаем стикер
			sticker := tgbotapi.NewSticker(update.Message.Chat.ID, update.Message.Sticker.FileID)
			sticker.ReplyToMessageID = update.Message.MessageID
			_, err := bot.Send(sticker)
			if err != nil {
				log.Printf("Ошибка отправки стикера: %v", err)
			}

		// Фото
		case len(update.Message.Photo) > 0:
			photo := update.Message.Photo[len(update.Message.Photo)-1]
			// Пересылаем фото
			photoMsg := tgbotapi.NewPhoto(update.Message.Chat.ID, tgbotapi.FileID(photo.FileID))
			photoMsg.ReplyToMessageID = update.Message.MessageID
			photoMsg.Caption = update.Message.Caption

			_, err := bot.Send(photoMsg)
			if err != nil {
				log.Printf("Ошибка отправки фото: %v", err)
			}

		// Файлы
		case update.Message.Document != nil:
			// Пересылаем файлы
			doc := tgbotapi.NewDocument(update.Message.Chat.ID, tgbotapi.FileID(update.Message.Document.FileID))
			doc.ReplyToMessageID = update.Message.MessageID
			doc.Caption = update.Message.Caption

			_, err := bot.Send(doc)
			if err != nil {
				log.Printf("Ошибка отправки документа: %v", err)
			}

		//Гифки
		case update.Message.Animation != nil:
			// Пересылаем Гифку
			animation := tgbotapi.NewAnimation(update.Message.Chat.ID, tgbotapi.FileID(update.Message.Animation.FileID))
			animation.ReplyToMessageID = update.Message.MessageID
			animation.Caption = update.Message.Caption

			_, err := bot.Send(animation)
			if err != nil {
				log.Printf("Ошибка отправки анимации: %v", err)
			}

		// Голосовые сообщения
		case update.Message.Voice != nil:
			// Пересылаем голосовое сообщение
			voice := tgbotapi.NewVoice(update.Message.Chat.ID, tgbotapi.FileID(update.Message.Voice.FileID))
			voice.ReplyToMessageID = update.Message.MessageID
			voice.Caption = update.Message.Caption

			_, err := bot.Send(voice)
			if err != nil {
				log.Printf("Ошибка отправки голосового сообщения: %v", err)
			}

		// Видео сообщения
		case update.Message.Video != nil:
			// Пересылаем видео
			video := tgbotapi.NewVideo(update.Message.Chat.ID, tgbotapi.FileID(update.Message.Video.FileID))
			video.ReplyToMessageID = update.Message.MessageID
			video.Caption = update.Message.Caption

			_, err := bot.Send(video)
			if err != nil {
				log.Printf("Ошибка отправки видео: %v", err)
			}

		// Аудио файлы
		case update.Message.Audio != nil:
			// Пересылаем аудио
			audio := tgbotapi.NewAudio(update.Message.Chat.ID, tgbotapi.FileID(update.Message.Audio.FileID))
			audio.ReplyToMessageID = update.Message.MessageID
			audio.Caption = update.Message.Caption

			_, err := bot.Send(audio)
			if err != nil {
				log.Printf("Ошибка отправки аудио: %v", err)
			}

		// Если тип сообщения отличается от описанных
		default:
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Извините, я пока не умею работать с этим типом сообщений")
			msg.ReplyToMessageID = update.Message.MessageID
			_, err := bot.Send(msg)
			if err != nil {
				log.Printf("Ошибка отправки сообщения: %v", err)
			}
		}
	}
}
