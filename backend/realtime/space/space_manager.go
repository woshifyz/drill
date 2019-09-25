package space

import "net"

type MemorySpaceManager struct {
	connToRoomIdsMap map[net.Conn]map[string]bool
	roomIdToConnsMap map[string]map[net.Conn]bool
	connToNickname   map[net.Conn]string
}

func NewMemorySpaceManager() *MemorySpaceManager {
	return &MemorySpaceManager{
		make(map[net.Conn]map[string]bool),
		make(map[string]map[net.Conn]bool),
		make(map[net.Conn]string),
	}
}

func (manager *MemorySpaceManager) AddConn(conn net.Conn, roomId string, nickname string) {
	existActiveIds, ok := manager.connToRoomIdsMap[conn]
	if !ok {
		existActiveIds = map[string]bool{roomId: true}
	} else {
		existActiveIds[roomId] = true
	}
	manager.connToRoomIdsMap[conn] = existActiveIds

	conns, ok := manager.roomIdToConnsMap[roomId]
	if !ok {
		conns = map[net.Conn]bool{conn: true}
	} else {
		conns[conn] = true
	}
	manager.roomIdToConnsMap[roomId] = conns

	manager.connToNickname[conn] = nickname
}

func (manager *MemorySpaceManager) RoomMembers(roomId string) ([]net.Conn) {
	res, ok := manager.roomIdToConnsMap[roomId]
	if !ok {
		return []net.Conn{}
	}
	conns := make([]net.Conn, 0, len(res))
	for k := range res {
		conns = append(conns, k)
	}
	return conns
}

func (manager *MemorySpaceManager) CloseConn(conn net.Conn) {
	existRoomIds, ok := manager.connToRoomIdsMap[conn]
	if ok {
		for roomId, _ := range existRoomIds {
			if r, ok := manager.roomIdToConnsMap[roomId]; ok {
				delete(r, conn)
			}
		}
	}
	delete(manager.connToRoomIdsMap, conn)
	delete(manager.connToNickname, conn)
}

func (manager *MemorySpaceManager) GetAllRoomsOfConn(conn net.Conn) []string {
	existRoomIds, ok := manager.connToRoomIdsMap[conn]
	if !ok {
		return []string{}
	}
	roomIds := make([]string, 0, len(existRoomIds))
	for k := range existRoomIds {
		roomIds = append(roomIds, k)
	}
	return roomIds
}

func (manager *MemorySpaceManager) ConnsToNicknames(conns []net.Conn) []string {
	nicknames := make([]string, 0, len(conns))
	for _, conn := range conns {
		name, ok := manager.connToNickname[conn]
		if ok {
			nicknames = append(nicknames, name)
		}
	}
	return nicknames
}

func (manager *MemorySpaceManager) GetConnNickname(conn net.Conn) string {
	x, ok := manager.connToNickname[conn]
	if !ok {
		return "noname"
	}
	return x
}
