package video

import (
	"errors"
	"time"

	"github.com/TianYi2000/simple-tiktok/models"
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

func FillVideoListFields(userId int64, videos *[]*models.Video) (*time.Time, error) {
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
	//如果用户为登录状态，则更新该视频是否被该用户点赞的状态
	latestTime, _ := FillVideoListFields(q.userId, &q.videos) //不是致命错误，不返回

	//准备好时间戳
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
