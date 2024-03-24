package rlpa

import (
	"encoding/hex"
	"errors"
	"io"
	"log/slog"
	"math/rand/v2"
	"net"
	"time"
)

type Server interface {
	Listen(address string) error
	Shutdown() error
}

type server struct {
	listener *net.TCPListener
	manager  Manager
}

func NewServer(manager Manager) Server {
	return &server{manager: manager}
}

func (s *server) Listen(address string) error {
	var err error
	tcpAddr, _ := net.ResolveTCPAddr("tcp", address)
	s.listener, err = net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		return err
	}
	slog.Info("rLPA server is running on", "address", address)

	for {
		conn, err := s.listener.AcceptTCP()
		if err != nil {
			if err == io.EOF || errors.Is(err, net.ErrClosed) {
				return nil
			}
			conn.Close()
			return err
		}

		conn.SetKeepAlive(true)
		conn.SetKeepAlivePeriod(30 * time.Second)
		go s.handleConn(conn)
	}
}

func (s *server) handleConn(tcpConn *net.TCPConn) {
	id := s.id()
	conn := NewConn(id, tcpConn)
	s.manager.Add(id, conn)
	// TODO: only accept connects from estk.me cards
	slog.Info("new connection from", "id", id)
	defer conn.Close()
	defer s.manager.Remove(id)

	for {
		tag, data, err := conn.Read()
		if err != nil {
			if err == io.EOF || errors.Is(err, net.ErrClosed) {
				return
			}
			slog.Error("error reading from connection", "error", err)
			continue
		}

		// Some workaround, should only be called once
		switch tag {
		case TagManagement:
			s.manager.HandleCallback(CallbackTypeConnSetType, id, string(ConnTypeManagement))
			break
		case TagDownloadProfile:
			s.manager.HandleCallback(CallbackTypeConnSetType, id, string(ConnTypeDownloadProfile))
			break
		case TagProcessNotification:
			s.manager.HandleCallback(CallbackTypeConnSetType, id, string(ConnTypeProcessNotification))
			break
		default:
			break
		}

		if tag == TagClose {
			slog.Info("client closed connection", "id", id)
			return
		}
		if tag == TagAPDU {
			slog.Info("received data from", "id", id, "tag", tag, "data", hex.EncodeToString(data))
		} else {
			slog.Info("received data from", "id", id, "tag", tag, "data", string(data))
		}
		go conn.Dispatch(tag, data)
	}
}

func (s *server) id() string {
	seeds := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")
	id := make([]rune, 6)
	for i := range id {
		id[i] = seeds[rand.IntN(len(seeds))]
	}
	if _, err := s.manager.Get(string(id)); err.Error() == ErrConnNotFound {
		return string(id)
	} else {
		return s.id()
	}
}

func (s *server) Shutdown() error {
	for _, conn := range s.manager.GetAll() {
		conn.Close()
		s.manager.Remove(conn.Id)
	}
	return s.listener.Close()
}
