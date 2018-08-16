package updater

import (
	"fmt"
	"os"
	"os/signal"

	log "github.com/sirupsen/logrus"
	"github.com/yosssi/gmq/mqtt"
	"github.com/yosssi/gmq/mqtt/client"
)

// AdaConfig - adafruit configurator
type AdaConfig struct {
	AdafruitHost    string
	AdafruitPort    string
	AdafruitUser    string
	AdafruitToken   string
	AdaFruitTopic   string
	AdaFruitMessage string
}

// ExportStats - push statistic to external graphic backend
func ExportStats(ada AdaConfig) {

	log.Printf("Send stats to AdaFruit, number of users : %s", ada.AdaFruitMessage)

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt, os.Kill)

	cli := client.New(&client.Options{
		// Define the processing of the error handler.
		ErrorHandler: func(err error) {
			fmt.Println(err)
		},
	})
	err := cli.Connect(&client.ConnectOptions{
		Network:      "tcp",
		Address:      ada.AdafruitHost + ":" + ada.AdafruitPort,
		UserName:     []byte(ada.AdafruitUser),
		Password:     []byte(ada.AdafruitToken),
		CleanSession: true,
	})
	if err != nil {
		panic(err)
	}
	log.Println("AdaFruit service connected")
	defer cli.Terminate()

	err = cli.Publish(&client.PublishOptions{
		QoS:       mqtt.QoS0,
		TopicName: []byte(ada.AdaFruitTopic),
		Retain:    true,
		Message:   []byte(ada.AdaFruitMessage),
	})

	if err != nil {
		panic(err)
	}

	<-sigc

	if err := cli.Disconnect(); err != nil {
		panic(err)
	}
}
