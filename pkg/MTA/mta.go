package mta

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"

	"github.com/a-ferraro007/improved-train/pkg/station"
)

type Stations struct {
	stations []*station.Station
}

type MtaApi struct {
	StationsFile    string
	Feed_URLS       []string
	Stations        map[string]*Stations
	StopsToStations map[string]interface{}
	Routes          map[string]interface{}
	MaxTrains       int16
	MaxMinutes      int16
	ExpiresSeconds  int16
}

func (mta *MtaApi) Init() {}

// _stops, stopsErr := readCsv("../../static_transit/stops.txt")
// if stopsErr != nil {
// 	log.Println(stopsErr)
// 	return
// }
// _transfers, transferErr := readCsv("../../static_transit/transfers.txt")
// if transferErr != nil {
// 	log.Println(transferErr)
// 	return
// }

// stops := make(map[string]map[string]string)
// transfers := make(map[string]map[string]struct{})

// for _, row := range _stops {
// 	if val, ok := row["parent_station"]; ok && val == "" {
// 		// log.Println(row)
// 		continue
// 	}
// 	// log.Println(row["stop_id"])
// 	stops[row["parent_station"]] = map[string]string{
// 		"name": row["stop_name"],
// 		"lat":  row["stop_lat"],
// 		"lon":  row["stop_lon"],
// 	}

// }
// log.Println(stops["101"])
// // log.Printf("id: %s\n stop: %v\n", row["stop_id"], stops[row["stop_id"]])

// for _, transfer := range _transfers {
// 	fromStopID := transfer["from_stop_id"] //.(string)
// 	toStopID := transfer["to_stop_id"]     //.(string)
// 	// if transfer[fromStopID] == transfer[toStopID] {
// 	// 	fromTo[fromStopID] = make(map[string]struct{})
// 	// 	fromTo[fromStopID][toStopID] = struct{}{}
// 	// 	continue
// 	// }
// 	// log.Println(t)

// 	if _, exists := transfers[fromStopID]; !exists {
// 		transfers[fromStopID] = make(map[string]struct{})
// 		transfers[fromStopID][toStopID] = struct{}{}
// 	} else {
// 		transfers[fromStopID][toStopID] = struct{}{}
// 		// log.Printf("from: %s\n to: %s\n", fromStopID, toStopID)
// 		// log.Println(transfers[fromStopID])
// 	}
// }

// log.Println(transfers["F12"])

// for parent_id := range transfers {

// 	// 	if parent_id > min(transfers[parent_id]) {
// 	// 		continue
// 	// 	}

// 	for stopId := range transfers[parent_id] {
// 		log.Printf("\n id: %s\n stop: %v \n parent: %s", stopId, stops[stopId], parent_id)
// 		// log.Printf("from: %s\n to: %s\n", parent_id, transfers[parent_id])
// 		delete(stops, stopId)
// 	}
// }

// csvFile, err := os.Create(mta.StationsFile)
// if err != nil {
// 	log.Fatalf("failed creating file: %s", err)
// }
// defer csvFile.Close()

// row := []string{"stop_id", "name", "lat", "lon", "parent_id"}
// writer := csv.NewWriter(csvFile)
// writer.Write(row)
// // log.Println(transfers)
// for from := range transfers {
// 	// if from > min(transfers[from]) {
// 	// 	continue
// 	// }

// 	// for to := range toSet {
// 	// 	if from == to {
// 	// 		log.Printf("EQALS from: %s, to: %s, toSet: %s", from, to, toSet)
// 	// 		continue
// 	// 	}
// 	// log.Printf("from: %s, to: %s, toSet: %s", from, to, toSet)
// 	// log.Println(stops["725N"])
// 	for stopId := range transfers[from] {
// 		// log.Printf("stopid: %s, from: %s,  toSet: %s, min: %s", stopId, from, toSet, min(transfers[from]))
// 		var stop map[string]string
// 		S := fmt.Sprintf("%sS", stopId)
// 		N := fmt.Sprintf("%sN", stopId)
// 		if len(stops[N]) != 0 {
// 			stop = stops[N]
// 		} else if len(stops[S]) != 0 {
// 			stop = stops[S]
// 		}

// 		if stop != nil {
// 			delete(stops, stopId)
// 			row := []string{stopId, stop["name"], stop["lat"], stop["lon"], from}
// 			writer.Write(row)
// 		}
// 	}

// 	// log.Println(len(stops))

// 	// for stopId := range stops {
// 	// 	row := []string{stopId, stops[stopId]["name"], stops[stopId]["lat"], stops[stopId]["lon"], from}
// 	// 	writer.Write(row)
// 	// }
// 	// }
// }
// writer.Flush()
// mta.InitJson()

func (mta *MtaApi) InitJson() {
	file, err := os.Open(mta.StationsFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	records, err := csv.NewReader(file).ReadAll()
	if err != nil {
		// handle error
		log.Fatal(err)
	}
	for _, record := range records {
		log.Println(record)
	}
}

func readCsv(file string) ([]map[string]string, error) {
	f, err := os.Open(file)
	var data []map[string]string

	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	reader := csv.NewReader(f)

	header, err := reader.Read()
	if err != nil {
		fmt.Println("Error reading header:", err)
		return data, err
	}

	for {
		record, err := reader.Read()
		if err != nil {
			break
		}

		rowMap := make(map[string]string)
		for i, field := range record {
			rowMap[header[i]] = field
		}
		data = append(data, rowMap)
	}
	log.Println(data[0])
	return data, err
}

func min(data map[string]struct{}) string {
	var min string
	for key := range data {
		if len(min) == 0 || min < key {
			min = key
		}
	}
	return min
}

//  map[112:map[A09:{}] 125:map[A24:{}] 127:map[725:{} 902:{} A27:{} R16:{}] 132:map[D19:{} L02:{}] 140:map[R27:{}] 222:map[415:{}] 228:map[A36:{} E01:{} R25:{}] 229:map[418:{} A38:{} M22:{}] 232:map[423:{} R28:{}] 235:map[D24:{} R31:{}] 239:map[S04:{}] 254:map[L26:{}] 414:map[D11:{}] 415:map[222:{}] 418:map[229:{} A38:{} M22:{}] 423:map[232:{} R28:{}] 629:map[B08:{} R11:{}] 630:map[F11:{}] 631:map[723:{} 901:{}] 635:map[L03:{} R20:{}] 637:map[D21:{}] 639:map[M20:{} Q01:{} R23:{}] 640:map[M21:{}] 710:map[G14:{}] 718:map[R09:{}] 719:map[F09:{} G22:{}] 723:map[631:{} 901:{}] 724:map[D16:{}] 725:map[127:{} 902:{} A27:{} R16:{}] 901:map[631:{} 723:{}] 902:map[127:{} 725:{} A27:{} R16:{}] A09:map[112:{}] A12:map[D13:{}] A24:map[125:{}] A27:map[127:{} 725:{} 902:{} R16:{}] A31:map[L01:{}] A32:map[D20:{}] A36:map[228:{} E01:{} R25:{}] A38:map[229:{} 418:{} M22:{}] A41:map[R29:{}] A45:map[S01:{}] A51:map[J27:{} L22:{}] B08:map[629:{} R11:{}] B16:map[N04:{}] D11:map[414:{}] D13:map[A12:{}] D16:map[724:{}] D17:map[R17:{}] D19:map[132:{} L02:{}] D20:map[A32:{}] D21:map[637:{}] D24:map[235:{} R31:{}] E01:map[228:{} A36:{} R25:{}] F09:map[719:{} G22:{}] F11:map[630:{}] F15:map[M18:{}] F23:map[R33:{}] G14:map[710:{}] G22:map[719:{} F09:{}] G29:map[L10:{}] J27:map[A51:{} L22:{}] L01:map[A31:{}] L02:map[132:{} D19:{}] L03:map[635:{} R20:{}] L10:map[G29:{}] L17:map[M08:{}] L22:map[A51:{} J27:{}] L26:map[254:{}] M08:map[L17:{}] M18:map[F15:{}] M20:map[639:{} Q01:{} R23:{}] M21:map[640:{}] M22:map[229:{} 418:{} A38:{}] N04:map[B16:{}] Q01:map[639:{} M20:{} R23:{}] R09:map[718:{}] R11:map[629:{} B08:{}] R16:map[127:{} 725:{} 902:{} A27:{}] R17:map[D17:{}] R20:map[635:{} L03:{}] R23:map[639:{} M20:{} Q01:{}] R25:map[228:{} A36:{} E01:{}] R27:map[140:{}] R28:map[232:{} 423:{}] R29:map[A41:{}] R31:map[235:{} D24:{}] R33:map[F23:{}] S01:map[A45:{}] S04:map[239:{}]]
