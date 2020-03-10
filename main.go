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

func getenv(token string) string {
	return os.Getenv(token)
}

type IChiRouter interface {
	InitRouter() *chi.Mux
}

type router struct{}

func (router *router) InitRouter() *chi.Mux {
	r := chi.NewRouter()
	fmt.Println("Start service", 8080)
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
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

var session bool = false

type UserInfo struct {
	ID          string
	Name        string
	Old         int
	Job         string
	WorkAt      string
	Home        string
	Lonely      bool
	PreQuestion int
	Quesioned   map[int]bool
	Answered    map[int]bool
	PreAnswer   string
}

func main() {
	go slackRun()
	http.ListenAndServe(":8080", ChiRouter().InitRouter())
}

func slackRun() {
	token := getenv("TOKEN")
	if token == "" {
		token = "xoxb-691975367441-783537757120-t1skkNkDSpB39v224AIDdY7Y"
	}
	fmt.Println("TOKEN", token)
	api := slack.New(token)
	rtm := api.NewRTM()
	rand.Seed(time.Now().UnixNano())

	var users map[string]UserInfo = make(map[string]UserInfo)

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
				matchedMode, _ := regexp.MatchString("mode", text)
				matchedHCC, _ := regexp.MatchString(strings.ToLower(info.User.ID), text)
				matchedHello, _ := regexp.MatchString("hello|hi", text)
				matchedName, _ := regexp.MatchString("tên|ten|name", text)
				matchedRepo, _ := regexp.MatchString("repo|source|code|git", text)
				matchedJira, _ := regexp.MatchString("ticket|task|jira|team", text)
				matchedTag, _ := regexp.MatchString("tag", text)
				matchedPlay, _ := regexp.MatchString("choi", text)
				matchedYoutube, _ := regexp.MatchString("youtube|video", text)
				matchedMail, _ := regexp.MatchString("mail|thu", text)
				matchedSlack, _ := regexp.MatchString("slack", text)
				matchedGG, _ := regexp.MatchString("gg|search", text)
				matchedFB, _ := regexp.MatchString("fb|mxh|chat", text)
				matchedZalo, _ := regexp.MatchString("zalo", text)
				matchedHowTo, _ := regexp.MatchString("how to|how|error", text)
				matchedFind, _ := regexp.MatchString("f:", text)

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
					"nói đi <:)>",
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
					"khi cô đơn E nhớ ai?",
					"Hiiiiiiiiiiiiiiiiiiiiii",
					"Chào cô bé may mắn",
					"Rãnh đọc báo đi `https://vnexpress.net/suc-khoe`",
					"Đọc tin tức nè `https://vnexpress.net/thoi-su`",
					"Thế giới nay có gì `https://vnexpress.net/the-gioi` đọc đi",
					"Chứng khoáng như thế nào rồi `https://vnexpress.net/kinh-doanh` nè",
					"có phải E là thiên thần corona",
					"hôm nay thế nào baby?",
					"love u bặc bặc",
				}
				if matchedHCC {
					if matchedHello {
						rtm.SendMessage(rtm.NewOutgoingMessage("hello cc", ev.Channel))
					} else if matchedMode {
						if session {
							session = false
							rtm.SendMessage(rtm.NewOutgoingMessage("Mode: Auto", ev.Channel))
						} else {
							session = true
							rtm.SendMessage(rtm.NewOutgoingMessage("Mode: Session", ev.Channel))
							users[ev.User] = UserInfo{
								ID:        ev.User,
								Quesioned: map[int]bool{},
								Answered:  map[int]bool{},
							}
						}
					} else if matchedRepo {
						rtm.SendMessage(rtm.NewOutgoingMessage("Ở đây <https://github.com/huylqbk>", ev.Channel))
					} else if matchedName {
						rtm.SendMessage(rtm.NewOutgoingMessage("Tên mình là HCC, vai trò là supporter ", ev.Channel))
					} else if matchedTag {
						rtm.SendMessage(rtm.NewOutgoingMessage("<@ULL51M6LF> "+"<@ULJ9FQMPS> "+"<@UMYG2C805> ơi", ev.Channel))

					} else if matchedPlay {
						rtm.SendMessage(rtm.NewOutgoingMessage("chơi với <@ULL51M6LF> đi nè", ev.Channel))
					} else if matchedJira {
						rtm.SendMessage(rtm.NewOutgoingMessage("Đây <https://teams.microsoft.com>", ev.Channel))
					} else if text == "" {
						rtm.SendMessage(rtm.NewOutgoingMessage("hi <@"+ev.User+">", ev.Channel))

					} else if matchedYoutube {
						rtm.SendMessage(rtm.NewOutgoingMessage("open youtube <https://www.youtube.com/?gl=VN>", ev.Channel))
					} else if matchedMail {
						rtm.SendMessage(rtm.NewOutgoingMessage("open mail ne <https://outlook.office.com/mail/inbox>", ev.Channel))
					} else if matchedSlack {
						rtm.SendMessage(rtm.NewOutgoingMessage("open slack cty <https://app.slack.com/client/T891EANLE>", ev.Channel))
					} else if matchedGG {
						rtm.SendMessage(rtm.NewOutgoingMessage("search gi <https://google.com.vn>", ev.Channel))
					} else if matchedFB {
						rtm.SendMessage(rtm.NewOutgoingMessage("news feed <https://facebook.com>", ev.Channel))
					} else if matchedZalo {
						rtm.SendMessage(rtm.NewOutgoingMessage("zalo o day <https://chat.zalo.me>", ev.Channel))
					} else if matchedHowTo {
						str := text[12:]
						find := strings.Replace(str, " ", "+", -1)
						rtm.SendMessage(rtm.NewOutgoingMessage("solution: <https://stackoverflow.com/search?q="+find+">", ev.Channel))
					} else {
						rtm.SendMessage(rtm.NewOutgoingMessage("là sao?", ev.Channel))
					}
				} else if session {
					info := users[ev.User]
					if info.Name == "" && info.PreQuestion != 1 {
						q := "Ten ban la gi?"
						rtm.SendMessage(rtm.NewOutgoingMessage(q, ev.Channel))
						info.PreQuestion = 1
						info.Quesioned[1] = true
						users[ev.User] = info
					} else {
						if _, ok := info.Answered[1]; !ok {
							info.Name = text
							info.Answered[1] = true
							users[ev.User] = info
						}
						rtm.SendMessage(rtm.NewOutgoingMessage("Chuc "+info.Name+" vui ve!!! :)", ev.Channel))
					}
				} else if matchedFind {
					str := text[2:]
					find := strings.Replace(str, " ", "+", -1)
					rtm.SendMessage(rtm.NewOutgoingMessage("KQ=> <https://google.com/search?q="+find+">", ev.Channel))

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
