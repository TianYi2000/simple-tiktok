package controller

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/TianYi2000/simple-tiktok/models"
	"github.com/gin-gonic/gin"
)

// MaxVideoNum 每次最多返回的视频流数量
const (
	MaxVideoNum = 30
)

type FeedVideoList struct {
	Videos   []*models.Video `json:"video_list,omitempty"`
	NextTime int64           `json:"next_time,omitempty"`
}

func QueryFeedVideoList(userId int64, latestTime time.Time) (*FeedVideoList, error) {
	return NewQueryFeedVideoListFlow(userId, latestTime).Do()
}

type QueryFeedVideoListFlow struct {
	userId     int64
	latestTime time.Time

	videos   []*models.Video
	nextTime int64

	feedVideo *FeedVideoList
}

func NewQueryFeedVideoListFlow(userId int64, latestTime time.Time) *QueryFeedVideoListFlow {
	return &QueryFeedVideoListFlow{userId: userId, latestTime: latestTime}
}

func (q *QueryFeedVideoListFlow) Do() (*FeedVideoList, error) {
	//所有传入的参数不填也应该给他正常处理
	q.checkNum()

	if err := q.prepareData(); err != nil {
		return nil, err
	}
	if err := q.packData(); err != nil {
		return nil, err
	}
	return q.feedVideo, nil
}

func (q *QueryFeedVideoListFlow) checkNum() {
	//上层通过把userId置零，表示userId不存在或不需要
	if q.userId > 0 {
		//这里说明userId是有效的，可以定制性的做一些登录用户的专属视频推荐
	}

	if q.latestTime.IsZero() {
		q.latestTime = time.Now()
	}
}

func FillVideoListFields(videos *[]*models.Video) (*time.Time, error) {
	size := len(*videos)
	if videos == nil || size == 0 {
		return nil, errors.New("util.FillVideoListFields videos为空")
	}

	latestTime := (*videos)[size-1].CreatedAt //获取最近的投稿时间
	return &latestTime, nil
}

func (q *QueryFeedVideoListFlow) prepareData() error {
	err := models.NewVideoDAO().QueryVideoListByLimitAndTime(MaxVideoNum, q.latestTime, &q.videos)
	if err != nil {
		return err
	}
	latestTime, _ := FillVideoListFields(&q.videos)
	if latestTime != nil {
		q.nextTime = (*latestTime).UnixNano() / 1e6
		return nil
	}
	q.nextTime = time.Now().Unix() / 1e6
	return nil
}

func (q *QueryFeedVideoListFlow) packData() error {
	q.feedVideo = &FeedVideoList{
		Videos:   q.videos,
		NextTime: q.nextTime,
	}
	return nil
}

type FeedResponse struct {
	models.CommonResponse
	*FeedVideoList
}

func Feed(c *gin.Context) {
	p := NewVideoList(c)
	err := p.Get()
	if err != nil {
		p.FeedVideoListError(err.Error())
	}
	return
}

type VideoList struct {
	*gin.Context
}

func NewVideoList(c *gin.Context) *VideoList {
	return &VideoList{Context: c}
}

func (p *VideoList) Get() error {
	rawTimestamp := p.Query("latest_time")
	var latestTime time.Time
	intTime, err := strconv.ParseInt(rawTimestamp, 10, 64)
	if err == nil {
		latestTime = time.Unix(0, intTime*1e6) //注意：前端传来的时间戳是以ms为单位的
	}
	videoList, err := QueryFeedVideoList(0, latestTime)
	if err != nil {
		return err
	}
	p.FeedVideoListOk(videoList)
	return nil
}

// DoHasToken 如果是登录状态，则生成UserId字段
func (p *VideoList) DoHasToken(token string) error {
	return nil
}

func (p *VideoList) FeedVideoListError(msg string) {
	p.JSON(http.StatusOK, FeedResponse{CommonResponse: models.CommonResponse{
		StatusCode: 1,
		StatusMsg:  msg,
	}})
}

func (p *VideoList) FeedVideoListOk(videoList *FeedVideoList) {
	p.JSON(http.StatusOK, FeedResponse{
		CommonResponse: models.CommonResponse{
			StatusCode: 0,
		},
		FeedVideoList: videoList,
	},
	)
}
