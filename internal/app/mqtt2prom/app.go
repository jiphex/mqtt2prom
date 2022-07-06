package mqtt2prom

import (
	"encoding/json"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/jiphex/mqtt2prom/internal/pkg/glow"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

func handlePublish(c mqtt.Client, msg mqtt.Message) {
	log.Warnf("unhandled publish to topic: %s", msg.Topic())
}

type Metric map[string]float64

func (s *Server) HandleZigbee2MQTT(c mqtt.Client, msg mqtt.Message) {
	// log.Warn(msg.Topic())
	if strings.HasPrefix(msg.Topic(), "zigbee2mqtt/bridge") {
		log.Trace("skipping bridge message")
		return
	}
	log.Infof("msg to topic: %s - %s", msg.Topic(), string(msg.Payload()))
	var pkt Metric
	err := json.Unmarshal(msg.Payload(), &pkt)
	if err != nil {
		log.WithError(err).Error("unable to unmarshal json")
		return
	}
	log.Infof("msg: %+v", pkt)
	if s.z2metrics == nil {
		s.z2metrics = make(map[string]*prometheus.GaugeVec)
	}
	for k := range pkt {
		if s.z2metrics[k] == nil {
			s.z2metrics[k] = prometheus.NewGaugeVec(
				prometheus.GaugeOpts{
					Namespace: "zigbee2mqtt",
					Subsystem: "sensor",
					Name:      k,
				}, []string{"sensor_id"},
			)
			prometheus.MustRegister(s.z2metrics[k])
		}
	}
	for k, v := range pkt {
		obs := s.z2metrics[k]
		obs.WithLabelValues(strings.SplitN(msg.Topic(), "/", 2)[1]).Set(v)
	}
	// if pkt["battery"] != 0 {
	// }
	// log.Infof("msg: %+v", string(msg.Payload()))
}

func (s *Server) HandleGlow(c mqtt.Client, msg mqtt.Message) {
	log.WithField("topic", msg.Topic()).Debug("glow packet from device")
	parsedTopic := glow.TopicPattern.FindStringSubmatch(msg.Topic())
	switch parsedTopic[2] {
	case "SENSOR":
		switch parsedTopic[3] {
		case "electricitymeter":
			var m glow.Electricitymeter
			if err := json.Unmarshal(msg.Payload(), &m); err != nil {
				log.WithError(err).Error("unable to parse electrictymeter packet")
			}
			log.WithField("m", m).Debug("METER")
			if s.mpanMetrics == nil {
				s.mpanMetrics = make(map[string]*prometheus.GaugeVec)
			}
			if s.mpanMetrics[m.Reading.Energy.Import.Mpan] == nil {
				s.mpanMetrics[m.Reading.Energy.Import.Mpan] = prometheus.NewGaugeVec(prometheus.GaugeOpts{
					Namespace: "glow",
					Subsystem: "meter_reading",
					Name:      "power_kw",
					ConstLabels: prometheus.Labels{
						"unit": "kW",
					},
				}, []string{"mpan"})

				prometheus.MustRegister(s.mpanMetrics[m.Reading.Energy.Import.Mpan])
			}
			s.mpanMetrics[m.Reading.Energy.Import.Mpan].WithLabelValues(m.Reading.Energy.Import.Mpan).Set(m.Reading.Power.Value)
		}
	case "STATE":
		log.Debug("glow state packet")
	}
}

var (
	// reg = prometheus.NewRegistry()
	d = prometheus.NewDesc("xxx", "xhelp", nil, prometheus.Labels{})
)

type Server struct {
	z2metrics   map[string]*prometheus.GaugeVec
	mpanMetrics map[string]*prometheus.GaugeVec
}

func RunMain() {
	log.SetLevel(log.DebugLevel)
	log.Info("mqtt2prom")
	mopts := mqtt.NewClientOptions().AddBroker("tcp://10.100.100.210:1883").SetClientID("mqtt2prom-xxpp")
	mopts.SetKeepAlive(2 * time.Second)
	mopts.SetDefaultPublishHandler(handlePublish)
	mopts.SetPingTimeout(5 * time.Second)
	c := mqtt.NewClient(mopts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		log.WithError(token.Error()).Fatal("unable to connect to mqtt")
	}
	log.Info("connected to broker")
	s := &Server{}
	if token := c.Subscribe("zigbee2mqtt/#", 1, s.HandleZigbee2MQTT); token.Wait() && token.Error() != nil {
		log.WithError(token.Error()).Fatalf("unable to subscribe to Zigbee2MQTT")
	}
	if token := c.Subscribe("glow/#", 1, s.HandleGlow); token.Wait() && token.Error() != nil {
		log.WithError(token.Error()).Fatalf("unable to subscribe to Glow")
	}
	intChan := make(chan os.Signal, 1)
	signal.Notify(intChan, os.Interrupt)
	http.Handle("/metrics", promhttp.Handler())
	go http.ListenAndServe(":2115", nil)
	done := false
	for !done {
		select {
		case <-intChan:
			log.Info("goodbye")
			done = true
		}
	}
}
