package main

import (
	"fmt"
	"github.com/Syfaro/telegram-bot-api"
	"github.com/dafanasev/go-yandex-translate"
	"golang.org/x/net/proxy"
	"log"
	"net/http"
	"os"
	"reflect"
)

func main() {
	os.Setenv("YANDEX_APIKEY", "YANDEX_API_KEY")
	os.Setenv("TG_BOT_TOKEN", "YOUR_TELEGRAM_BOT_TOKEN")
	os.Setenv("PROXY_HOST", "YOUR_PROXY_SOCKS5")
	//proxyStr := "178.159.36.10:9050"
	//proxyURL, err := url.Parse(proxyStr)
	//if err != nil {
	//	log.Println(err)
	//}
	trnslt := translate.New(os.Getenv("YANDEX_APIKEY"))
	dialSocksProxy, err := proxy.SOCKS5("tcp", os.Getenv("PROXY_HOST"), nil, proxy.Direct)

	if err != nil {
		fmt.Println("Error connecting to proxy:", err)
	}

	tr := &http.Transport{Dial: dialSocksProxy.Dial}

	// Create client
	myClient := &http.Client{
		Transport: tr,
	}

	//mClient := &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxyURL)}}
	//http.DefaultTransport = &http.Transport{Proxy: http.ProxyURL(proxyURL)}

	//bot, err := tgbotapi.NewBotAPI("311927713:AAE8BOsoajS7TTMU87swuEfkPhmIlBV5_Xo")
	bot, err := tgbotapi.NewBotAPIWithClient(os.Getenv("TG_BOT_TOKEN"), myClient)

	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true
	if bot.Debug {
		log.Printf("Authorized on account %s", bot.Self.UserName)
	}

	var ucfg tgbotapi.UpdateConfig = tgbotapi.NewUpdate(1)
	ucfg.Timeout = 60
	updates, err := bot.GetUpdatesChan(ucfg)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		if bot.Debug {
			log.Printf("[%s]%s", update.Message.From.UserName, update.Message.Text)
		}

		if reflect.TypeOf(update.Message.Text).Kind() == reflect.String && update.Message.Text != "" {

			if update.Message.Text == "/start" {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Привет! Напиши мне любую фразу на русском языке, и я отвечу тебе сообщением с перводом на английском. =) Попробуй)")
				bot.Send(msg)
				continue
			} else {
				translated := trans(update.Message.Text, trnslt)
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, translated)
				bot.Send(msg)
				continue
			}

		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Не, давай текст")

		//msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		//msg.ReplyToMessageID = update.Message.MessageID

		bot.Send(msg)
	}
}

func sendGet(client http.Client, url string) (resp *http.Response) {
	resp, err := client.Get(url)
	if err != nil {
		fmt.Println("Error", err)
	}

	return resp
}

func trans(text string, tr *translate.Translator) string {
	response, err := tr.GetLangs("ru")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(response.Langs)
		fmt.Println(response.Dirs)
	}

	translation, err := tr.Translate("en", text)
	if err != nil {
		fmt.Println(err)
	}

	return translation.Result()
}
