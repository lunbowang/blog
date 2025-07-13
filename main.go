package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"web_app/controller"
	"web_app/dao/mysql"
	"web_app/dao/redis"
	"web_app/logger"
	"web_app/pkg/snowflake"
	"web_app/routes"
	"web_app/settings"

	"github.com/spf13/viper"

	"go.uber.org/zap"
)

// @title bluebell
// @version 1.0
// @description 简单贴吧
// @termsOfService http://swagger.io/terms/

// @contact.name 1492568061@qq.com
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8081
// @BasePath /api/v1
func main() {
	//1.加载配置
	if err := settings.Init(); err != nil {
		fmt.Printf("init settings failed, err:%v\n", err)
		return
	}
	//2.初始化日志
	if err := logger.Init(settings.Conf.LogConfig, settings.Conf.Mode); err != nil {
		fmt.Printf("init logger failed, err:%v\n", err)
		return
	}
	defer zap.L().Sync()
	zap.L().Debug("logger init success...")
	//3.初始化MYSQL连接
	defer mysql.Close()
	if err := mysql.Init(settings.Conf.MySQLConfig); err != nil {
		fmt.Printf("init mysql failed, err:%v\n", err)
		return
	}
	//4.初始化Redis连接
	defer redis.Close()

	if err := snowflake.Init(settings.Conf.StartTime, settings.Conf.MachineID); err != nil {
		fmt.Printf("init snowflake failed, err:%v\n", err)
		return
	}

	if err := redis.Init(settings.Conf.RedisConfig); err != nil {
		fmt.Printf("init redis failed, err:%v\n", err)
		return
	}

	//初始化gin框架内置的校验器使用的翻译器
	if err := controller.Init("zh"); err != nil {
		fmt.Printf("init validator trans failed,err:%v\n", err)
		return
	}

	//5.注册路由
	r := routes.Setup(settings.Conf.Mode)
	//6.启动服务（优雅关机）
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", viper.GetInt("port")),
		Handler: r,
	}
	go func() {
		//开启一个goroutine启动服务
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	//等待中断信号来优雅地关闭服务器，为关闭服务器操作设置一个5秒地超时
	quit := make(chan os.Signal, 1) //创建一个接收信号地通道
	//kill 默认会发送 syscall.SIGTERM 信号
	//kill -2 发送 syscall.SIGINT 信号，我们常用的Ctrl+C就是触发系统SIGINT信号
	//kill -9 发送 syscall.SIGKILL 信号，但是不能被捕获，所以不需要添加它
	//signal.Notify把收到的syscall.SIGINT或syscall.SIGTERM信号转发给quit
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) //此处不会阻塞
	<-quit                                               //阻塞在此，当接收到上述两个信号时才会往下执行
	zap.L().Info("Shutdown Server ...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	//5秒内优雅关闭服务（将未处理完的请求处理完毕再关闭服务），超过5秒就超时退出
	if err := srv.Shutdown(ctx); err != nil {
		zap.L().Fatal("Server Shutdown:", zap.Error(err))
	}
	zap.L().Info("Server exiting")
}
