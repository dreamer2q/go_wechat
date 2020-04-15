package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"time"
	"wechat/api"
	"wechat/api/media"
	"wechat/api/request"
	"wechat/config"
)

func main() {

	c := &request.Config{
		AppID:        config.AppID,
		AppSecret:    config.AppSecret,
		AppToken:     config.AppToken,
		AesEncodeKey: config.AesEncodeKey,
		Callback:     "wx",
		Timeout:      10 * time.Second,
	}
	r := request.New(c)
	//r.SetToken(`32_q2Yyqlsha3jO9rC_IAqtBIKUwHi0hnP8_z2x07SIMQhP6lJeD6eotET5Ol_BRl-R5O1Ssp72XUF7VF7Y-hs3VXi8IShvYrApH5PoOXiIafVPXnthR1KZtJ7wrzngrwXEqeeieh2PQft0mWSTFFZfAFAAAU`)

	m := media.New(r)

	//err := m.SendPreview(&media.Message{
	//	ToWxName: "dreamer2q",
	//	MsgWrapper: &media.MpNews{
	//		MediaID: "bGL97-3t2bSmRVZowC4qmL2nV3flhA-NNRLXe9XvdQg",
	//	},
	//})
	////bGL97-3t2bSmRVZowC4qmIyYbtpLHluojxs6bIsFGDk
	//ret, err := m.UploadArticle(
	//	&media.ArticleWrapper{
	//		Articles: []media.NewsItem{
	//			{
	//				Title:        "New Title Test",
	//				ThumbMediaID: "bGL97-3t2bSmRVZowC4qmMYa_8Xx6dSQ4HhvFFM9DSU",
	//				ShowCoverPic: 1,
	//				Author:       "Admin",
	//				Digest:       "Im the digest content",
	//				Content:      config.ContentTest,
	//				SourceURL:    "https://dreamer2q.wang/",
	//				Comment:      1,
	//			},
	//		},
	//	})
	//mCounter, err := m.MaterialCounter()
	//mList, err := m.GetMaterialList(media.TypImage, 0, 10)
	//
	picBytes, err := ioutil.ReadFile("assets/123.png")
	if err != nil {
		log.Panic(err)
	}
	ret, err := m.UploadMaterial("123.png", bytes.NewReader(picBytes), true, media.TypImage)
	if err != nil {
		log.Panic(err)
	}
	fmt.Printf("%v\n", ret)
	//fmt.Printf("%v", err)
	//
	//fmt.Printf("%v %v %v", mCounter, mList, err)
	//fmt.Printf("%v", ret)

	wx := api.New(&api.Config{
		AppID:        config.AppID,
		AppSecret:    config.AppSecret,
		AppToken:     config.AppToken,
		AesEncodeKey: config.AesEncodeKey,
		Callback:     "/wx",
	})

	wx.EventHandle = func(m *api.MessageReceive) api.MessageReply {
		log.Printf("EventHandler: user: %s: %s\n", m.FromUserName, m.Event)
		switch m.Event {
		case "subscribe":
			log.Printf("subscribe")
			return api.Text{Content: fmt.Sprintf("%s: welcome", m.FromUserName)}
		case "unsubscribe":
			log.Printf("unsubscribe")
		default:
			log.Printf("default")
		}
		return nil
	}

	wx.MessageHandle = func(m *api.MessageReceive) api.MessageReply {
		return api.Articles{
			Items: []api.ArticleItem{
				{
					Title:       "title1",
					Description: "description1",
					PicURL:      "https://i.loli.net/2019/11/17/QKv1kdFmquzsfYr.jpg",
					URL:         "https://dreamer2q.wang/index.php/study/2020/04/59/",
				},
			},
		}
		//return api.Image{
		//	MediaId: "bGL97-3t2bSmRVZowC4qmECPCNzHe8Jure1-H6QERUw",
		//}
		//return api.Text{Content: "Hello"}
	}

	log.Panicln(wx.Run(":80"))
}
