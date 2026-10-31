package server

import (
	"github.com/mahdi-1991/CoCo-tunnel/internal/packet/engine/flog"
	"github.com/mahdi-1991/CoCo-tunnel/internal/packet/engine/protocol"
	"github.com/mahdi-1991/CoCo-tunnel/internal/packet/engine/tnet"
)

func (s *Server) handlePing(strm tnet.Strm) {
	p := protocol.Proto{Type: protocol.PPONG}
	if err := p.Write(strm); err != nil {
		flog.Errorf("failed to send pong on stream %d: %v", strm.SID(), err)
	}
}
