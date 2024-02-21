package bot;

import (
	"bytes"
	"encoding/json"
	"net/http"
)

// SendWebhookMessage send message to discord
func SendWebhookMessage(url string, msg string) {
	data := map[string]string{"content": msg}
	jsonData, _ := json.Marshal(data)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	client.Do(req)
}