package config

type Config struct {
	PublicWsUrl            string                  `json:"ws_url"`
	HttpPort               int                     `json:"http_port"`
	WsPort                 int                     `json:"ws_port"`
	RoomBackend            string                  `json:"room_backend"`
	RedisRoomBackendConfig *RoomBackendRedisConfig `json:"redis_room_backend,omitempty"`
	EnableTurn             bool                    `json:"enable_turn"`
	TurnConfig             *TurnConfig             `json:"turn_config,omitempty"`
}

type RoomBackendRedisConfig struct {
	Url      string `json:"url"`
	Password string `json:"password,omitempty"`
	Db       int    `json:"db"`
	Timeout  int    `json:"timeout"`
}

type TurnConfig struct {
	User     string `json:"user"`
	Password string `json:"password"`
	UdpPort  int    `json:"udp_port"`
	Realm    string `json:"realm"`
}

var GlobalConfig Config
