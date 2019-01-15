package alerts

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/prometheus/alertmanager/notify"
)

var signaldAddr = os.Getenv("SIGNALD_BIND_ADDR")
var senderNumber = os.Getenv("SENDER_NUMBER")
var receiverGroupId = os.Getenv("RECEIVER_GROUP_ID")

type SignalRequest struct {
	Type             string `json:"type"`
	Username         string `json:"username,omitempty"`
	MessageBody      string `json:"messageBody,omitempty"`
	RecipientGroupID string `json:"recipientGroupId,omitempty"`
}

func Handler(w http.ResponseWriter, req *http.Request) {
	// Verify that request is a POST.
	if req.Method != http.MethodPost {
		log.Printf("Client used unsupported method: %s: %s", req.Method, req.RemoteAddr)
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Read request body.
	alertBytes, err := ioutil.ReadAll(req.Body)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to read request body: %s", err)
		log.Printf(errMsg)
		http.Error(w, errMsg, http.StatusInternalServerError)
		return
	}

	// The WebhookMessage is dependent on alertmanager version. Parse it.
	msg := &notify.WebhookMessage{}
	if err := json.Unmarshal(alertBytes, msg); err != nil {
		log.Printf("Failed to parse webhook message from %s: %s", req.RemoteAddr, err)
		log.Printf("%s", string(alertBytes))
		http.Error(w, "Failed to parse alert", http.StatusBadRequest)
		return
	}
	// log.Print(pretty.Sprint(msg))

	// Handle the webhook message.
	if err := processAlert(msg); err != nil {
		errMsg := fmt.Sprintf("Failed to handle alert: %s", err)
		log.Printf(errMsg)
		http.Error(w, errMsg, http.StatusInternalServerError)
		return
	}
	log.Printf("Completed alert: %s", msg)
	w.WriteHeader(http.StatusOK)
	// Empty response.
}

func processAlert(msg *notify.WebhookMessage) error {
	conn, err := net.Dial("tcp", signaldAddr)
	req := SignalRequest{"send", senderNumber, formatAlert(msg), receiverGroupId}
	b, err := json.Marshal(req)
	log.Print("Sending ", string(b))
	e := json.NewEncoder(conn)
	e.Encode(req)
	return err
}

func formatAlert(msg *notify.WebhookMessage) string {
	return fmt.Sprintf("%s: %s", msg.Data.GroupLabels["alertname"], msg.Data.ExternalURL)
}
