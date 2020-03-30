package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image"
	"image/jpeg"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
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

var randomWords = []string{}

type DataWords struct {
	Data []string `json:"data"`
}

func randomWord2000() []string {
	url := "https://www.randomlists.com/data/words.json"
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		return nil
	}
	res, err := client.Do(req)
	defer res.Body.Close()
	var data DataWords
	json.NewDecoder(res.Body).Decode(&data)
	return data.Data
}

func main() {
	randomWords = randomWord2000()
	go slackRun()
	http.ListenAndServe(":8080", ChiRouter().InitRouter())
}

func slackRun() {
	token := getenv("TOKEN")
	if token == "" {
		token = "xoxb-691975367441-783537757120-Vwv5QRBhXd5iOSNFuoCtNmND"
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
				matchedHello, _ := regexp.MatchString("hello|helo|halo", text)
				matchedHappy, _ := regexp.MatchString("hehe|hihi|kk", text)
				matchedBye, _ := regexp.MatchString("pipi|bye|bibi", text)
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
				matchedEng, _ := regexp.MatchString("e:", text)
				matchedWord, _ := regexp.MatchString("w:", text)
				matchedImg, _ := regexp.MatchString("i:", text)
				matchedCC, _ := regexp.MatchString("cc|fuck|dm|lon|cl|dcm|dcmm|lol|cặt", text)
				matchedTranslate, _ := regexp.MatchString("d:", text)
				matchedTranslateEng, _ := regexp.MatchString("t:", text)
				matchedCorona, _ := regexp.MatchString("corona", text)
				matchedHy, _ := regexp.MatchString("hy|huy", text)
				matchedIntro, _ := regexp.MatchString("giới thiệu|gioi thieu", text)

				autoReply := []string{
					"Chờ xíu có người online rồi nói chuyện",
					"Có chuyện buồn không?",
					"Có chuyện j vui không?",
					"nói tiếp đi, con nghe nè",
					"chùi ui",
					"nói cái gì, sao nữa?",
					"có ai iu con không?",
					"đợi xíu Hy rep liền",
					"cô đơn quá nè, nói gì đi",
					"vui quá vui quá",
					"sao sao",
					"kệ mẹ E",
					"nói cái nòi j vậy, tiếp deee",
					"ờ, rồi sao",
					"hôm nay ăn gì nhỉ?",
					"thay mặt Hy , đang nghe",
					"nói đi <:)>",
					"thu đi để lại lá vàng \n <@" + ev.User + "> đi để lại bàng hoàng trong Hy",
					"Mình biết làm thơ đó",
					"có gì không nào",
					"đang nghe nè",
					"<@ULL51M6LF> nói j đi",
					"<@" + ev.User + "> đáng iu quá",
					"cái j nói đi, lát Hy mới onl",
					"ccc đố là gì",
					" ăn nhìu không",
					"lâu quá không gặp",
					"nice to meat you",
					"ăn l*n rồi",
					"có j vui today",
					"corona muôn nơi",
					"cập nhật tình hình corona nào",
					"phải chăng  con quá đáng iu",
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
					"rãnh quá haaaa",
					"nói xàm xàm đi",
					"Đi ăn gì đi ha",
					"Muốn đi du lịch đâu đó quá",
					"Rãnh ngồi học tiếng anh đi",
					"sometime I miss u so much, hehe",
					"Sao tự kỉ đến mức nói với tôi vậy",
					"Mình cũng thông minh lắm đó",
					"Mình hong biết nữa",
					"kể nghe câu chuyện vui đi",
					"hôm nay nhiều task để làm không?",
				}

				hellos := []string{
					"Chào cậu, mình là  con con",
					"Chào cái cc",
					"hello cô bé",
					"Chào bé <@" + ev.User + "> đáng iu",
					"hi <@" + ev.User + ">",
					"Tớ cũng chào cậu",
				}

				byes := []string{
					"hẹn gặp lại <@" + ev.User + ">",
					"không tiễn, bibi",
					"see you again",
					"huhu, đi đâu vậy",
					"pipi, Đừng quên  này",
					"Rãnh nhắn tin với  tiếp nhé",
					"Nơi này cô đơn lắm, rãnh nhớ ghé",
					"chia ly từ đây, bye uuu",
					"miss you so much!!!",
					"byeeeeeeeeeeeeeeeee !!!!!! <:)>",
				}

				ccs := []string{
					"fuck u baby",
					"cc đây nè 8=========>",
					"where do you fuck me?",
					"fuck me quickly fuck me, please!!",
					"cc cl dm...",
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
					} else if matchedIntro {
						rtm.SendMessage(rtm.NewOutgoingMessage("Mình là đầy tới ở đây, tên ConCon, Hy Lê là cha mình", ev.Channel))
					} else if matchedRepo {
						rtm.SendMessage(rtm.NewOutgoingMessage("Ở đây <https://github.com/huylqbk>", ev.Channel))
					} else if matchedName {
						rtm.SendMessage(rtm.NewOutgoingMessage("Tên mình là con con, vai trò là siêu nhân giải cứu thế giới ", ev.Channel))
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
						rtm.SendMessage(rtm.NewOutgoingMessage("là sao? không hiểu", ev.Channel))
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

				} else if matchedEng {
					w := getRandomWork()
					if w == "" {
						rtm.SendMessage(rtm.NewOutgoingMessage("server lỗi rồi, học sau nhé", ev.Channel))
					} else {
						rtm.SendMessage(rtm.NewOutgoingMessage("Từ <"+w+"> nghĩa là gì?", ev.Channel))
					}
				} else if matchedTranslate {
					text = text[2:]
					r, e := getTranlateVN(text)
					if e != nil {
						rtm.SendMessage(rtm.NewOutgoingMessage("server lỗi rồi, dịch sau nhé", ev.Channel))
					} else {
						rtm.SendMessage(rtm.NewOutgoingMessage("English=> "+r, ev.Channel))
					}
				} else if matchedTranslateEng {
					text = text[2:]
					r, e := getTranlateEng(text)
					if e != nil {
						rtm.SendMessage(rtm.NewOutgoingMessage("server lỗi rồi, dịch sau nhé", ev.Channel))
					} else {
						rtm.SendMessage(rtm.NewOutgoingMessage("Tiếng Việt=> "+r, ev.Channel))
					}
				} else if matchedWord {
					if len(randomWords) == 0 {
						rtm.SendMessage(rtm.NewOutgoingMessage("server lỗi rồi, dùng sau nhé", ev.Channel))
					} else {
						w1 := randomWords[rand.Intn(len(randomWords))]
						w2 := randomWords[rand.Intn(len(randomWords))]
						w3 := randomWords[rand.Intn(len(randomWords))]
						w4 := randomWords[rand.Intn(len(randomWords))]
						w5 := randomWords[rand.Intn(len(randomWords))]
						r := fmt.Sprintf("%s, %s, %s, %s, %s", w1, w2, w3, w4, w5)
						m, _ := getTranlateEng(r)
						rtm.SendMessage(rtm.NewOutgoingMessage("Random=> "+r, ev.Channel))
						rtm.SendMessage(rtm.NewOutgoingMessage("Means=> "+m, ev.Channel))
					}

				} else if matchedCorona {

					corona := getCoronaNews()
					news := []string{}
					for i, c := range corona.Data.TopTrueNews {
						new := fmt.Sprintf("Tin So %d\nTitle: %s\nAuthor: %s\nURL: %s\n\n\n", i+1, c.Title, c.Author, c.Url)
						news = append(news, new)
					}
					rtm.SendMessage(rtm.NewOutgoingMessage("Tin tuc Corona Moi, cap nhat lien tuc\n"+strings.Join(news, ""), ev.Channel))

				} else if matchedImg {
					rtm.SendMessage(rtm.NewOutgoingMessage("here <https://www.facebook.com/search/top/?q="+text[2:]+"&epa=SEARCH_BOX>", ev.Channel))
				} else if matchedCC {
					rtm.SendMessage(rtm.NewOutgoingMessage(ccs[rand.Intn(len(ccs))], ev.Channel))
				} else if matchedHello {
					rtm.SendMessage(rtm.NewOutgoingMessage(hellos[rand.Intn(len(hellos))], ev.Channel))
				} else if matchedBye {
					rtm.SendMessage(rtm.NewOutgoingMessage(byes[rand.Intn(len(byes))], ev.Channel))
				} else if matchedHappy {
					rtm.SendMessage(rtm.NewOutgoingMessage("cười cười cc", ev.Channel))
				} else if matchedHy {
					rtm.SendMessage(rtm.NewOutgoingMessage("Hy là người đẹp trai, dễ mến và tạo ra tôi đó", ev.Channel))

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

func getRandomWork() string {
	var word []string
	r, err := http.Get("https://random-word-api.herokuapp.com/word?number=1")
	if err != nil {
		return ""
	}
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return ""
	}
	err = json.Unmarshal(body, &word)
	if err != nil {
		return ""
	}
	return word[0]
}

func getRandomImg() image.Image {
	r, err := http.Get("https://source.unsplash.com/random")
	if err != nil {
		return nil
	}
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	var fieldMapForBasic map[string]*json.RawMessage
	json.Unmarshal(body, &fieldMapForBasic)
	image, _ := json.Marshal(fieldMapForBasic["image"])
	coI := strings.Index(string(image), ",")
	rawImage := string(image)[coI+1:]
	unbased, _ := base64.StdEncoding.DecodeString(string(rawImage))
	jpgI, errJpg := jpeg.Decode(bytes.NewReader(unbased))
	fmt.Println(string(body))
	if errJpg != nil {
		return nil
	}
	return jpgI
}

type ResultTranlate struct {
	Status bool   `json:"status"`
	Msg    string `json:"msg"`
	Html   string `json:"html"`
}

func getTranlateVN(s string) (string, error) {
	data := url.Values{
		"text": {s},
	}
	r, err := http.PostForm("https://www.tienganh123.com/ajax/translate_sentences/result", data)
	if err != nil {
		return "", err
	}
	var result ResultTranlate
	err = json.NewDecoder(r.Body).Decode(&result)
	if err != nil {
		return "", err
	}

	text := result.Html[616:]
	text = text[:strings.Index(text, "</div>")]

	return text, nil
}

func getTranlateEng(s string) (string, error) {
	input := strings.Replace(s, " ", "+", -1)
	input = input + "&lang_from=en&lang_to=vi"

	url := "https://dich.tienganh123.com/translate_sentences/result"
	method := "POST"

	payload := strings.NewReader("text=" + input + "&lang_from=en&lang_to=vi")

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		return "", err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")

	res, err := client.Do(req)
	defer res.Body.Close()
	var result ResultTranlate
	err = json.NewDecoder(res.Body).Decode(&result)
	if err != nil {
		return "", err
	}

	text := result.Html[616:]
	text = text[:strings.Index(text, "</div>")]
	return text, nil
}

type Corona struct {
	Data ListNews `json:"data"`
}

type ListNews struct {
	TopTrueNews []TopTrueNews
}

type TopTrueNews struct {
	Type    string `json:"type"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Url     string `json:"url"`
	Author  string `json:"author"`
}

func getCoronaNews() Corona {
	url := "https://corona-api.kompa.ai/graphql"
	method := "POST"

	payload := strings.NewReader("{\"operationName\": \"topTrueNews\",\"variables\": {},\"query\": \"query topTrueNews{topTrueNews {id type title content url siteName publishedDate author picture}}\"}")

	client := &http.Client{}
	req, _ := http.NewRequest(method, url, payload)
	req.Header.Add("authority", "corona-api.kompa.ai")
	req.Header.Add("content-type", "application/json")
	req.Header.Add("origin", "https://corona.kompa.ai")

	res, err := client.Do(req)
	if err != nil {
		return Corona{}
	}
	defer res.Body.Close()
	var corona Corona
	json.NewDecoder(res.Body).Decode(&corona)
	return corona
}
