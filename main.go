package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	bot "github.com/MixinNetwork/bot-api-go-client"
	"github.com/julienschmidt/httprouter"
	qrcode "github.com/skip2/go-qrcode"
	"github.com/timshannon/badgerhold"
	"swap.xin/durable"
)

var database *durable.Database

var ClientID = ""
var SessionID = ""
var PrivateKey = ""
var Pin = ""
var PinToken = ""

func doGet(url string) ([]byte, error) {
	httpCli := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	resp, err := httpCli.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	return body, err
}

var exitTimer *time.Timer

var ctx = context.Background()

type User struct {
	Id   string `yaml:"id"`
	Name string `yaml:"name"`
	Uuid string `yaml:"uuid"`
}

var v struct {
	List []User `yaml:"list"`
}

var ymlData []byte

func AllHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	keyword := ps.ByName("keyword")
	if keyword == "favicon.ico" {
		return
	}

	var swapUser *SwapUser

	var findSwapUser []*SwapUser
	database.Find(&findSwapUser, badgerhold.Where("Keyword").Eq(keyword))
	if len(findSwapUser) > 0 {
		swapUser = findSwapUser[0]
	} else {
		http.Redirect(w, r, "/all", http.StatusFound)

		return
	}

	w.Header().Set("Content-Type", "text/html; charset = utf-8")

	htmlStr := fmt.Sprintf(`<!DOCTYPE html>
	<html lang="en">
		<head>
			<meta charset="UTF-8" />
			<meta name="viewport" content="width=device-width, initial-scale=1.0, viewport-fit=cover"
			/>
			<title>
				Mixin名片
			</title>
			<style>
				body, html { margin: 0; padding: 0; font-size: 16px; font-family: 'Quicksand',
				Helvetica, Arial, sans-serif; -webkit-font-smoothing: antialiased; -moz-osx-font-smoothing:
				grayscale; box-sizing: border-box; height: 100vh; } a { text-decoration:
				none; color: #135d98; font-weight: bold; cursor: pointer; } body { color:
				#2c3e50; background: url('https://mixin.one/assets/4aa9012849cc84b35004a005ab43ca2a.jpg')
				#366bd5; background-size: cover; background-repeat: no-repeat; } .wrapper
				{ max-width: 1200px; margin: 0 auto; padding: 6rem 1rem 2rem; box-sizing:
				border-box; line-height: 1.5; } @media (max-width: 760px) { .wrapper {
				padding: 2rem 1rem 2rem; } body { background: #366bd5; } html { font-size:
				14px; } } .qrcode-img { user-select: none; display: block; width: 12rem;
				margin: 0 auto; } .card { position: relative; background: #fff; border-radius:
				0.5rem; padding: 2rem 1rem 2rem; text-align: center; font-size: 1rem; }
				.avatar { width: 6rem; height: 6rem; border-radius: 6rem; border: 2px solid
				#fff; background: #fafafa; margin: 0 auto; overflow: hidden; } .avatar-img
				{ width: 6rem; height: 6rem; margin: 0 auto; } .qrcode-toggle { user-select:
				none; position: absolute; right: 1rem; top: 1rem; font-size: 0.9rem; }
				.name { font-weight: bold; font-size: 1.2rem; } .id { color: #555; } .reg
				{ } .uuid { } .hide { display: none; } .guide { text-align: center; font-size:
				0.8rem; padding: 0.5rem; color: #fff; } .guide a { color: #fff; text-decoration:
				underline; } .guide span { margin-right: 0.5rem; }
			</style>
		</head>
		<body>
			<div class="wrapper">
				<div class="card">
					<div class="qrcode-toggle" id="qrToggle">
						<a id="tt">
							转账
						</a>
						<a id="tc" class="hide">
							联系
						</a>
					</div>
					<div class="avatar">
						<img class="avatar-img" src="` + swapUser.AvatarURL + `" alt="头像" />
					</div>
					<div class="name">
						` + swapUser.Name + `
					</div>
					<div class="id">
						` + swapUser.Id + `
					</div>
					<div class="qrcode">
						<div id="qc">
							<img class="qrcode-img" src="` + swapUser.UserQrcode + `" alt="个人二维码" />
							<!-- <div>联系</div> -->
						</div>
						<div id="qt" class="hide">
							<img class="qrcode-img" src="` + swapUser.TransferQrcode + `" alt="转账二维码"
							/>
							<!-- <div>转账</div> -->
						</div>
					</div>
					<div class="uuid">
						` + swapUser.MixinId + `
					</div>
					<div class="reg">
						注册于: ` + strings.Split(swapUser.RegTime, ".")[0] + `
					</div>
				</div>
				<div class="guide">
					下载
					<a href="https://mixin.one/messenger" target="_blank">
						Mixin Messenger
					</a>
					<span>
						，以下任一方式可建立联系: 1. APP 内搜索
						<b>
							` + swapUser.Id + `
						</b>
					</span>
					<span>
						2. 使用 APP 扫码
					</span>
					<span>
						3. APP 内打开当前链接直达
					</span>
				</div>
			</div>
			<script>
				let inMixin = false
	if (
		window.webkit &&
		window.webkit.messageHandlers &&
		window.webkit.messageHandlers.MixinContext
	) {
		inMixin = prompt('MixinContext.getContext()')
	}
	if (window.MixinContext && window.MixinContext.getContext) {
		inMixin = window.MixinContext.getContext()
		const inMixinData = JSON.parse(inMixin)
		if (inMixinData.platform === 'Desktop') {
			inMixin = false
		}
	}
	if (inMixin) {
		window.open('mixin://users/` + swapUser.MixinId + `')
	}
	document.getElementById('qrToggle').onclick = () => {
		if (document.getElementById('tc').className === 'hide') {
			document.getElementById('tc').className = ''
			document.getElementById('tt').className = 'hide'
			document.getElementById('qc').className = 'hide'
			document.getElementById('qt').className = ''
		} else {
			document.getElementById('tt').className = ''
			document.getElementById('tc').className = 'hide'
			document.getElementById('qt').className = 'hide'
			document.getElementById('qc').className = ''
		}
	}
			</script>
		</body>
	
	</html>`)
	w.Write([]byte(htmlStr))
	return

}

func IndexHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	http.Redirect(w, r, "/all", http.StatusFound)
}

type SwapUser struct {
	Keyword        string `json:"keyword"`
	Name           string `json:"name"`
	Id             string `json:"id"`
	AvatarURL      string `json:"avatarURL"`
	RegTime        string `json:"regTime"`
	UpdateTime     int64  `json:"updateTime"`
	UserType       string `json:"userType"`
	MixinId        string `json:"mixinId"`
	OwnerMixinId   string `json:"ownerMixinId"`
	UserQrcode     string `json:"userQrcode"`
	TransferQrcode string `json:"transferQrcode"`
}

var client *bot.BlazeClient

// Handler is an implementation for interface bot.BlazeListener
// check out the url for more details: https://github.com/MixinNetwork/bot-api-go-client/blob/master/blaze.go#L89.
type Handler struct{}

// OnMessage is a general method of bot.BlazeListener
func (r Handler) OnMessage(ctx context.Context, msgView bot.MessageView, botID string) error {
	// I handle PLAIN_TEXT message only and make sure respond to current conversation.
	if msgView.Category == bot.MessageCategoryPlainText &&
		msgView.ConversationId == bot.UniqueConversationId(ClientID, msgView.UserId) {
		var data []byte
		var err error
		if data, err = base64.StdEncoding.DecodeString(msgView.Data); err != nil {
			log.Panicf("Error: %s\n", err)
			return err
		}
		inst := string(data)

		newMsg := bot.MessageView{
			UserId:         msgView.UserId,
			ConversationId: msgView.ConversationId,
		}
		instList := strings.Split(inst, " ")
		if len(instList) == 3 {
			if instList[0] == "n" {
				keyword := instList[2]
				id := instList[1]

				var findSwapUser []*SwapUser

				database.Find(&findSwapUser, badgerhold.Where("Keyword").Eq(keyword))

				if len(findSwapUser) > 0 {
					client.SendPlainText(ctx, newMsg, "当前关键词已被其他用户使用，如需申诉请联系：916716")
					return nil
				}

				currentUserData, _ := bot.GetUser(ctx, msgView.UserId, ClientID, SessionID, PrivateKey)

				userData, err := bot.SearchUser(ctx, id, ClientID, SessionID, PrivateKey)
				if err != nil {
					client.SendPlainText(ctx, newMsg, "当前id未搜索到相关信息")
					return nil
				}

				userType := ""
				ownerMixinId := ""

				if userData.App.CreatorId != "" && userData.App.CreatorId == currentUserData.UserId {
					userType = "app"
					ownerMixinId = currentUserData.UserId

				} else if userData.UserId == currentUserData.UserId {
					userType = "user"

				} else {
					client.SendPlainText(ctx, newMsg, "无权操作当前id")
					return nil
				}

				var findOldSwapUser []*SwapUser
				database.Find(&findOldSwapUser, badgerhold.Where("MixinId").Eq(userData.UserId))
				if len(findOldSwapUser) > 0 {
					client.SendPlainText(ctx, newMsg, "请先解除该账号绑定")
					return nil
				}

				preStr := ""

				userImage, _ := qrcode.Encode("mixin://users/"+userData.UserId, qrcode.Medium, 256)
				transferImage, _ := qrcode.Encode("mixin://transfer/"+userData.UserId, qrcode.Medium, 256)

				userQrcode := "data:image/png;base64," + base64.StdEncoding.EncodeToString(userImage)

				transferQrcode := "data:image/png;base64," + base64.StdEncoding.EncodeToString(transferImage)

				swapUser := &SwapUser{
					Keyword:   keyword,
					Name:      userData.FullName,
					Id:        userData.IdentityNumber,
					RegTime:   userData.CreatedAt,
					AvatarURL: userData.AvatarURL,
					// UpdateTime:     "",
					UserType:       userType,
					OwnerMixinId:   ownerMixinId,
					MixinId:        userData.UserId,
					UserQrcode:     userQrcode,
					TransferQrcode: transferQrcode,
				}

				database.Insert(swapUser)

				client.SendPlainText(ctx, newMsg, "成功获得专属链接：https://swap.xin/"+preStr+keyword)
			} else {
				client.SendPlainText(ctx, newMsg, "命令输入有误")
			}
		} else if inst == "info" {
			currentUserData, _ := bot.GetUser(ctx, msgView.UserId, ClientID, SessionID, PrivateKey)

			var findSwapUser []*SwapUser
			database.Find(&findSwapUser, badgerhold.Where("MixinId").Eq(currentUserData.UserId))

			var findOldSwapUser []*SwapUser
			database.Find(&findOldSwapUser, badgerhold.Where("OwnerMixinId").Eq(currentUserData.UserId))

			if len(findSwapUser) == 0 && len(findOldSwapUser) == 0 {
				client.SendPlainText(ctx, newMsg, "当前账号名下无链接")
				return nil
			} else {
				if len(findSwapUser) != 0 {
					client.SendPlainText(ctx, newMsg, "当前账号名片专属链接：https://swap.xin/"+findSwapUser[0].Keyword)
				}
				if len(findOldSwapUser) != 0 {
					str := ""
					for _, item := range findOldSwapUser {
						str = str + ("机器人" + item.Id + "专属链接：https://swap.xin/" + item.Keyword + "\n")
					}
					client.SendPlainText(ctx, newMsg, str)
				}
			}
		} else if len(instList) == 2 {
			id := instList[1]
			if instList[0] == "r" || instList[0] == "u" {
				currentUserData, _ := bot.GetUser(ctx, msgView.UserId, ClientID, SessionID, PrivateKey)
				userData, err := bot.SearchUser(ctx, id, ClientID, SessionID, PrivateKey)
				if err != nil {
					client.SendPlainText(ctx, newMsg, "当前id未搜索到相关信息")
					return nil
				}

				userType := ""
				ownerMixinId := ""
				if userData.App.CreatorId != "" && userData.App.CreatorId == currentUserData.UserId {
					userType = "app"
					ownerMixinId = currentUserData.UserId
				} else if userData.UserId == currentUserData.UserId {
					userType = "user"
				} else {
					client.SendPlainText(ctx, newMsg, "无权操作当前id")
					return nil
				}

				var findSwapUser []*SwapUser
				database.Find(&findSwapUser, badgerhold.Where("MixinId").Eq(userData.UserId))
				keyword := findSwapUser[0].Keyword

				var swapUser SwapUser
				database.Delete(&swapUser, badgerhold.Where("MixinId").Eq(userData.UserId))
				if instList[0] == "r" {
					userImage, _ := qrcode.Encode("mixin://users/"+userData.UserId, qrcode.Medium, 256)
					transferImage, _ := qrcode.Encode("mixin://transfer/"+userData.UserId, qrcode.Medium, 256)

					userQrcode := "data:image/png;base64," + base64.StdEncoding.EncodeToString(userImage)

					transferQrcode := "data:image/png;base64," + base64.StdEncoding.EncodeToString(transferImage)

					swapUser := &SwapUser{
						Keyword:   keyword,
						Name:      userData.FullName,
						Id:        userData.IdentityNumber,
						RegTime:   userData.CreatedAt,
						AvatarURL: userData.AvatarURL,
						// UpdateTime:     "",
						UserType:       userType,
						MixinId:        userData.UserId,
						OwnerMixinId:   ownerMixinId,
						UserQrcode:     userQrcode,
						TransferQrcode: transferQrcode,
					}

					database.Insert(swapUser)

					client.SendPlainText(ctx, newMsg, "刷新成功，名片专属链接：https://swap.xin/"+keyword)
				} else if instList[0] == "u" {
					client.SendPlainText(ctx, newMsg, "解绑成功")
				}
			}
		} else {
			client.SendPlainText(ctx, newMsg, "使用命令: \n\n1. 查看名下所有链接: info\n\n2. 注册链接: n [自己或自己机器人的id] [关键词]\n举例: n 916716 jim\n\n3. 刷新用户信息: r [id]\n\n4. 解除绑定: u [id]\n\n* 若关键词被占用请联系作者申诉")
		}

	}
	return nil
}

// Respond to user.
func Respond(ctx context.Context, msgView bot.MessageView, msg string) {
	if err := client.SendPlainText(ctx, msgView, msg); err != nil {
		log.Panicf("Error: %s\n", err)
	}
}

func main() {
	database = durable.OpenDatabaseClient()

	log.Println("start")
	client = bot.NewBlazeClient(ClientID, SessionID, PrivateKey)
	handler := Handler{}
	go func() {
		c := time.Tick(1 * time.Second)
		for {
			if err := client.Loop(ctx, handler); err != nil {
				log.Printf("Error: %v\n", err)
				os.Exit(3)
			}
			<-c
		}
	}()

	router := httprouter.New()
	router.GET("/:keyword", AllHandler)
	router.GET("/", IndexHandler)

	log.Fatal(http.ListenAndServe(":9919", router))
}
