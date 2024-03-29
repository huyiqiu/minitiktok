package controller

import (
	"minitiktok/service"
	"minitiktok/utils"
	"minitiktok/utils/logger"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type FeedResponce struct {
	StatusCode int         `json:"status_code"`
	Msg        string      `json:"status_msg"`
	NextTime   int64       `json:"next_time,omitempty"`
	VideoList  interface{} `json:"video_list"`
}

func QueryVideoInfo(lastTime string, userToken string) *FeedResponce {
	videoInfo, err := service.QueryFeedFlow(lastTime, userToken)
	if err != nil {
		return &FeedResponce{
			StatusCode: -1,
			Msg:        err.Error(),
		}
	}
	return &FeedResponce{
		StatusCode: 0,
		Msg:        "success",
		NextTime:   time.Now().Unix(),
		VideoList:  videoInfo,
	}
}

func QueryPublishInfo(token string, userId int) *FeedResponce {
	videoInfo, err := service.QueryPublishFlow(userId, token)
	if err != nil {
		return &FeedResponce{
			StatusCode: -1,
			Msg:        err.Error(),
		}
	}
	return &FeedResponce{
		StatusCode: 0,
		Msg:        "success",
		NextTime:   time.Now().Unix(),
		VideoList:  videoInfo,
	}
}

func Feed(c *gin.Context) {
	lastTime := c.DefaultQuery("latest_time", strconv.FormatInt(time.Now().Unix(), 10))
	userToken := c.DefaultQuery("token", "null")
	if userToken != "null"{
		_, err := utils.ValidateToken(userToken)
		if err != nil {
			c.JSON(http.StatusOK, &ErrorResponce{
				StatusCode: 1,
				StatusMsg:  "token过期,请重新登录",
			})
		}
	}
	c.JSON(http.StatusOK, QueryVideoInfo(lastTime, userToken))
}

func PublishFlow(c *gin.Context) {
	userId, err := strconv.Atoi(c.Query("user_id"))
	if err != nil {
		logger.Warn(err)
	}
	userToken := c.DefaultQuery("token", "null")
	if len(userToken) == 0 {
		userToken = "null"
	}
	if userToken != "null"{
		_, err := utils.ValidateToken(userToken)
		if err != nil {
			c.JSON(http.StatusOK, &ErrorResponce{
				StatusCode: 1,
				StatusMsg:  "token过期,请重新登录",
			})
		}
	}
	c.JSON(http.StatusOK, QueryPublishInfo(userToken, userId))
}
