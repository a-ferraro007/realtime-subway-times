package utils

import (
	"log"
	"sort"
	"strings"
	"time"

	"github.com/MobilityData/gtfs-realtime-bindings/golang/gtfs"
	"github.com/a-ferraro007/improved-train/pkg/types"
)

// ConvertToTrainSliceAndParse Function
func ConvertToTrainSliceAndParse(stopTimeUpdates []*types.StopTimeUpdate) types.TrainsByDirection {
	trainsByDirection := types.TrainsByDirection{North: make([]*types.Train, 0), South: make([]*types.Train, 0)}
	for _, trip := range stopTimeUpdates {
		train := &types.Train{}
		train.StopTimeUpdate = trip
		if train.StopTimeUpdate.ArrivalTime == nil {
			log.Default().Println("Nil ArrivalTime: ", train.StopTimeUpdate.ArrivalTime)
			continue
		}
		train.StopTimeUpdate.AddDelay()
		train.StopTimeUpdate.ConvertArrivalNoDelay()
		train.StopTimeUpdate.ConvertArrivalWithDelay()
		train.StopTimeUpdate.ConvertTimeToMinutesNoDelay()
		train.StopTimeUpdate.ConvertTimeToMinutesWithDelay()
		//train.Train.ConvertDeparture()

		if train.StopTimeUpdate.TimeInMinutes < 0 {
			log.Default().Println("Negative TimeInMinute: ", train.StopTimeUpdate.TimeInMinutes)
			continue
		}

		idSplit := strings.Split(trip.ID, "")
		direction := strings.ToLower(idSplit[len(idSplit)-1])
		switch direction {
		case "n":
			train.DirectionV2 = "N"
			trainsByDirection.North = append(trainsByDirection.North, train)
		case "s":
			train.DirectionV2 = "S"
			trainsByDirection.South = append(trainsByDirection.South, train)
		default:
			log.Default().Println("Error: Direction unknown: ", direction)
		}
	}

	return trainsByDirection
}

// ParseTripUpdate Function
func ParseTripUpdate(gtfsStopTimeUpdate *gtfs.TripUpdate_StopTimeUpdate, ret *types.StopTimeUpdate, stopID string) bool {
	if gtfsStopTimeUpdate != nil && strings.Contains(gtfsStopTimeUpdate.GetStopId(), stopID) {
		ret.ID = gtfsStopTimeUpdate.GetStopId()

		if gtfsStopTimeUpdate.GetDeparture() != nil {
			ret.DepartureTime = gtfsStopTimeUpdate.GetDeparture().Time
			ret.GtfsDeparture = gtfsStopTimeUpdate.GetDeparture()
		}

		if gtfsStopTimeUpdate.GetArrival() != nil {
			ret.ArrivalTime = gtfsStopTimeUpdate.GetArrival().Time
			if gtfsStopTimeUpdate.GetArrival().Delay != nil {
				ret.Delay = *gtfsStopTimeUpdate.GetArrival().Delay
			}
		}
		return true
	}

	return false
}

func ReturnLimit(trainsByDirection types.TrainsByDirection, limit int) types.TrainsByDirection {
	if limit == 0 || limit > len(trainsByDirection.South) || limit > len(trainsByDirection.North) {
		return trainsByDirection
	}
	return types.TrainsByDirection{
		North: trainsByDirection.North[0:limit],
		South: trainsByDirection.South[0:limit],
	}
}

// DefaultSort Function
func DefaultSort(parsed types.TrainsByDirection) types.TrainsByDirection {
	log.Println("Default sort")

	sort.SliceStable(parsed.North, func(i, j int) bool {
		return parsed.North[i].StopTimeUpdate.TimeInMinutes < parsed.North[j].StopTimeUpdate.TimeInMinutes
	})

	sort.SliceStable(parsed.South, func(i, j int) bool {
		return parsed.South[i].StopTimeUpdate.TimeInMinutes < parsed.South[j].StopTimeUpdate.TimeInMinutes
	})

	return parsed
}

// DescendingSort Function
func DescendingSort(parsed types.TrainsByDirection) types.TrainsByDirection {
	log.Println("Descending sort", time.Now())
	return parsed
}

// TestGen Function
func TestGen(parsed types.TrainsByDirection) types.TrainsByDirection {
	return parsed
}
