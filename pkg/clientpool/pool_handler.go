package clientpool

import (
	"log"
	"sync"

	"github.com/MobilityData/gtfs-realtime-bindings/golang/gtfs"
	"github.com/a-ferraro007/improved-train/pkg/types"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// TODO: PoolMap should hold map[string][]*Pool and hold a slice of pools for each subwayLine
// PoolMap Struct
type PoolMap struct {
	Mutex sync.RWMutex
	Map   map[string]*Pool
}

// Pools Map
var Pools PoolMap

// Init Function
func Init() {
	log.Default().Println("Initialize Pools")
	Pools.Map = make(map[string]*Pool)
}

// HandleNewConnection Function
func HandleNewConnection(subwayLine string, stopID string, conn *websocket.Conn) {
	if Pools.Map[subwayLine] == nil {
		log.Println("nil")
		createPool(subwayLine)
		insertIntoPool(subwayLine, stopID, conn)
	} else {
		insertIntoPool(subwayLine, stopID, conn)
	}
}

func createPool(subwayLine string) *Pool {
	Pools.Mutex.Lock()
	defer Pools.Mutex.Unlock()

	pool := newPool(subwayLine)
	Pools.Map[subwayLine] = pool

	go pool.run()
	go pool.fetchData()
	return pool
}

func insertIntoPool(subwayLine string, stopID string, conn *websocket.Conn) {
	Pools.Mutex.Lock()
	defer Pools.Mutex.Unlock()
	pool := Pools.Map[subwayLine]

	client := &Client{
		UUID:       uuid.New(),
		Pool:       pool,
		Conn:       conn,
		Send:       make(chan []*gtfs.TripUpdate_StopTimeUpdate),
		StopID:     stopID,
		SubwayLine: subwayLine,
		Config:     types.Config{StopID: stopID, SubwayLine: subwayLine, Sort: "ascending"},
		Fetching:   false,
	}
	client.SortConfig()

	cache := make([]*gtfs.TripUpdate_StopTimeUpdate, 0)
	cache = pool.CachedStopTimeUpdate[client.Config.SubwayLine]
	pool.Register <- client
	go client.read()
	go client.write(&cache)
	log.Printf("Inserted ClientId: %v in Pool: %v\n", client.UUID, pool.SubwayLine)
}

// DeletePool function
func (p *PoolMap) DeletePool(subwayLine string) {
	p.Mutex.Lock()
	defer p.Mutex.Unlock()

	p.Map[subwayLine].Done <- true
	delete(p.Map, subwayLine)
	log.Printf("Deleted Pool: %v, Pool Map: %v\n", subwayLine, p.Map)
}
