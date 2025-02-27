package tokenreceiver

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/cli/browser"
	"github.com/rs/cors"
)

type ReceiverServer struct {
	server *http.Server
	token  string
	port   string
}

func New() *ReceiverServer {
	return &ReceiverServer{}
}

func getPort(listener net.Listener) string {
	_, port, _ := net.SplitHostPort(listener.Addr().String())
	return port
}

func getURL(target string) string {
	return fmt.Sprintf("http://localhost:3000/account/upctl-login/%s", target)
}

func (s *ReceiverServer) GetLoginURL() string {
	return getURL(s.port)
}

func (s *ReceiverServer) Start() error {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /ping", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})
	mux.HandleFunc("POST /callback", func(w http.ResponseWriter, req *http.Request) {
		token := req.URL.Query().Get("token")
		if token == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		s.token = token
		w.WriteHeader(http.StatusNoContent)
	})

	handler := cors.Default().Handler(mux)
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return fmt.Errorf("failed to create receiver server: %w", err)
	}

	go func() {
		defer listener.Close()
		s.server = &http.Server{
			Handler:           handler,
			ReadHeaderTimeout: time.Second,
		}
		_ = s.server.Serve(listener)
	}()
	s.port = getPort(listener)
	return nil
}

func (s *ReceiverServer) OpenBrowser() error {
	return browser.OpenURL(s.GetLoginURL())
}

func (s *ReceiverServer) Wait(ctx context.Context) (string, error) {
	ticker := time.NewTicker(time.Second * 2)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			_ = s.server.Shutdown(context.TODO())
			return "", ctx.Err()
		case <-ticker.C:
			if s.token != "" {
				_ = s.server.Shutdown(context.TODO())
				return s.token, nil
			}
		}
	}
}
