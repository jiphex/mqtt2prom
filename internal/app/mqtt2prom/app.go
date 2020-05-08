package mqtt2prom

import (
	"encoding/json"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

func handlePublish(c mqtt.Client, msg mqtt.Message) {
	log.Warnf("unhandled publish to topic: %s", msg.Topic())
}

type Metric map[string]float64

func (s *Server) HandleMQTT(c mqtt.Client, msg mqtt.Message) {
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
	if s.metrics == nil {
		s.metrics = make(map[string]*prometheus.GaugeVec)
	}
	for k := range pkt {
		if s.metrics[k] == nil {
			s.metrics[k] = prometheus.NewGaugeVec(
				prometheus.GaugeOpts{
					Namespace: "zigbee2mqtt",
					Subsystem: "sensor",
					Name:      k,
				}, []string{"sensor_id"},
			)
			prometheus.MustRegister(s.metrics[k])
		}
	}
	for k, v := range pkt {
		obs := s.metrics[k]
		obs.WithLabelValues(strings.SplitN(msg.Topic(), "/", 2)[1]).Set(v)
	}
	// if pkt["battery"] != 0 {
	// }
	// log.Infof("msg: %+v", string(msg.Payload()))
}

var (
	// reg = prometheus.NewRegistry()
	d = prometheus.NewDesc("xxx", "xhelp", nil, prometheus.Labels{})
)

type Server struct {
	metrics map[string]*prometheus.GaugeVec
}

func RunMain() {
	log.SetLevel(log.DebugLevel)
	log.Info("mqtt2prom")
	mopts := mqtt.NewClientOptions().AddBroker("tcp://10.100.100.210:1883").SetClientID("mqtt2prom-xx")
	mopts.SetKeepAlive(2 * time.Second)
	mopts.SetDefaultPublishHandler(handlePublish)
	mopts.SetPingTimeout(5 * time.Second)
	c := mqtt.NewClient(mopts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		log.WithError(token.Error()).Fatal("unable to connect to mqtt")
	}
	log.Info("connected to broker")
	s := &Server{}
	if token := c.Subscribe("zigbee2mqtt/#", 1, s.HandleMQTT); token.Wait() && token.Error() != nil {
		log.WithError(token.Error()).Fatalf("unable to subscribe")
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
