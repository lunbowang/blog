package redis

import (
	"errors"
	"math"
	"strconv"
	"time"

	"github.com/go-redis/redis"
)

/* 投票的几种情况
direction=1时，有两种情况：
	1.之前没有投过票，现在投赞成票
	2.之前投反对票，现在改投赞成票
direction=0时，有两种情况:
	1.之前投过赞成票，现在要取消投票
	2.之前透过反对票，现在要取消投票
direction=-1时，有两种情况：
	1.之前没有投过票，现在投反对票
	2.之前投赞成票，现在改投反对票

投票的限制：
每个贴子自发表之日一个星期之内允许用户投票，超过一个星期就不允许再投票了。
	1.到期之后将redis中保存的赞成票数及反对票数存到mysql表中
	2.到期之后删除那个 KeyPostVotedZSetPF
*/

const (
	oneWeekInSeconds = 7 * 24 * 3600
	scorePerVote     = 432 //每一票值多少分
)

var (
	ErrVoteTimeExpire = errors.New("投票时间已过")
	ErrVoteRepeated   = errors.New("不允许重复投票")
)

func CreatePost(postID, communityID int64) error {

	pipeline := client.TxPipeline()
	//贴子时间
	pipeline.ZAdd(getRedisKey(KeyPostTimeZSet), redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: postID,
	})

	//帖子分数
	pipeline.ZAdd(getRedisKey(KeyPostScoreZSet), redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: postID,
	})
	//把帖子id加到社区的set
	cKey := getRedisKey(KeyCommunitySetPF + strconv.Itoa(int(communityID)))
	pipeline.SAdd(cKey, postID)
	_, err := pipeline.Exec()
	return err
}

func VoteForPost(userID, postID string, value float64) error {
	//1.判断投票限制
	//去redis取贴子发布时间
	postTime := client.ZScore(getRedisKey(KeyPostTimeZSet), postID).Val()
	if float64(time.Now().Unix())-postTime > oneWeekInSeconds {
		return ErrVoteTimeExpire
	}
	//2.更新贴子的分数
	//先查看当前用户给当前贴子的投票记录
	ov := client.ZScore(getRedisKey(KeyPostVotedZSetPF+postID), userID).Val()
	//更新：如果这一次投票的值和以前一致，就提示不允许重复投票
	if value == ov {
		return ErrVoteRepeated
	}
	var op float64
	if value > ov {
		op = 1
	} else {
		op = -1
	}
	diff := math.Abs(ov - value) //计算两次投票的差值
	pipeline := client.TxPipeline()
	pipeline.ZIncrBy(getRedisKey(KeyPostScoreZSet), op*diff*scorePerVote, postID)
	//3.记录用户为该帖子投票的数据
	if value != 0 {
		client.ZRem(getRedisKey(KeyPostVotedZSetPF+postID), postID)
	} else {
		client.ZAdd(getRedisKey(KeyPostVotedZSetPF+postID), redis.Z{
			Score:  value, //赞成票还是反对票
			Member: userID,
		})
	}
	_, err := pipeline.Exec()
	return err
}
