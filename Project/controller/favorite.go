package controller

import (
	"Project/dao"
	"Project/models"
	"Project/utils"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

// FavoriteAction 点赞接口
func FavoriteAction(c *gin.Context) {
	// 获取请求参数
	token := c.DefaultQuery("token", "")
	videoId, _ := strconv.ParseInt(c.DefaultQuery("video_id", "-1"), 10, 64)
	actionType, _ := strconv.ParseInt(c.DefaultQuery("action_type", "-1"), 10, 64)
	// 参数获取失败
	if token == "" || videoId == -1 || actionType == -1 {
		c.JSON(http.StatusOK, models.Response{
			StatusCode: -1,
			StatusMsg:  "failed to obtain parameters!",
		})
		return
	}
	var userId int64
	myClaims, err := utils.ParseToken(token)
	if err != nil { // token 解析失败
		c.JSON(http.StatusOK, models.Response{
			StatusCode: -2,
			StatusMsg:  "token error!",
		})
		return
	} else { // 如果 token 解析成功，获取 userId
		userId, _ = strconv.ParseInt(myClaims.Uid, 10, 64)
	}
	// 数据库操作
	err = dao.FavoriteAction(userId, videoId, actionType)
	if err != nil {
		// 操作失败
		c.JSON(http.StatusOK, models.Response{
			StatusCode: -3,
			StatusMsg:  "update Mysql error!",
		})
		return
	}
	// 操作成功
	c.JSON(http.StatusOK, models.Response{
		StatusCode: 0,
		StatusMsg:  "success!",
	})
}

// FavoriteList 点赞列表接口
func FavoriteList(c *gin.Context) {
	// 获取请求参数
	uid, _ := strconv.ParseInt(c.DefaultQuery("user_id", "-1"), 10, 64)
	token := c.DefaultQuery("token", "")
	// 参数获取失败
	if uid == -1 || token == "" {
		c.JSON(http.StatusOK, models.VideoListResponse{
			Response: models.Response{
				StatusCode: -1,
				StatusMsg:  "failed to obtain parameters!"},
			VideoList: nil,
		})
		return
	}

	// 校验 token
	// token 无效
	if ok, err := utils.CheckToken(token, strconv.FormatInt(uid, 10)); !ok {
		// 记录校验失败 err
		log.Println(err)
		c.JSON(http.StatusOK, models.VideoListResponse{
			Response: models.Response{
				StatusCode: -2,
				StatusMsg:  "token error!"},
			VideoList: nil,
		})
		return
	}

	// 获取视频点赞列表
	videos := dao.GetFavoriteList(uid)

	c.JSON(http.StatusOK, models.VideoListResponse{
		Response: models.Response{
			StatusCode: 0,
			StatusMsg:  "success!"},
		VideoList: videos,
	})
	return
}
