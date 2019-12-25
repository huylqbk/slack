package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/nlopes/slack"
)

func getenv(name string) string {
	v := os.Getenv(name)
	if v == "" {
		panic("missing required environment variable " + name)
	}
	return v
}

func main() {
	token := "xoxb-639493558084-785415366871-n86QqVnsdNmRIGEotiuceAyt"
	api := slack.New(token)
	rtm := api.NewRTM()
	go rtm.ManageConnection()

	for {
		select {
		case msg := <-rtm.IncomingEvents:
			fmt.Print("Event Received: ")
			switch ev := msg.Data.(type) {

			case *slack.MessageEvent:
				info := rtm.GetInfo()

				text := ev.Text
				text = strings.TrimSpace(text)
				text = strings.ToLower(text)

				matchedHello, _ := regexp.MatchString("hello", text)
				matchedName, _ := regexp.MatchString("tên|ten|name", text)
				matchedRepo, _ := regexp.MatchString("repo|source|code", text)
				matchedJira, _ := regexp.MatchString("ticket|task|jira", text)
				matchedTag, _ := regexp.MatchString("tag", text)

				if ev.User != info.User.ID {
					if matchedHello {
						rtm.SendMessage(rtm.NewOutgoingMessage("hello cc", ev.Channel))
					} else if matchedRepo {
						rtm.SendMessage(rtm.NewOutgoingMessage("Ở đây `https://bitbucket.org/exgo-tech`", ev.Channel))
					} else if matchedName {
						rtm.SendMessage(rtm.NewOutgoingMessage("Tên mình là Exgo, vai trò là supporter ", ev.Channel))
					} else if matchedTag {
						rtm.SendMessage(rtm.NewOutgoingMessage("@huy.le", ev.Channel))

					} else if matchedJira {
						rtm.SendMessage(rtm.NewOutgoingMessage("Đây `https://exgo.atlassian.net/secure/BrowseProjects.jspa`", ev.Channel))
					} else {
						rtm.SendMessage(rtm.NewOutgoingMessage("là sao?", ev.Channel))
					}
				}

			case *slack.RTMError:
				fmt.Printf("Error: %s\n", ev.Error())

			case *slack.InvalidAuthEvent:
				fmt.Printf("Invalid credentials")
				break

			default:

			}
		}
	}
}
