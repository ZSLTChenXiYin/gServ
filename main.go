package main

import (
	"context"
	"fmt"
	"gServ/core/gameserv"
	"gServ/core/httpserv"
	"gServ/core/log"
	"gServ/core/repository"
	"gServ/core/tcpserv"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	log.StdInfof("gServ 初始化成功")

	// 创建HTTP服务器
	http_server := httpserv.GetHTTPServer()
	// 创建TCP服务器
	tcp_server := tcpserv.GetTCPListener()

	// 运行服务器
	runHttpServer(http_server)
	log.StdInfof("gServ HTTP 服务启动成功")
	runTcpServer(tcp_server)
	log.StdInfof("gServ TCP 服务启动成功")

	// 运行房间自动清理协程
	room_queue := make(chan struct{}, 1)
	runRoomAutoClean(room_queue)
	log.StdInfof("gServ 游戏房间自动清理协程启动成功")
	// 运行验证码五分钟过期自动清理协程
	captcha_queue := make(chan struct{}, 1)
	runCaptchaAutoClean(captcha_queue)
	log.StdInfof("gServ 验证码自动清理协程启动成功")

	// 创建退出信号通道
	quit := make(chan os.Signal, 1)
	// 等待中断信号
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// 停止服务器
	err := stopHttpServer(http_server)
	if err != nil {
		log.StdErrorf("gServ HTTP 服务停止失败: %v", err)
	}
	log.StdInfof("gServ HTTP 服务停止成功")
	err = stopTcpServer(tcp_server)
	if err != nil {
		log.StdErrorf("gServ TCP 服务停止失败: %v", err)
	}
	log.StdInfof("gServ TCP 服务停止成功")

	// 停止自动清理协程
	stopRoomAutoClean(room_queue)
	log.StdInfof("gServ 游戏房间自动清理协程停止成功")
	stopCaptchaAutoClean(captcha_queue)
	log.StdInfof("gServ 验证码自动清理协程停止成功")
}

func runHttpServer(server *http.Server) {
	// 在 goroutine 中启动服务器
	go func() {
		log.StdDebugf("gServ HTTP 服务监听启动中")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.StdErrorf("gServ HTTP 服务监听失败: %v", err)
		}
	}()
}

func runTcpServer(server net.Listener) {
	go func() {
		log.StdDebugf("gServ TCP 服务监听启动中")
		for {
			conn, err := server.Accept()
			if err != nil {
				log.StdErrorf("gServ TCP 服务连接客户端失败: %v", err)
				continue
			}

			go tcpserv.HandleConnection(conn)
		}
	}()
}

func runRoomAutoClean(room_queue chan struct{}) {
	go func() {
		log.StdDebugf("gServ 游戏房间自动清理协程启动中")

		// 创建一个定时任务，每隔5分钟清理一次房间
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()

		for {
			select {
			case <-room_queue:
				return
			case <-ticker.C:
				// 清理房间
				gameserv.CleanRooms()
			}
		}
	}()
}

func runCaptchaAutoClean(captcha_queue chan struct{}) {
	go func() {
		log.StdDebugf("gServ 验证码自动清理协程启动中")

		// 创建一个定时任务，每隔5分钟清理一次验证码
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()

		for {
			select {
			case <-captcha_queue:
				return
			case <-ticker.C:
				// 清理验证码
				err := repository.DeleteUnusedEmailCaptchas()
				if err != nil {
					log.StdWarnf("删除未使用的验证码失败: %v", err)
				}
			}
		}
	}()
}

func stopHttpServer(server *http.Server) error {
	log.StdDebugf("gServ HTTP 服务停止中")

	// 设置优雅关闭超时时间
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 关闭服务器
	if err := server.Shutdown(ctx); err != nil {
		return fmt.Errorf("HTTP服务关闭失败: %v", err)
	}

	return nil
}

func stopTcpServer(server net.Listener) error {
	log.StdDebugf("gServ TCP 服务停止中")

	// 停止TCP服务器
	if err := server.Close(); err != nil {
		return fmt.Errorf("TCP服务关闭失败: %v", err)
	}

	return nil
}

func stopRoomAutoClean(room_queue chan struct{}) {
	log.StdDebugf("gServ 游戏房间自动清理协程停止中")

	room_queue <- struct{}{}
	close(room_queue)
}

func stopCaptchaAutoClean(captcha_queue chan struct{}) {
	log.StdDebugf("gServ 验证码自动清理协程停止中")

	captcha_queue <- struct{}{}
	close(captcha_queue)
}
