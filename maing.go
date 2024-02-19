type Update struct {
	UpdateId int     `json:"update_id"`
	Message  Message `json:"message"`
}
type Message struct {
	Text string `json:"text"`
	Chat Chat   `json:"chat"`
}
type Chat struct {
	Id int `json:"id"`
}

func parsTelegramRequest(r *http.Request) (*Update, error) {
	var update Update
	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		log.Printf("could not decode incoming update %s ", err.Error())
		return nil, err
	}
	return &update, nil
}
func handleTelegramWebHook(w http.ResponseWriter, r *http.Request) {
	var update, err = parsTelegramRequest(r)
	if err != nil {
		log.Printf("error parsing update, %s", err.Error())
		return
	}
	var telegramResponseBody, errTelegram = sendTextToTelegramChat(update.Message.Chat.Id, lyric)
	if errTelegram != nil {
		log.Printf("got error %s from telegram, response body is %s", errTelegram.Error(), telegramResponseBody)
	} else {
		log.Printf("punchline %s successfully distributed to chat id %d", lyric, update.Message.Chat.Id)
	}
}

func sendTextToTelegramChat(chatId int, text string) (string, error) {
	log.Printf("Sending %s to chat_id: %d", text, chatId)
	var telegramApi string = "https://api.telegram.org/bot" + os.Getenv("TELEGRAM_BOT_TOKEN") + "/sendMessage"
	response, err := http.PostForm(
		telegramApi,
		url.values{
			"chat_id": {strconv.Itoa(chatId)},
			"text":    {text},
		})
	if err != nil {
		log.Printf("error when posting text to the chat: %s", err.Error())
		return "", err
	}
	defer response.Body.Close()
	var bodyBytes, errRead = ioutilReadAll(response.Body)
	if errRead != nil {
		log.printf("error in parsing telegram answer %s", errRead.Error())
		return "", err
	}
	bodyString := string(bodyBytes)
	log.Printf("Body of Telegram Response:%s", bodyString)
	return bodyString, nil
}