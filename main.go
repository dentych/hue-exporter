package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"os"
)

var (
	gaugeVecLightOn = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "hue",
		Name:      "light_on",
	}, []string{"id", "name"})
	gaugeVecBrightness = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "hue",
		Name:      "brightness",
	}, []string{"id", "name"})
)

var clientID string

func main() {
	clientID = os.Getenv("HUE_CLIENT_ID")
	if clientID == "" {
		log.Fatalf("HUE_CLIENT_ID is missing")
	}
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	http.DefaultClient.Transport = transport

	http.Handle("/", retrieveMiddleware(promhttp.Handler()))
	log.Println("Listening on :9000")
	err := http.ListenAndServe(":9000", nil)
	if err != nil {
		log.Printf("Listen and serve err: %s\n", err)
	}
}

func retrieveMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		defer next.ServeHTTP(writer, request)
		lights, err := RetrieveLightStatuses()
		if err != nil {
			log.Printf("Failed to retrieve light statuses: %s\n", err)
			return
		}

		for _, light := range lights {
			var onValue, brightness float64
			if light.On {
				onValue = 1
				brightness = light.Brightness
			}
			gaugeVecLightOn.With(prometheus.Labels{"id": light.ID, "name": light.Name}).Set(onValue)
			gaugeVecBrightness.With(prometheus.Labels{"id": light.ID, "name": light.Name}).Set(brightness)
		}
	})
}

func RetrieveLightStatuses() ([]Light, error) {
	log.Println("Retrieving light statuses!")
	req, err := http.NewRequest(http.MethodGet, "https://192.168.1.200/clip/v2/resource/light", nil)
	if err != nil {
		return nil, fmt.Errorf("new request err: %w", err)
	}
	req.Header.Set("hue-application-key", clientID)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send http request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("non-successful status code: %d", resp.StatusCode)
	}

	lightJson := LightJson{}
	err = json.NewDecoder(resp.Body).Decode(&lightJson)
	if err != nil {
		return nil, fmt.Errorf("json decode error: %w", err)
	}

	var lights []Light
	for _, v := range lightJson.Data {
		light := Light{
			ID:         v.ID,
			Name:       v.Metadata.Name,
			Brightness: v.Dimming.Brightness,
			On:         v.On.On,
			Type:       v.Type,
		}
		lights = append(lights, light)
	}
	return lights, nil
}

type DimmingJson struct {
	Brightness float64 `json:"brightness"`
}

type MetadataJson struct {
	Name string `json:"name"`
}

type OnJson struct {
	On bool `json:"on"`
}

type DataJson struct {
	ID       string       `json:"id"`
	Dimming  DimmingJson  `json:"dimming"`
	Metadata MetadataJson `json:"metadata"`
	On       OnJson       `json:"on"`
	Type     string       `json:"type"`
}

type LightJson struct {
	Data []DataJson `json:"data"`
}

type Light struct {
	ID         string
	Name       string
	Brightness float64
	On         bool
	Type       string
}

func (l LightJson) String() string {
	output := ""
	for _, d := range l.Data {
		output += fmt.Sprintf("Light: Name '%s', ID '%s', On: '%t'\n", d.Metadata.Name, d.ID, d.On.On)
	}
	return output
}
