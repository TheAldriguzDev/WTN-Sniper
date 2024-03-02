package notifier

import (
	console_logger "WTN-Sniper/src/packages/logger"
	"bytes"
	"encoding/json"
	"fmt"

	"net/http"
)

var l = console_logger.New("Notifier", true)

type Notifier struct {
	*Queue
}

func New(max_retries int) *Notifier {
	queue := CreateQueue(max_retries)
	return &Notifier{
		Queue: queue,
	}
}

func SendWebhook(discord_data Discord_Params) error {
	// Check if the webhook url exist
	if discord_data.Webhook_url == "" {
		return fmt.Errorf("Expected discord webhook, received null")
	}

	// Convert the message to a json
	jsonData, err := json.Marshal(discord_data.Webhook_data)

	if err != nil {
		return err
	}

	// Create the request
	req, err := http.NewRequest("POST", discord_data.Webhook_url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	// Set request header
	req.Header.Set("Content-Type", "application/json")

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 && resp.StatusCode != 204 {
		return fmt.Errorf("Error while loggin the webhook: [%v]", resp.StatusCode)
	}

	defer resp.Body.Close()
	return nil
}
