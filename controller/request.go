package controller

import (
	"errors"
	"strconv"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

const CtxUserIDKey = "userID"

var ErrorUserNotLogin = errors.New("用户未登录")

// getCurrentUserID 获取当前登录的用户ID
func getCurrentUserID(c *gin.Context) (userID int64, err error) {
	uid, ok := c.Get(CtxUserIDKey)
	if !ok {
		err = ErrorUserNotLogin
		return
	}
	userID, ok = uid.(int64)
	if !ok {
		err = ErrorUserNotLogin
		return
	}
	return
}

func GetPostParam(r *gin.Context) (int64, int64, error) {
	// 获取 query 参数
	pagestr := r.DefaultQuery("page", "1")  // 设置默认值为 1
	sizestr := r.DefaultQuery("size", "10") // 设置默认值为 10

	// 转换 page 参数
	page, err := strconv.ParseInt(pagestr, 10, 64)
	if err != nil {
		zap.L().Error("GetPostParam failed: invalid page value", zap.String("page", pagestr), zap.Error(err))
		return 0, 0, err
	}

	// 转换 size 参数
	size, err := strconv.ParseInt(sizestr, 10, 64)
	if err != nil {
		zap.L().Error("GetPostParam failed: invalid size value", zap.String("size", sizestr), zap.Error(err))
		return 0, 0, err
	}

	// 设置默认值
	if size == 0 {
		size = 10
	}
	if page == 0 {
		page = 1
	}

	return page, size, nil

	//pagestr := r.Query("page")
	//sizestr := r.Query("size")
	//
	//page, err := strconv.ParseInt(pagestr, 10, 64)
	//if err != nil {
	//	zap.L().Error("GetPostParam filed:", zap.Error(err))
	//	return 0, 0, err
	//}
	//size, err1 := strconv.ParseInt(sizestr, 10, 64)
	//if err1 != nil {
	//	zap.L().Error("GetPostParam filed:", zap.Error(err))
	//	return 0, 0, err
	//}
	//if size == 0 {
	//	size = 10
	//}
	//if page == 0 {
	//	page = 1
	//}
	//return page, size, nil
}
