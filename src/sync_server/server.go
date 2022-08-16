package syncserver

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"pangya/src/internal/logger"
	"pangya/src/internal/pangya"
	"pangya/src/internal/sync"

	"go.uber.org/zap"
)

type Server interface {
	Listen(port int) error
	AddClient(server pangya.ServerInfo, conn net.Conn)
	AddHandler(id string, ph sync.PacketHandler)
	GameServerList() []*Client
}

type Client struct {
	conn net.Conn
	Info pangya.ServerInfo
}

type syncServer struct {
	clients  map[string]map[string]*Client
	handlers map[string]sync.PacketHandler
}

func New() Server {
	return &syncServer{
		clients:  make(map[string]map[string]*Client),
		handlers: make(map[string]sync.PacketHandler),
	}
}

func (svc *syncServer) AddClient(server pangya.ServerInfo, conn net.Conn) {
	if svc.clients[server.Type] == nil {
		svc.clients[server.Type] = make(map[string]*Client)
	}
	svc.clients[server.Type][conn.RemoteAddr().String()] = &Client{
		Info: server,
		conn: conn,
	}
}

func (svc *syncServer) AddHandler(id string, ph sync.PacketHandler) {
	svc.handlers[id] = ph
}

func (svc *syncServer) Listen(port int) error {
	portStr := fmt.Sprintf(":%d", port)
	tcp, err := net.Listen("tcp", portStr)
	if err != nil {
		return err
	}
	logger.Log.Sugar().Infof("listening on port %s", portStr)
	defer tcp.Close()

	for {
		conn, err := tcp.Accept()
		if err != nil {
			return err
		}
		logger.Log.Sugar().Infof("accepted connection from %s", conn.RemoteAddr().String())

		go svc.handleConnection(conn)
	}
}

func (svc *syncServer) GameServerList() []*Client {
	var clients []*Client
	for _, client := range svc.clients["GameServer"] {
		clients = append(clients, client)
	}
	return clients
}

func (svc *syncServer) handleConnection(conn net.Conn) {
	for {
		buf := make([]byte, 1_024)
		l, err := conn.Read(buf)
		if err != nil {
			logger.Log.Sugar().Error(err)
			conn.Close()

			remoteAddr := conn.RemoteAddr().String()
			for serverType := range svc.clients {
				if _, found := svc.clients[serverType][remoteAddr]; found {
					delete(svc.clients[serverType], remoteAddr)
					logger.Log.Sugar().Infof("unregistered %s from %s", serverType, conn.RemoteAddr())
				}
			}
			return
		}

		if l == 0 {
			logger.Log.Debug("packet size is 0")
			continue
		}

		data := bytes.Trim(buf, "\x00")
		logger.Log.Debug("received packet", zap.String("payload", string(data)))

		var base sync.PacketBase
		err = json.Unmarshal(data, &base)
		if err != nil {
			logger.Log.Sugar().Error(err.Error())
			continue
		}

		h, found := svc.handlers[base.ID]
		if !found {
			logger.Log.Sugar().Warnf("packet %s not implemented", base.ID)
			continue
		}

		if err := h.Action(conn, data); err != nil {
			logger.Log.Error(err.Error())
		}
	}
}
