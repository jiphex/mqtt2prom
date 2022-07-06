package glow

import "time"

type Electricitymeter struct {
	Reading struct {
		Timestamp time.Time `json:"timestamp"`
		Energy    struct {
			Export struct {
				Cumulative float64 `json:"cumulative"`
				Units      string  `json:"units"`
			} `json:"export"`
			Import struct {
				Cumulative float64 `json:"cumulative"`
				Day        float64 `json:"day"`
				Week       float64 `json:"week"`
				Month      float64 `json:"month"`
				Units      string  `json:"units"`
				Mpan       string  `json:"mpan"`
				Supplier   string  `json:"supplier"`
				Price      struct {
					Unitrate       float64 `json:"unitrate"`
					Standingcharge float64 `json:"standingcharge"`
				} `json:"price"`
			} `json:"import"`
		} `json:"energy"`
		Power struct {
			Value float64 `json:"value"`
			Units string  `json:"units"`
		} `json:"power"`
	} `json:"electricitymeter"`
}

type State struct {
	Ethmac string `json:"ethmac"`
	Eui    string `json:"eui"`
	Han    struct {
		Lqi    int64  `json:"lqi"`
		Rssi   int64  `json:"rssi"`
		Status string `json:"status"`
	} `json:"han"`
	Hardware     string `json:"hardware"`
	Smetsversion string `json:"smetsversion"`
	Software     string `json:"software"`
	Timestamp    string `json:"timestamp"`
	Wifistamac   string `json:"wifistamac"`
	Zigbee       string `json:"zigbee"`
}
