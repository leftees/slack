package slack

import (
	"fmt"
	"regexp"
)

func (bot *Bot) RespondRegexp(re *regexp.Regexp, handler BotAction) {
	namePattern := fmt.Sprintf("\\A%s|<@%s>:? ", bot.Name, bot.ID)
	nameRe := regexp.MustCompile(namePattern)
	closure := func(self *Bot, event map[string]interface{}) (*Message, Status) {
		text := event["text"].(string)
		match := nameRe.FindStringIndex(text)
		if match == nil {
			return nil, CONTINUE
		}
		unmatchedText := text[match[1]+1:]
		if re.MatchString(unmatchedText) {
			return handler(self, event)
		}
		return nil, CONTINUE
	}
	messageHandlers, ok := bot.Handlers["message"]
	if !ok {
		messageHandlers = make([]BotAction, 0)
	}
	messageHandlers = append(messageHandlers, closure)
	bot.Handlers["message"] = messageHandlers
}

func (bot *Bot) Respond(text string, handler BotAction) {
	re := regexp.MustCompile(text)
	bot.RespondRegexp(re, handler)
}
