package realtime

import (
	"drill/config"
	"drill/realtime/space"
	"encoding/json"
	"fmt"
	"time"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"net/http"
)

type WsDoor struct {
	port int
}

func NewWsDoor(port int) *WsDoor {
	return &WsDoor{
		port: port,
	}
}

type JoinRoomRequest struct {
	Nickname string `json:"nickname" msgpack:"nickname"`
	RoomId   string `json:"room_id" msgpack:"room_id"`
}

type BroadcastRequest struct {
	RoomId  string `json:"room_id" msgpack:"room_id"`
	Msg     string `json:"msg" msgpack:"msg"`
	Version int32  `json:"version" msgpack:"version"`
	Full    bool   `json:"full" msgpack:"full"`
	All     string `json:"all" msgpack:"all"`
}

type JoinMessage struct {
	Action         string   `json:"action" msgpack:"action"`
	Content        []string `json:"content" msgpack:"content"`
	Members        []string `json:"members" msgpack:"members"`
	CurrentVersion int32    `json:"current_version" msgpack:"current_version"`
}

type RefreshMemberMessage struct {
	Action  string   `json:"action" msgpack:"action"`
	Members []string `json:"members" msgpack:"members"`
}

type EditorDeltaMessage struct {
	Action         string `json:"action" msgpack:"action"`
	Delta          string `json:"delta" msgpack:"delta"`
	CurrentVersion int32  `json:"current_version" msgpack:"current_version"`
	Version        int32  `json:"version" msgpack:"version"`
	Sender         string `json:"sender" msgpack:"sender"`
}

func DecodeMsg(s []byte, v interface{}) error {
	return json.Unmarshal(s, v)
}

func EncodeMsg(v interface{}) []byte {
	resp, _ := json.Marshal(v)
	return resp
}

//func reverse(numbers []byte) {
//	for i, j := 0, len(numbers)-1; i < j; i, j = i+1, j-1 {
//		numbers[i], numbers[j] = numbers[j], numbers[i]
//	}
//}
//
//func DecodeMsg(s []byte, v interface{}) error {
//	reverse(s)
//	return msgpack.Unmarshal(s, v)
//}
//
//func EncodeMsg(v interface{}) []byte {
//	resp, _ := msgpack.Marshal(v)
//	reverse(resp)
//	return resp
//}

func UpdateRoomInfo(backend RoomBackend, roomInfo *RoomInfo, content string, full bool, all string) *RoomInfo {
	roomInfo.mux.Lock()
	defer roomInfo.mux.Unlock()
	roomInfo.CurrentVersion += 1
	if full {
		roomInfo.Content = []string{all}
	} else {
		roomInfo.Content = append(roomInfo.Content, content)
	}

	backend.Store(roomInfo.RoomId, roomInfo)

	return roomInfo
}

func (door *WsDoor) Start() {
	var backend RoomBackend
	if config.GlobalConfig.RoomBackend == "memory" {
		backend = NewMemoryRoomBackend()
	} else if config.GlobalConfig.RoomBackend == "redis" {
		backend = NewRedisRoomBackend(
			config.GlobalConfig.RedisRoomBackendConfig.Url,
			config.GlobalConfig.RedisRoomBackendConfig.Password,
			config.GlobalConfig.RedisRoomBackendConfig.Db,
			config.GlobalConfig.RedisRoomBackendConfig.Timeout)
	}

	spaceManager := space.NewMemorySpaceManager()

	fmt.Println("start listen websocket on port ", door.port)

	http.ListenAndServe(fmt.Sprintf(":%d", door.port), http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, _, _, err := ws.UpgradeHTTP(r, w)
		if err != nil {
			// handle error
		}
		go func() {
			defer conn.Close()

			fmt.Println(time.Now().String(), "conn: ", conn)
			for {
				msg, _, err := wsutil.ReadClientData(conn)
				if err != nil {
					fmt.Println(conn, err)
					roomIds := spaceManager.GetAllRoomsOfConn(conn)
					spaceManager.CloseConn(conn)

					for _, roomId := range roomIds {
						conns := spaceManager.RoomMembers(roomId)
						members := spaceManager.ConnsToNicknames(conns)
						resp := EncodeMsg(RefreshMemberMessage{"refresh_member", members})
						for _, memberConn := range conns {
							wsutil.WriteServerMessage(memberConn, ws.OpText, resp)
						}
					}
					break
				}

				sMsg := string(msg)
				action, params := sMsg[0], sMsg[2:]
				if action == 'j' {
					var request JoinRoomRequest
					err := DecodeMsg([]byte(params), &request)
					if err != nil {
						continue
					}

					spaceManager.AddConn(conn, request.RoomId, request.Nickname)
					memberConns := spaceManager.RoomMembers(request.RoomId)
					members := spaceManager.ConnsToNicknames(memberConns)

					roomInfo := backend.Load(request.RoomId)

					resp := EncodeMsg(JoinMessage{"join", roomInfo.Content, members, roomInfo.CurrentVersion})
					refreshResp := EncodeMsg(RefreshMemberMessage{"refresh_member", members})
					for _, memberConn := range memberConns {
						if memberConn == conn {
							wsutil.WriteServerMessage(memberConn, ws.OpText, resp)
						} else {
							wsutil.WriteServerMessage(memberConn, ws.OpText, refreshResp)
						}
					}
				} else if action == 'e' {
					var request BroadcastRequest
					err := DecodeMsg([]byte(params), &request)
					if err != nil {
						continue
					}

					conns := spaceManager.RoomMembers(request.RoomId)
					currentUser := spaceManager.GetConnNickname(conn)

					roomInfo := backend.Load(request.RoomId)

					UpdateRoomInfo(backend, roomInfo, request.Msg, request.Full, request.All)

					resp := EncodeMsg(EditorDeltaMessage{"editor_delta", request.Msg, roomInfo.CurrentVersion, request.Version, currentUser})

					for _, memberConn := range conns {
						wsutil.WriteServerMessage(memberConn, ws.OpText, resp)
					}
				}
			}
		}()
	}))
}
