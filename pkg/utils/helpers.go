package utils

import (
	"log"
	"strings"
	"time"

	"github.com/MobilityData/gtfs-realtime-bindings/golang/gtfs"
	"github.com/a-ferraro007/improved-train/pkg/types"
)

// ConvertToTrainSliceAndParse Function
func ConvertToTrainSliceAndParse(stopTimeUpdates []*types.StopTimeUpdate) types.TrainsByDirection {
	trainsByDirection := types.TrainsByDirection{North: make([]*types.Train, 0), South: make([]*types.Train, 0)}
	for _, stopTime := range stopTimeUpdates {
		train := &types.Train{
			StopTimeUpdate: stopTime,
		}

		train.StopTimeUpdate.ProcessStopTimeUpdate()

		if train.StopTimeUpdate.SecondsUntilArrival <= -30 {
			continue
		}

		if train.StopTimeUpdate.SecondsUntilArrival <= 30 {
			train.StopTimeUpdate.IsArriving = true
		}

		direction := strings.ToLower(strings.Split(stopTime.ID, "")[len(strings.Split(stopTime.ID, ""))-1])
		switch direction {
		case "n":
			train.Direction = "N"
			trainsByDirection.North = append(trainsByDirection.North, train)
		case "s":
			train.Direction = "S"
			trainsByDirection.South = append(trainsByDirection.South, train)
		default:
			log.Default().Println("Error: Direction unknown: ", direction)
		}
	}

	return trainsByDirection
}

// ParseTripUpdate Function
func ParseTripUpdate(trip *gtfs.TripDescriptor, gtfsStopTimeUpdate *gtfs.TripUpdate_StopTimeUpdate, ret *types.StopTimeUpdate, stopID string) bool {
	if gtfsStopTimeUpdate != nil && strings.Contains(gtfsStopTimeUpdate.GetStopId(), stopID) {
		log.Default().Println(gtfsStopTimeUpdate.GetStopId())
		ret.ID = gtfsStopTimeUpdate.GetStopId()
		ret.Trip = trip

		departure := gtfsStopTimeUpdate.GetDeparture()
		if departure != nil {
			if departure.Delay != nil {
				ret.DepartureDelay.Delay = departure.GetDelay()
			}
			if departure.Uncertainty != nil {
				ret.DepartureDelay.Uncertainty = departure.GetUncertainty()
			}
			if departure.Time != nil {
				ret.DepartureTime = departure.GetTime()
			}
		}

		arrival := gtfsStopTimeUpdate.GetArrival()
		if arrival != nil {
			if arrival.Delay != nil {
				ret.ArrivalDelay.Delay = arrival.GetDelay()
			}
			if arrival.Uncertainty != nil {
				ret.ArrivalDelay.Uncertainty = arrival.GetUncertainty()
			}
			if arrival.Time != nil {
				ret.ArrivalTime = arrival.GetTime()
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

	// sort.SliceStable(parsed.North, func(i, j int) bool {
	// 	return parsed.North[i].StopTimeUpdate.ArrivalTimeInMinutesWithDelay < parsed.North[j].StopTimeUpdate.ArrivalTimeInMinutesWithDelay
	// })

	// sort.SliceStable(parsed.South, func(i, j int) bool {
	// 	return parsed.South[i].StopTimeUpdate.ArrivalTimeInMinutesWithDelay < parsed.South[j].StopTimeUpdate.ArrivalTimeInMinutesWithDelay
	// })

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
