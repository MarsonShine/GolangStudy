package logcore

import (
	"sync"

	"go.uber.org/zap/zapcore"
)

type ServerWriteSync struct {
	sync.Mutex
	ws zapcore.WriteSyncer
}

func (s *ServerWriteSync) Write(bs []byte) (int, error) {
	s.Lock()
	n, err := s.ws.Write(bs)
	s.Unlock()
	return n, err
}

func (s *ServerWriteSync) Sync() error {
	s.Lock()
	err := s.ws.Sync()
	s.Unlock()
	return err
}
