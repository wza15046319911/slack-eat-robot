package api

import (
	"context"
	"eat-and-go/config"
	"eat-and-go/gorm"
	"eat-and-go/handler"
	"eat-and-go/model"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"go.mongodb.org/mongo-driver/bson"
	"io/ioutil"
	"log"
	r "math/rand"
	"strings"
	"time"
)

type SlackDefination struct {
	SlackAPI         *slack.Client
	Recommendation   *model.DianpingItem
	MessageTimestamp string
}

var api *SlackDefination

const (
	EmptyMentionMessage  = "芜湖，总得发点什么，那就给你随便推荐一个吧！"
	HelpMessage          = "@我，并且试试这样问我：火锅、寿司、韩国料理"
	NoMatchMessage       = "芜湖，好像没有匹配到，那就给你随便推荐一个吧！"
	AfterResponseMessage = "✅ 感谢反馈～"
)

func InitSlack() {
	slackConfig := config.GetConfig().Slack
	botToken := slackConfig.SlackBotToken
	api = &SlackDefination{
		SlackAPI: slack.New(botToken),
	}
	fmt.Printf("slack初始化成功\n")
}

//DispatchSlackEvent 处理初次连接到slack的请求以及app mention的请求
// @Summary       api
// @Description   Handle APP_MENTION event or other call back event sent from slack
// @Tags          event
// @Accept        json
// @Success       200   {object} handler.Response
// @Router        /slack [post]
func DispatchSlackEvent(context *gin.Context) {
	body, err := ioutil.ReadAll(context.Request.Body)
	if err != nil {
		log.Fatal("read err", err)
		return
	}
	eventsAPIEvent, err := slackevents.ParseEvent(body, slackevents.OptionNoVerifyToken())
	if err != nil {
		log.Fatal("event api err:", err)
		return
	}
	fmt.Println()
	if eventsAPIEvent.Type == slackevents.URLVerification {
		var c *slackevents.ChallengeResponse
		err := json.Unmarshal(body, &c)
		if err != nil {
			log.Fatal("500 internal error")
			return
		}
		fmt.Printf(c.Challenge)
		// Respond with 200 OK
		context.String(200, c.Challenge)
	} else if eventsAPIEvent.Type == slackevents.CallbackEvent {
		innerEvent := eventsAPIEvent.InnerEvent
		switch ev := innerEvent.Data.(type) {
		case *slackevents.AppMentionEvent:
			err := SlackMentionHandler(ev)
			if err != nil {
				log.Fatal("message sending failed, err: ", err)
			}
			fmt.Printf("信息已经成功发送到频道 %s\n", ev.Channel)
			attachment := getAttachments()
			optionAttachments := slack.MsgOptionAttachments(attachment)
			_, ts, _ := api.SlackAPI.PostMessage(ev.Channel, slack.MsgOptionText("", false), optionAttachments)
			api.MessageTimestamp = ts
		}
	}
}

//SlackMentionHandler 处理app metion的事件
//提取关键词进行店铺推荐。如果关键词为空/无效，或者没有匹配到内容，则随机返回数据库中的一家店铺link
//如果有匹配，随机返回所有满足条件的一家店铺link。
func SlackMentionHandler(ev *slackevents.AppMentionEvent) error {
	var messageFromChannel = extractText(ev.Text)
	var messageSendToChannel string
	randomResults, err := RandomSelect()
	if err != nil {
		return err
	}
	if messageFromChannel == "" {
		// 如果没有发消息，随即返回一个店铺，并且给出操作提示
		api.SlackAPI.PostMessage(ev.Channel, slack.MsgOptionText(EmptyMentionMessage, false))
		api.Recommendation = &randomResults[0]
		api.SlackAPI.PostMessage(ev.Channel, slack.MsgOptionText(HelpMessage, false))
	} else {
		// 有消息，试着匹配
		documents, err := getAllDocuments()
		if err != nil {
			return err
		}
		selectedItem := SelectByKeyword(messageFromChannel, documents)
		if err != nil {
			return err
		}
		api.Recommendation = &selectedItem
		if api.Recommendation.DetailLink == "" {
			// 如果没有匹配，随即返回
			api.SlackAPI.PostMessage(ev.Channel, slack.MsgOptionText(NoMatchMessage, false))
			api.Recommendation = &randomResults[0]
		}
	}
	messageSendToChannel = api.Recommendation.DetailLink
	api.SlackAPI.PostMessage(ev.Channel, slack.MsgOptionText(messageSendToChannel+"\n\n", false))
	return nil
}

//ButtonAttachmentsHandler 处理按钮点击事件
//该方法会收集用户满意度并保存至数据库，为日后的数据分析和算法学习
//作准备
func ButtonAttachmentsHandler(context *gin.Context) {
	var payload slack.InteractionCallback
	formValue, ok := context.GetPostForm("payload")
	if !ok {
		handler.SendResponse400(context, errors.New("payload not found"), nil)
	} else {
		err := json.Unmarshal([]byte(formValue), &payload)
		if err != nil {
			handler.SendResponse400(context, err, nil)
			return
		}
		//var userName = payload.User.Name
		var userId = payload.User.ID
		//var userChoice = payload.ActionCallback.AttachmentActions[0].Value

		fmt.Printf("Recommended item: %s", api.Recommendation.ShopID)
		var user model.SlackRecommendation
		db := gorm.DB.Self
		user = model.SlackRecommendation{
			UserId:            userId,
			UserName:          payload.User.Name,
			RecommendShopId:   api.Recommendation.ShopID,
			RecommendShopName: api.Recommendation.ShopName,
			RecommendShopLink: api.Recommendation.DetailLink,
			Preference:        payload.ActionCallback.AttachmentActions[0].Value,
		}
		res := db.Create(&user)
		if res.Error != nil {
			handler.SendResponse400(context, err, nil)
			return
		}
		optionalAttachment := slack.MsgOptionAttachments(
			slack.Attachment{
				CallbackID: "accept_or_reject",
				ID:         0,
				AuthorName: "slack eat-and-go",
				Text:       AfterResponseMessage,
				Actions:    []slack.AttachmentAction{},
			})
		fmt.Printf("成功更新用户%s的喜好\n", payload.User.Name)
		api.SlackAPI.UpdateMessage(
			payload.Channel.ID,
			api.MessageTimestamp,
			slack.MsgOptionText("", false),
			optionalAttachment,
		)
	}

}

//RandomSelect 随机返回一个店铺的信息
func RandomSelect() ([]model.DianpingItem, error) {
	var results []model.DianpingItem
	var result model.DianpingItem
	var db = gorm.Collections.RegoCollection
	pipeline := []bson.D{{{"$sample", bson.D{{"size", 1}}}}}
	cur, err := db.Aggregate(context.TODO(), pipeline)
	if err != nil {
		return nil, err
	}
	for cur.Next(context.TODO()) {
		err = cur.Decode(&result)
		if err != nil {
			return nil, err
		}
		results = append(results, result)
	}
	return results, nil
}

//getAllDocuments 返回dianping数据表中所有数据
func getAllDocuments() ([]model.DianpingItem, error) {
	var results []model.DianpingItem
	var result model.DianpingItem
	var db = gorm.Collections.RegoCollection
	cur, err := db.Find(context.TODO(), bson.D{{}})
	if err != nil {
		return results, err
	}
	for cur.Next(context.TODO()) {
		err := cur.Decode(&result)
		if err != nil {
			fmt.Println(result.ShopID)
			return results, err
		}
		results = append(results, result)
	}
	return results, nil
}

//SelectByKeyword 根据关键字进行店铺不完全匹配，随即返回一个满足的店铺。
//根据标签1、店铺名称和推荐菜，寻找匹配的店铺
func SelectByKeyword(keyword string, data []model.DianpingItem) model.DianpingItem {
	var satisfiedItems []model.DianpingItem
	for _, item := range data {
		var tag = item.TagOne
		var shopName = item.ShopName
		var recommend = item.Recommend
		if strings.Index(tag, keyword) != -1 {
			satisfiedItems = append(satisfiedItems, item)
		} else if strings.Index(shopName, keyword) != -1 {
			satisfiedItems = append(satisfiedItems, item)
		} else if strings.Index(recommend, keyword) != -1 {
			satisfiedItems = append(satisfiedItems, item)
		}
	}
	fmt.Printf("共有 %d 家店铺满足条件\n", len(satisfiedItems))
	if len(satisfiedItems) == 0 {
		return model.DianpingItem{}
	}
	r.Seed(time.Now().Unix())
	var randomIndex = r.Intn(len(satisfiedItems))
	return satisfiedItems[randomIndex]
}

func extractText(message string) string {
	authUser := config.GetConfig().Slack.AuthUser
	newMessage := strings.Replace(message, authUser, "", -1)
	newMessage = strings.Trim(newMessage, " ")
	return newMessage
}

//getAttachments slack吃饭机器人推送的问卷调查
func getAttachments() slack.Attachment {
	return slack.Attachment{
		Color: "#3AA3E3",
		//Fallback:   "We don't know the future",
		CallbackID: "accept_or_reject",
		ID:         0,
		AuthorName: "slack eat-and-go",
		Title:      "满意度小调查",
		Text:       "对我反馈的结果满意嘛？",
		Actions: []slack.AttachmentAction{
			{
				Name:  "yes",
				Text:  "yes!",
				Type:  "button",
				Value: "yes",
			},
			{
				Name:  "no",
				Text:  "no?",
				Type:  "button",
				Value: "no",
				Style: "danger",
			},
		},
	}
}
