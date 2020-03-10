package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/go-chi/chi"
	"github.com/slack-go/slack"
)

func getenv(name string) string {
	v := os.Getenv(name)
	if v == "" {
		panic("missing required environment variable " + name)
	}
	return v
}

type IChiRouter interface {
	InitRouter() *chi.Mux
}

type router struct{}

func (router *router) InitRouter() *chi.Mux {
	r := chi.NewRouter()
	fmt.Println("Start service", 8080)
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		mess := map[string]interface{}{"success": true}
		json.NewEncoder(w).Encode(mess)
	})

	return r
}

var (
	m          *router
	routerOnce sync.Once
)

func ChiRouter() IChiRouter {
	if m == nil {
		routerOnce.Do(func() {
			m = &router{}
		})
	}
	return m
}

func main() {
	token := "xoxb-691975367441-783537757120-PgxkjCMcnT9SsojhtB9s8enw"
	api := slack.New(token)
	rtm := api.NewRTM()
	rand.Seed(time.Now().UnixNano())

	http.ListenAndServe(":80", ChiRouter().InitRouter())

	go rtm.ManageConnection()

	for {
		select {
		case msg := <-rtm.IncomingEvents:
			switch ev := msg.Data.(type) {

			case *slack.MessageEvent:
				info := rtm.GetInfo()

				text := ev.Text
				text = strings.TrimSpace(text)
				text = strings.ToLower(text)
				fmt.Println("Mess:", text)
				matchedHCC, _ := regexp.MatchString(strings.ToLower(info.User.ID), text)
				matchedHello, _ := regexp.MatchString("hello", text)
				matchedName, _ := regexp.MatchString("tên|ten|name", text)
				matchedRepo, _ := regexp.MatchString("repo|source|code", text)
				matchedJira, _ := regexp.MatchString("ticket|task|jira", text)
				matchedTag, _ := regexp.MatchString("tag", text)
				matchedPlay, _ := regexp.MatchString("choi", text)

				autoReply := []string{
					"Chờ xíu có người online rồi nói chuyện",
					"Có chuyện buồn không?",
					"Có chuyện j vui không?",
					"nói tiếp đi, Heo con nghe nè",
					"chùi ui",
					"nói cái gì, sao nữa?",
					"có ai iu Hy Heo không?",
					"đợi xíu Hy Heo rep liền",
					"cô đơn quá nè, nói gì đi",
					"vui quá vui quá",
					"sao sao",
					"kệ mẹ E",
					"nói cái nòi j vậy, tiếp deee",
					"ờ, rồi sao",
					"hôm nay ăn gì nhỉ?",
					"thay mặt Hy Heo, đang nghe",
					"nói đi",
					"thu đi để lại lá vàng \n Heo đi để lại bàng hoàng trong Hy",
					"Heo biết làm thơ đó",
					"có gì không nào",
					"đang nghe nè",
					"<@ULL51M6LF> nói j đi",
					"<@ULJ9FQMPS> đáng iu quá",
					"cái j nói đi, lát Hy mới onl",
					"ccc đố là gì",
					"Heo ăn nhìu không",
					"lâu quá không gặp",
					"nice to meat you",
					"ăn l*n rồi",
					"có j vui today",
					"corona muôn nơi",
					"cập nhật tình hình corona nào",
					"phải chăng E quá đáng iu",
				}

				if matchedHCC {
					if matchedHello {
						rtm.SendMessage(rtm.NewOutgoingMessage("hello cc", ev.Channel))
					} else if matchedRepo {
						rtm.SendMessage(rtm.NewOutgoingMessage("Ở đây `https://bitbucket.org/`", ev.Channel))
					} else if matchedName {
						rtm.SendMessage(rtm.NewOutgoingMessage("Tên mình là HCC, vai trò là supporter ", ev.Channel))
					} else if matchedTag {
						rtm.SendMessage(rtm.NewOutgoingMessage("<@ULL51M6LF> "+"<@ULJ9FQMPS> "+"<@UMYG2C805> ơi", ev.Channel))

					} else if matchedPlay {
						rtm.SendMessage(rtm.NewOutgoingMessage("chơi với <@ULL51M6LF> đi nè", ev.Channel))
					} else if matchedJira {
						rtm.SendMessage(rtm.NewOutgoingMessage("Đây `https://huylqbk.github.io`", ev.Channel))
					} else if text == "" {
						rtm.SendMessage(rtm.NewOutgoingMessage("hi <@"+ev.User+">", ev.Channel))
					} else {
						rtm.SendMessage(rtm.NewOutgoingMessage("là sao?", ev.Channel))
					}
				} else {
					rtm.SendMessage(rtm.NewOutgoingMessage(autoReply[rand.Intn(len(autoReply))], ev.Channel))
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
