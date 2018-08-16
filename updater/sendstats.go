package updater

import (
	"os"
	"os/signal"

	log "github.com/sirupsen/logrus"
	"github.com/yosssi/gmq/mqtt"
	"github.com/yosssi/gmq/mqtt/client"
)

// AdaConfig - adafruit configurator
type AdaConfig struct {
	AdafruitHost  string
	AdafruitPort  string
	AdafruitUser  string
	AdafruitToken string
	AdaFruitTopic string
}

// ExportStats - push statistic to external graphic backend
func ExportStats(ada AdaConfig, ch <-chan string) {

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt, os.Kill)

	cli := client.New(&client.Options{
		// Define the processing of the error handler.
		ErrorHandler: func(err error) {
			log.Println(err)
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

	for {
		select {
		case message := <-ch:
			log.Printf("Catch the message, %s", message)
			err = cli.Publish(&client.PublishOptions{
				QoS:       mqtt.QoS0,
				TopicName: []byte(ada.AdaFruitTopic),
				Retain:    true,
				Message:   []byte(message),
			})

			if err != nil {
				panic(err)
			}

		case <-sigc:
			log.Println("AdaFruit service disconnected")
			if err := cli.Disconnect(); err != nil {
				panic(err)
			}

		default:
			continue
		}

	}

}
