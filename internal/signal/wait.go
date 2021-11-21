package signal

import (
	"os"
	"os/signal"
)

func Wait(signals ...os.Signal) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, signals...)

	<-stop
}
