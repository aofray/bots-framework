package viber_bot

import (
	"net/http"
	"fmt"
	"github.com/strongo/bots-api-viber"
	"net/url"
	"github.com/strongo/app/log"
	"github.com/julienschmidt/httprouter"
)

func (h ViberWebhookHandler) SetWebhook(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	c := h.Context(r)
	client := h.GetHttpClient(c)
	botCode := r.URL.Query().Get("code")
	if botCode == "" {
		http.Error(w, "Missing required parameter: code", http.StatusBadRequest)
		return
	}
	botSettings, ok := h.botsBy(c).ByCode[botCode]
	if !ok {
		http.Error(w, fmt.Sprintf("Bot not found by code: %v", botCode), http.StatusBadRequest)
		return
	}
	bot := viberbotapi.NewViberBotApiWithHttpClient(botSettings.Token, client)
	//bot.Debug = true

	webhookUrl := fmt.Sprintf("https://%v/bot/viber/callback/%v", r.Host, url.QueryEscape(botSettings.Code))

	//eventTypes := []string {"delivered", "seen", "failed", "subscribed",  "unsubscribed", "conversation_started"}
	eventTypes := []string {"failed", "subscribed",  "unsubscribed", "conversation_started"}

	if _, err := bot.SetWebhook(webhookUrl, eventTypes); err != nil {
		log.Errorf(c, "%v", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	} else {
		w.Write([]byte("Webhook set"))
	}
}

