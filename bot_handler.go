package main

import (
	"fmt"
	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func MessageHandler(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	if update.Message == nil {
		return
	}

	switch {
	case update.Message.Text != "":
		HandleTextMessage(bot, update)
	case update.Message.Sticker != nil:
		HandleStickerMessage(bot, update)
	case update.Message.Photo != nil:
		HandlePhotoMessage(bot, update)
	case update.Message.Document != nil:
		HandleDocumentMessage(bot, update)
	case update.Message.Video != nil:
		HandleVideoMessage(bot, update)
	case update.Message.Voice != nil:
		HandleVoiceMessage(bot, update)
	case update.Message.Audio != nil:
		HandleAudioMessage(bot, update)
	case update.Message.VideoNote != nil:
		HandleVideoNoteMessage(bot, update)
	}
}

func HandleTextMessage(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	msg := update.Message
	text := strings.ToLower(msg.Text)

	var response string

	switch {
	case strings.Contains(text, "привет") || strings.Contains(text, "start"):
		response = fmt.Sprintf("👋 Привет, %s!", msg.From.FirstName)
	case strings.Contains(text, "помощь") || strings.Contains(text, "help"):
		response = "🆘 Помощь: я отвечаю на сообщения!"
	default:
		response = fmt.Sprintf("📝 Вы сказали: \"%s\"", msg.Text)
	}

	sendMessage(bot, msg.Chat.ID, response)
}

// обработка фото
func HandlePhotoMessage(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	msg := update.Message
	photo := msg.Photo[len(msg.Photo)-1] // Самая большая версия фото

	response := fmt.Sprintf("📸 Фото получено! Размер: %dx%dpx",
		photo.Width, photo.Height)

	sendMessage(bot, msg.Chat.ID, response)
}

// обработка стикеров
func HandleStickerMessage(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	msg := update.Message
	sticker := msg.Sticker

	response := fmt.Sprintf("🎭 Стикер! Эмодзи: %s", sticker.Emoji)

	sendMessage(bot, msg.Chat.ID, response)
}

// обработка файлов
func HandleDocumentMessage(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	msg := update.Message
	Document := msg.Document
	response := fmt.Sprintf("🎭 Файл!: %s", Document.FileID)

	sendMessage(bot, msg.Chat.ID, response)
}

// обработка видеосообщения
func HandleVideoMessage(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	msg := update.Message
	Video := msg.Video
	response := fmt.Sprintf("🎭 Файл!: %s", Video.FileID)

	sendMessage(bot, msg.Chat.ID, response)
}

// обработка голосового
func HandleVoiceMessage(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	msg := update.Message
	Voice := msg.Voice
	response := fmt.Sprintf("🎭 Файл!: %s", Voice.FileID)

	sendMessage(bot, msg.Chat.ID, response)
}

// обработка аудиосообщения
func HandleAudioMessage(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	msg := update.Message
	Audio := msg.Audio
	response := fmt.Sprintf("🎭 Файл!: %s", Audio.FileID)

	sendMessage(bot, msg.Chat.ID, response)
}

// обработка кружочков
func HandleVideoNoteMessage(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	msg := update.Message
	VideoNote := msg.VideoNote
	response := fmt.Sprintf("🎭 Файл!: %s", VideoNote.FileID)

	sendMessage(bot, msg.Chat.ID, response)
}

// handleUnknownMessage - обработка неизвестных типов сообщений
func HandleUnknownMessage(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	sendMessage(bot, update.Message.Chat.ID, "❓ Неизвестный тип сообщения")
}

// sendMessage - вспомогательная функция отправки сообщений
func sendMessage(bot *tgbotapi.BotAPI, chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	_, err := bot.Send(msg)
	if err != nil {
		log.Printf("Ошибка отправки сообщения: %v", err)
	}
}
