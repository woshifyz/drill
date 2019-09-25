package realtime

import (
	"drill/helper/lru"
	"encoding/json"
	"github.com/go-redis/redis"
	"sync"
	"time"
)

type RoomInfo struct {
	RoomId         string
	Content        []string
	CurrentVersion int32
	mux            sync.Mutex
}

func NewRoomInfo(roomId string) *RoomInfo {
	return &RoomInfo{
		roomId,
		make([]string, 0),
		0,
		sync.Mutex{},
	}
}

type RoomBackend interface {
	Store(roomId string, info *RoomInfo)
	Load(roomId string) *RoomInfo
}

type MemoryRoomBackend struct {
	localCache *lru.Cache
}

func (b *MemoryRoomBackend) Store(roomId string, info *RoomInfo) {
	b.localCache.Add(roomId, info)
}

func (b *MemoryRoomBackend) Load(roomId string) *RoomInfo {
	info, ok := b.localCache.Get(roomId)
	if ok {
		return info.(*RoomInfo)
	}

	return NewRoomInfo(roomId)
}

func NewMemoryRoomBackend() *MemoryRoomBackend {
	cache := lru.New(1000)
	return &MemoryRoomBackend{
		localCache: cache,
	}
}

type RedisRoomBackend struct {
	client     *redis.Client
	timeout    int
	localCache *lru.Cache
}

func NewRedisRoomBackend(addr, pwd string, db int, timeout int) *RedisRoomBackend {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pwd, // no password set
		DB:       db,  // use default DB
	})
	cache := lru.New(1000)
	return &RedisRoomBackend{
		client:     client,
		timeout:    timeout,
		localCache: cache,
	}
}

func (b *RedisRoomBackend) Store(roomId string, info *RoomInfo) {
	val, _ := json.Marshal(info)
	if roomId[len(roomId)-1] == 'P' {
		b.client.Set(roomId, val, 0)
	} else {
		b.client.Set(roomId, val, time.Duration(b.timeout)*time.Second)
	}
	b.localCache.Add(roomId, info)
}

func (b *RedisRoomBackend) Load(roomId string) *RoomInfo {
	info, ok := b.localCache.Get(roomId)
	if ok {
		return info.(*RoomInfo)
	}

	val, err := b.client.Get(roomId).Result()
	if err != nil {
		return NewRoomInfo(roomId)
	}
	var remoteInfo RoomInfo
	err = json.Unmarshal([]byte(val), &remoteInfo)
	if err != nil {
		return NewRoomInfo(roomId)
	}
	remoteInfo.mux = sync.Mutex{}
	return &remoteInfo
}
