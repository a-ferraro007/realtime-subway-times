package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/a-ferraro007/improved-train/pkg/clientpool"
	"github.com/a-ferraro007/improved-train/pkg/types"
	"github.com/a-ferraro007/improved-train/pkg/utils"
	"github.com/gorilla/websocket"
)

// var Pools clientpool.PoolMap
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func main() {
	log.Println("Train Time Server v0.3.0")
	clientpool.Init()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		(w).Header().Set("Access-Control-Allow-Origin", "*")
	})

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Default().Println("Error upgrading http connection: ", err)
			return
		}

		subwayLine := r.URL.Query()["subwayLine"]
		stopID := r.URL.Query()["stopID"]

		if len(stopID) == 0 {
			log.Default().Println("Missing stopId")
			return
		}

		if len(subwayLine) == 0 {
			log.Default().Println("Missing subwayLine")
			return
		}

		clientpool.HandleNewConnection(subwayLine[0], stopID[0], conn)
	})

	http.HandleFunc("/transit", func(w http.ResponseWriter, r *http.Request) {
		(w).Header().Set("Access-Control-Allow-Origin", "*")

		stopID := r.URL.Query()["stopID"][0]
		subwayLine := r.URL.Query()["subwayLine"][0]
		stopTimeUpdates := make([]*types.StopTimeUpdate, 0)
		if len(stopID) == 0 {
			log.Default().Println("Missing stopId")
			return
		}

		if len(subwayLine) == 0 {
			log.Default().Println("Missing subwayLine")
			return
		}

		data := utils.FetchTransitData(subwayLine)
		stopTimeUpdate := types.StopTimeUpdate{}
		for _, tripUpdate := range data {
			if utils.ParseTripUpdate(tripUpdate, &stopTimeUpdate, stopID) {
				stopTimeUpdates = append(stopTimeUpdates, &stopTimeUpdate)
			}
		}

		if len(stopTimeUpdates) <= 0 {
			json, _ := json.Marshal("empty")
			w.Header().Set("Content-Type", "application/json")
			w.Write(json)
		}

		trainsByDirection := utils.DefaultSort(utils.ConvertToTrainSliceAndParse(stopTimeUpdates))
		m := clientpool.Message{Message: types.NextTrain{TrainsByDirection: trainsByDirection}}

		json, _ := json.Marshal(m.Message)
		w.Header().Set("Content-Type", "application/json")
		w.Write(json)

	})

	// http.HandleFunc("/stations", func(w http.ResponseWriter, r *http.Request) {
	// 	log.Println(w, "Stations Endpoint")
	// 	(w).Header().Set("Access-Control-Allow-Origin", "*")
	// 	json, _ := json.Marshal(stations)
	// 	w.Header().Set("Content-Type", "application/json")
	// 	w.Write(json)
	// })

	log.Println("Server Running On Port :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
