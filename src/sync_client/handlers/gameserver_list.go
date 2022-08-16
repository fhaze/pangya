package handlers

import (
	"encoding/json"
	"net"
	"pangya/src/internal/sync"
	loginserver "pangya/src/login_server"
	syncclient "pangya/src/sync_client"
)

type ServerGameServerList struct {
	svc syncclient.Client
	ls  *loginserver.LoginServer
}

func (sgsl *ServerGameServerList) Action(conn net.Conn, pak []byte) error {
	var res sync.ServerPacketGameServerList
	if err := json.Unmarshal(pak, &res); err != nil {
		return err
	}

	sgsl.ls.GameServers = res.Servers

	return nil
}

func NewServerGameServerList(svc syncclient.Client, ls *loginserver.LoginServer) sync.PacketHandler {
	return &ServerGameServerList{svc: svc, ls: ls}
}
