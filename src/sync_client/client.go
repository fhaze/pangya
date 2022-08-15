package syncclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"pangya/src/internal/logger"
	"pangya/src/internal/pangya"
	"pangya/src/internal/sync"
	"strconv"
	"strings"

	"go.uber.org/zap"
)

type Client interface {
	Dial(addr string, port int) error
	AddHandler(id string, ph sync.PacketHandler)
}

type syncClient struct {
	handlers map[string]sync.PacketHandler
	srv      pangya.Server
}

func New(srv pangya.Server) Client {
	return &syncClient{
		handlers: make(map[string]sync.PacketHandler),
		srv:      srv,
	}
}

func (svc *syncClient) AddHandler(id string, ph sync.PacketHandler) {
	svc.handlers[id] = ph
}

func (svc *syncClient) Dial(addr string, port int) error {
	logger.Log.Sugar().Infof("trying to connect to sync server %s:%d", addr, port)

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", addr, port))
	if err != nil {
		return err
	}
	logger.Log.Sugar().Infof("connected to sync server %s:%d", addr, port)

	svc.handshake(svc.srv.ServerName(), conn)
	go svc.handleConnection(conn)

	return nil
}

func (svc *syncClient) handshake(server string, conn net.Conn) error {
	addr := strings.Split(conn.LocalAddr().String(), ":")
	port, err := strconv.Atoi(addr[1])
	if err != nil {
		return err
	}

	req := sync.ClientPacketHandshake{
		PacketBase: sync.PacketBase{
			ID: sync.PacketHandshake,
		},
		Server: server,
		IP:     addr[0],
		Port:   port,
	}

	var buf []byte
	buf, err = json.Marshal(req)
	if err != nil {
		return err
	}

	logger.Log.Sugar().Infof("trying to handshake to sync server")
	_, err = conn.Write(buf)
	if err != nil {
		return err
	}

	return nil
}

func (svc *syncClient) handleConnection(conn net.Conn) {
	for {
		buf := make([]byte, 1_024)
		l, err := conn.Read(buf)
		if err != nil {
			logger.Log.Sugar().Error(err)
			conn.Close()
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
