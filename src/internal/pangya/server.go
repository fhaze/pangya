package pangya

import (
	"encoding/hex"
	"fmt"
	"net"
	"pangya/src/internal/logger"
	"time"

	"github.com/pangbox/pangcrypt"
	"go.uber.org/zap"
)

type Server interface {
	Listen(port int) error
	AddHandler(id uint16, ph PacketHandler)
	ServerInfo() ServerInfo
}

type ServerInfo struct {
	Type     string `json:"type"`
	Name     string `json:"name,omitempty"`
	IP       string `json:"ip,omitempty"`
	Port     uint16 `json:"port,omitempty"`
	MaxUsers uint32 `json:"maxUsers,omitempty"`
	Flags    uint32 `json:"flags,omitempty"`
	Boosts   uint16 `json:"boosts,omitempty"`
}

type ServerConfig interface {
	OnClientConnect(con net.Conn) uint16
}

type pangyaServer struct {
	handlers map[uint16]PacketHandler
	conf     ServerConfig
}

func NewServer(conf ServerConfig) Server {
	return &pangyaServer{
		handlers: make(map[uint16]PacketHandler),
		conf:     conf,
	}
}

func (svc *pangyaServer) AddHandler(id uint16, ph PacketHandler) {
	svc.handlers[id] = ph
}

func (svc *pangyaServer) Listen(port int) error {
	portStr := fmt.Sprintf(":%d", port)
	tcp, err := net.Listen("tcp", portStr)
	if err != nil {
		return err
	}
	logger.Log.Sugar().Infof("listening on port %s", portStr)
	defer tcp.Close()

	for {
		conn, err := tcp.Accept()
		conn.SetDeadline(time.Now().Add(time.Second * 5))
		if err != nil {
			return err
		}
		logger.Log.Sugar().Infof("accepted connection from %s", conn.RemoteAddr().String())
		go svc.handleConnection(conn)
	}
}

func (svc *pangyaServer) ServerInfo() ServerInfo {
	return ServerInfo{Type: "GenericServer"}
}

func (svc *pangyaServer) handleConnection(conn net.Conn) {
	defer conn.Close()

	key := svc.conf.OnClientConnect(conn)

	buf := make([]byte, 1_024)
	l, err := conn.Read(buf)
	if err != nil {
		logger.Log.Sugar().Errorf("error reading from pangya client %s: %s", conn.RemoteAddr().String(), err)
		return
	}

	encryptedData := buf[:l]

	if len(encryptedData) == 0 {
		logger.Log.Debug("packet size is 0")
		return
	}

	data, err := pangcrypt.ClientDecrypt(encryptedData, byte(key))
	if err != nil {
		logger.Log.Error(
			"could not decrypt client packet",
			zap.Error(err),
		)
	}

	pak, err := PacketFromBytes(data)
	if err != nil {
		logger.Log.Error(
			"invalid pangya packet",
			zap.Error(err),
		)
		return
	}

	logger.Log.Debug(
		"received packet",
		zap.Int("length", len(pak.ToBytes())),
		zap.Int("packetID", int(pak.ID)),
		zap.String("packetPayload", hex.EncodeToString(pak.Payload)),
		zap.String("packetPayloadStr", string(pak.Payload)),
	)

	h, found := svc.handlers[pak.ID]
	if !found {
		logger.Log.Sugar().Warnf("packet %d not implemented", pak.ID)
		return
	}

	logger.Log.Sugar().Debugf("calling action for packet %d", pak.ID)
	if err := h.Action(conn, pak, key); err != nil {
		logger.Log.Error(err.Error())
	}
}
