package syncclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"pangya/src/internal/logger"
	"pangya/src/internal/sync"
	"strconv"
	"strings"
)

type Client interface {
	Dial(addr string, port int) error
	Handshake(server string) error
	Read() ([]byte, error)
}

type syncClient struct {
	conn net.Conn
}

func New() Client {
	return &syncClient{}
}

func (svc *syncClient) Dial(addr string, port int) error {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", addr, port))
	if err != nil {
		return err
	}

	logger.Log.Sugar().Infof("connected to sync server %s:%d", addr, port)
	svc.conn = conn

	return nil
}

func (svc *syncClient) Handshake(server string) error {
	addr := strings.Split(svc.conn.LocalAddr().String(), ":")
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
	_, err = svc.conn.Write(buf)
	if err != nil {
		return err
	}
	logger.Log.Sugar().Infof("sucessfully handshaked to sync server")

	return nil
}

func (svc *syncClient) Read() ([]byte, error) {
	buf := make([]byte, 1_024)
	_, err := svc.conn.Read(buf)
	if err != nil {
		return nil, err
	}
	data := bytes.Trim(buf, "\x00")
	return data, nil
}
