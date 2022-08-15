package pangya

import (
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
	ServerName() string
}

type pangyaServer struct {
	hello    func(net.Conn) uint16
	handlers map[uint16]PacketHandler
}

func NewServer(hello func(net.Conn) uint16) Server {
	return &pangyaServer{
		hello:    hello,
		handlers: make(map[uint16]PacketHandler),
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

func (svc *pangyaServer) ServerName() string {
	return "GenericServer"
}

func (svc *pangyaServer) handleConnection(conn net.Conn) {
	defer conn.Close()

	key := svc.hello(conn)

	buf := make([]byte, 1_024)
	l, err := conn.Read(buf)
	if err != nil {
		logger.Log.Sugar().Errorf("error reading from pangya client %s: %s", conn.RemoteAddr().String(), err)
		return
	}

	encryptedData := buf[:l]
	logger.Log.Debug("recived packet", zap.Int("packetLength", len(encryptedData)))

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
