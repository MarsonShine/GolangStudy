package main

type Server struct{}

func (s *Server) Call(peerId int, msg string, args interface{}, reply interface{}) error {
	return nil
}
