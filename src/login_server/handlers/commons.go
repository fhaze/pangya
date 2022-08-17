package handlers

import (
	"pangya/src/internal/pangya"
	loginserver "pangya/src/login_server"
	"pangya/src/models"
)

func packetGameServerList(ls *loginserver.LoginServer) pangya.Packet {
	w := pangya.NewPacket(0x0002)
	w.PutUint8(uint8(len(ls.GameServers)))
	for _, g := range ls.GameServers {
		appendGameServer(&w, g)
	}
	return w
}

func appendNeedSetNickname(w *pangya.Packet) {
	w.PutUint8(216)
	w.PutUint32(0xffffffff)
}

func appendNeedSelectCharacter(w *pangya.Packet) {
	w.PutUint8(217)
}
func appendLoginSuccess(w *pangya.Packet, acc models.Account) {
	w.PutUint8(0x00)                                                                                       // status ok
	w.PutLString(acc.Username)                                                                             // username
	w.PutUint32(uint32(acc.ID))                                                                            // id
	w.PutBytes([]byte{0x00, 0x00, 0x00, 0x00, 0x25, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}) // unknown
	w.PutUint16(0x0000)                                                                                    // nickname?
}

func appendInvalidCredentials(w *pangya.Packet) {
	w.PutUint8(227)      // status
	w.PutUint32(5100143) // invalid credentials
}

func appendGameServer(w *pangya.Packet, g pangya.ServerInfo) {
	w.PutString(g.Name, 40)
	w.PutUint32(uint32(g.ID))
	w.PutUint32(g.MaxUsers)
	w.PutUint32(1)
	w.PutString(g.IP, 18)
	w.PutUint16(g.Port)
	w.PutUint16(0)        // unknown
	w.PutUint16(g.Flags)  // flags
	w.PutUint16(0)        // unknown
	w.PutUint16(0)        // unknown
	w.PutUint16(0)        // unknown
	w.PutUint16(g.Boosts) // boosts
	w.PutUint16(0)        // unknown
	w.PutUint16(0)        // unknown
	w.PutUint16(0)        // unknown
	w.PutUint16(g.Icon)   // character icon
}
