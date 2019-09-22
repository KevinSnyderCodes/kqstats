package client

import (
	"github.com/gorilla/websocket"
	"github.com/killer-queen-stats/kqstats/internal/pkg/util"
	"github.com/sirupsen/logrus"
)

// Orchestrator connects to the websocket and takes a list of the different channels to broadcast on
// Orchestrator steps: receive message -> parse message (transform to JSON) -> broadcast message
type Orchestrator struct {
	conn     *websocket.Conn
	StopChan chan interface{}
	// TODO: Channel []*Plugin
	// TODO: Parser *MessageParser
}

// NewOrchestrator returns an error
func NewOrchestrator(info *util.ConnectionInfo) (*Orchestrator, error) {
	conn, err := util.Connect(info)

	if err != nil {
		return nil, err
	}
	orchestrator := &Orchestrator{
		conn:     conn,
		StopChan: make(chan interface{}, 1),
	}
	return orchestrator, nil
}

// ReadMessage reads messages from the websocket
func (o *Orchestrator) ReadMessage() {
	for {
		select {
		case <-o.StopChan:
			logrus.Infof("Shutting down")
			return
		default:
			_, message, err := o.conn.ReadMessage()
			if err != nil {
				// Close the connection and break
				o.conn.Close()
				o.Stop()
			}
			// TODO: Move this to stdout plugin
			logrus.Infof("%v", string(message))
		}
	}
}

// Stop will kill the orchestrator
func (o *Orchestrator) Stop() {
	o.StopChan <- true
}
