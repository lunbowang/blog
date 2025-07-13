package redis

const (
	Prefix             = "bluebell:"
	KeyPostTimeZSet    = "post:time"   //zset;贴子及发贴时间
	KeyPostScoreZSet   = "post:score"  //zset;贴子及投票的分数
	KeyPostVotedZSetPF = "post:voted:" //zset;记录用户及投票类型;参数是post_id

	KeyCommunitySetPF = "community:" //set;保存每个分区下帖子的id
)

// 给redis key加上前缀
func getRedisKey(key string) string {
	return Prefix + key
}
