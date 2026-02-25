package utils

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/MobilityData/gtfs-realtime-bindings/golang/gtfs"
	proto "github.com/golang/protobuf/proto"
)

func FetchTransitData(subwayLine string) []*gtfs.TripUpdate_StopTimeUpdate {
	client := &http.Client{}
	stopTimeUpdate := make([]*gtfs.TripUpdate_StopTimeUpdate, 0)

	reqURL := SUBWAY_LINE_REQUEST_URLS[subwayLine]
	req, err := http.NewRequest("GET", reqURL, nil)
	req.Header.Set("x-api-key", MTA_API_KEY)
	if err != nil {
		log.Default().Println("Error fetching transit data: ", err)
		return stopTimeUpdate
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Default().Println(err)
		return stopTimeUpdate
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Default().Println(err)
		return stopTimeUpdate
	}

	return parseStopTimeUpdate(body, stopTimeUpdate)
}

func parseStopTimeUpdate(body []byte, stopTimeUpdates []*gtfs.TripUpdate_StopTimeUpdate) []*gtfs.TripUpdate_StopTimeUpdate {
	feed := gtfs.FeedMessage{}

	err := proto.Unmarshal(body, &feed)
	if err != nil {
		log.Default().Println("Error parsing StopTimeUpdate: ", err)
		return stopTimeUpdates
	}

	for _, entity := range feed.Entity {
		tripUpdate := entity.TripUpdate
		if tripUpdate != nil {
			stopTimeUpdates = append(stopTimeUpdates, tripUpdate.GetStopTimeUpdate()...)
		}
	}
	return stopTimeUpdates
}
