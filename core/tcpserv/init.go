package tcpserv

import (
	"fmt"
	"gServ/core/config"
	"net"
)

var (
	tcp_listener net.Listener

	player_manager = NewPlayerManager()
)

func Init() error {
	var err error
	tcp_listener, err = net.Listen("tcp", fmt.Sprintf(":%d", config.GetConfig().Server.TCPPort))
	if err != nil {
		return err
	}

	return nil
}

func GetTCPListener() net.Listener {
	return tcp_listener
}
