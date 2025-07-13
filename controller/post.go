package controller

import (
	"strconv"
	"web_app/logic"
	"web_app/models"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// CreatePostHandler 创建帖子的处理函数
func CreatePostHandler(c *gin.Context) {
	//1. 获取参数及参数的校验
	p := new(models.Post)
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Debug("c.ShouldBindJSON(p) error", zap.Any("err", err))
		zap.L().Error("create post with invalid param")
		ResponseError(c, CodeInvalidParam)
		return
	}
	//从c取到当前发请求的用户ID
	userID, err := getCurrentUserID(c)
	if err != nil {
		ResponseError(c, CodeNeedLogin)
		return
	}
	p.AuthorID = userID
	//2. 创建帖子
	if err := logic.CreatePost(p); err != nil {
		zap.L().Error("logic.CreatePost failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	//3. 返回响应
	ResponseSuccess(c, nil)
}

// GetPostDetailHandler 获取帖子详情的处理函数
func GetPostDetailHandler(c *gin.Context) {
	//1. 获取参数（从URL中获取帖子的id）
	pidStr := c.Param("id")
	pid, err := strconv.ParseInt(pidStr, 10, 64)
	if err != nil {
		zap.L().Error("get post detail with invalid param", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}
	//2. 根据id取出帖子数据（查数据库）
	data, err := logic.GetPostById(pid)
	if err != nil {
		zap.L().Error("logic.GetPostById failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	//3. 返回响应
	ResponseSuccess(c, data)
}

func GetPostListHandler(r *gin.Context) {
	//数据
	page, size, err := GetPostParam(r)
	if err != nil {
		zap.L().Error("GetPostParam err", zap.Error(err))
		ResponseError(r, CodeInvalidParam)
		return
	}
	//业务处理
	date, err := logic.GetPostList(page, size)
	if err != nil {
		zap.L().Error("logic.GetPosts err", zap.Error(err))
		ResponseError(r, CodeServerBusy)
	}
	//返回
	ResponseSuccess(r, date)
}

// GetPostListHandler2 升级版帖子列表接口
// @Summary 升级版帖子列表接口
// @Description 可按社区按时间或分数排序查询帖子列表接口
// @Tags 帖子相关接口
// @Accept application/json
// @Produce application/json
// @Param Authorization header string false "Bearer 用户令牌"
// @Param object query models.ParamPostList false "查询参数"
// @Security ApiKeyAuth
// @Success 200 {object} _ResponsePostList
// @Router /posts2 [get]
// 根据前端传来的参数动态的获取帖子列表
// 按创建时间排序 或者 按照分数排序
// 1.获取请求的query string参数
// 2.去redis查询id列表
// 3.根据id去数据库查询帖子详细信息
func GetPostListHandler2(r *gin.Context) {
	//初始化结构体时指定初始化参数
	p := &models.ParamPostList{
		Page:  1,
		Size:  10,
		Order: models.OrderTime,
	}
	//r.ShouldBind() 根据请求的数据类型选择相应的方法去获取数据
	//r.ShouldBindJSON() 如果请求中携带的是json格式的数据，才能用这个方法获取到数据
	if err := r.ShouldBindQuery(p); err != nil {
		zap.L().Error("GetPostListHandler2 with invalid params", zap.Error(err))
		ResponseError(r, CodeInvalidParam)
		return
	}
	//业务处理
	data, err := logic.GetPostListNew(p)

	if err != nil {
		zap.L().Error("logic.GetPosts err", zap.Error(err))
		ResponseError(r, CodeServerBusy)
	}
	//返回
	ResponseSuccess(r, data)
}

//func GetCommunityPostListHandler(r *gin.Context) {
//	//初始化结构体时指定初始化参数
//	p := &models.ParamCommunityPostList{
//		ParamPostList: models.ParamPostList{
//			Page:  1,
//			Size:  10,
//			Order: models.OrderTime,
//		},
//	}
//	//r.ShouldBind() 根据请求的数据类型选择相应的方法去获取数据
//	//r.ShouldBindJSON() 如果请求中携带的是json格式的数据，才能用这个方法获取到数据
//	if err := r.ShouldBindQuery(p); err != nil {
//		zap.L().Error("GetCommunityPostListHandler with invalid params", zap.Error(err))
//		ResponseError(r, CodeInvalidParam)
//		return
//	}
//	//业务处理
//	date, err := logic.GetCommunityPostList(p)
//	if err != nil {
//		zap.L().Error("logic.GetPosts err", zap.Error(err))
//		ResponseError(r, CodeServerBusy)
//	}
//	//返回
//	ResponseSuccess(r, date)
//}
