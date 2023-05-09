package main

import (
	"fmt"
	"latex-to-image-bot/latexToImage"
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("GO_TELEGRAM_APITOKEN"))
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil { // ignore any non-Message updates
			continue
		}

		if update.Message.IsCommand() {
			cmdReplayMsg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
			// respond to the commands
			switch update.Message.Command() {
			case "start":
				cmdReplayMsg.Text = "Hi there! Thanks for starting LaTeX to image bot. With my help, you can quickly and easily convert your LaTeX code into images that you can share with others. Just send me your LaTeX code and I'll take care of the rest."
			case "help":
				cmdReplayMsg.Text = `Just send me your LaTeX code and I'll convert it into an image for you.
                
I you want a code snippet to paste into your LaTeX document, use the following equation:
                
$$\text{Entropy } H =  \sum\limits_{i=1}^{n} -p(m_{i})\log_{2}(p(m_{i}))$$`
			case "contact":
				cmdReplayMsg.Text = "To see the source code of this bot, please visit https://github.com/taham8875/latex-to-image-bot"
			default:
				cmdReplayMsg.Text = "I don't know that command, supported commands are: /start, /help, /contact"
			}

			if _, err := bot.Send(cmdReplayMsg); err != nil {
				log.Panic(err)
			}

			continue

		}
		// respond to the user with a message
		waitMsg := tgbotapi.NewMessage(update.Message.Chat.ID, "converting your latex code to image...")
		waitMsg.ReplyToMessageID = update.Message.MessageID
		if _, err := bot.Send(waitMsg); err != nil {
			log.Panic(err)
		}
		// convert the user input to pdf
		outputFilePath, err := latexToImage.ConvertLatexToImage(update.Message.Text)
		if err != nil {
			fmt.Printf("Error: %v", err)
			return
		}

		// when the user send a message, reply with photo
		msg := tgbotapi.NewPhoto(update.Message.Chat.ID, tgbotapi.FilePath(outputFilePath))

		if _, err := bot.Send(msg); err != nil {
			log.Panic(err)
		}

		// delete the image file
		if err = os.Remove(outputFilePath); err != nil {
			fmt.Printf("Error deleting the image file: %v", err)
			return
		}

	}
}
