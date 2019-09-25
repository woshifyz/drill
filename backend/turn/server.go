package turn

import (
	"github.com/pion/logging"
	"github.com/pion/turn"
	"net"
)

func createAuthHandler(usersMap map[string]string) turn.AuthHandler {
	return func(username string, srcAddr net.Addr) (string, bool) {
		if password, ok := usersMap[username]; ok {
			return password, true
		}
		return "", false
	}
}

type TurnServer struct {
	User     string
	Password string
	Port     int
	Realm    string
	server   *turn.Server
}

func NewTurnServer(user, password string, port int, realm string) *TurnServer {
	return &TurnServer{
		User:     user,
		Password: password,
		Port:     port,
		Realm:    realm,
	}
}

func (s *TurnServer) Start() error {
	usersMap := map[string]string{
		s.User: s.Password,
	}

	s.server = turn.NewServer(&turn.ServerConfig{
		Realm:         s.Realm,
		AuthHandler:   createAuthHandler(usersMap),
		ListeningPort: s.Port,
		LoggerFactory: logging.NewDefaultLoggerFactory(),
	})

	return s.server.Start()
}

func (s *TurnServer) Stop() error {
	if s.server == nil {
		return nil
	}
	return s.server.Close()
}
