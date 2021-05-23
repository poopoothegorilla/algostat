package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/fatih/color"
)

const (
	defaultHighPercentile = .9
	defaultLowPercentile  = .1
)

// Payload holds data points.
type Payload struct {
	Points []Point `json:"data"`
}

// Point represents a data point.
type Point struct {
	Time  uint64 `json:"time"`
	MVRVZ string `json:"mvrv-zscore"`
}

func main() {
	starttime := flag.Int64("unix_start_time", 0, "unix timestamp to start measurement")
	starttime = flag.Int64("ust", 0, "unix timestamp to start measurement (shorthand)")
	highPercentile := flag.Float64("high_percentile", defaultHighPercentile, "percentile to use as a high threshold")
	highPercentile = flag.Float64("hp", defaultHighPercentile, "percentile to use as a high threshold (shorthand)")
	lowPercentile := flag.Float64("low_percentile", defaultLowPercentile, "percentile to use as a low threshold")
	lowPercentile = flag.Float64("lp", defaultLowPercentile, "percentile to use as a low threshold (shorthand)")
	flag.Parse()

	payload, err := fetchPayload(*starttime)
	if err != nil {
		log.Fatal(err)
	}

	data := make([]float64, len(payload.Points))
	for i := range data {
		v, err := strconv.ParseFloat(payload.Points[i].MVRVZ, 64)
		if err != nil {
			log.Fatal(err)
		}
		data[i] = v
	}

	currentPayload, err := fetchPayload(time.Now().Add(-1 * time.Hour).Unix())
	if err != nil {
		log.Fatal(err)
	}

	high, low := calculateThresholds(data, *highPercentile, *lowPercentile)
	log.Println("High:", high, "Low:", low)
	point := currentPayload.Points[len(currentPayload.Points)-1]
	score, err := strconv.ParseFloat(point.MVRVZ, 64)
	if err != nil {
		log.Fatal(err)
	}

	switch {
	case score > high:
		color.Set(color.FgRed)
	case score < low:
		color.Set(color.FgGreen)
	}

	log.Println("Recent MVRV:", score)
	color.Unset()
}

func calculateThresholds(data []float64, highPercentile, lowPercentile float64) (high, low float64) {
	if !sort.Float64sAreSorted(data) {
		sort.Float64s(data)
	}

	highIDX := int(math.Round(highPercentile * float64(len(data))))
	lowIDX := int(math.Round(lowPercentile * float64(len(data))))

	return data[highIDX], data[lowIDX]
}

func fetchPayload(starttime int64) (Payload, error) {
	url := "https://new.algoexplorerapi.io/v2/ext/stats/general?interval=24h&indicators=mvrv-zscore"
	if starttime > 0 {
		url = fmt.Sprintf("https://new.algoexplorerapi.io/v2/ext/stats/general?time-start=%v&interval=6h&indicators=mvrv-zscore", starttime)
	}

	var payload Payload

	res, err := http.Get(url)
	if err != nil {
		return payload, err
	}
	points, err := io.ReadAll(res.Body)
	if err != nil {
		return payload, err
	}
	res.Body.Close()

	if err := json.Unmarshal(points, &payload); err != nil {
		return payload, err
	}

	return payload, nil
}
