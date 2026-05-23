package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

func askAI(text string) string {
	body := map[string]interface{}{
		"model": "gpt-4.1-mini",
		"input": "Kamu adalah admin WhatsApp agency desain. Jawab singkat, ramah, dan bantu customer. Pesan customer: " + text,
	}

	b, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", "https://api.openai.com/v1/responses", bytes.NewBuffer(b))
	req.Header.Set("Authorization", "Bearer "+os.Getenv("OPENAI_API_KEY"))
	req.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "Maaf kak, AI sedang gangguan."
	}
	defer res.Body.Close()

	data, _ := io.ReadAll(res.Body)

	var result map[string]interface{}
	json.Unmarshal(data, &result)

	if output, ok := result["output"].([]interface{}); ok && len(output) > 0 {
		content := output[0].(map[string]interface{})["content"].([]interface{})
		return content[0].(map[string]interface{})["text"].(string)
	}

	return "Siap kak, boleh jelaskan request desainnya?"
}

func sendWA(to, text string) {
	url := "https://graph.facebook.com/v20.0/" + os.Getenv("WA_PHONE_ID") + "/messages"

	payload := map[string]interface{}{
		"messaging_product": "whatsapp",
		"to": to,
		"type": "text",
		"text": map[string]string{"body": text},
	}

	b, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(b))
	req.Header.Set("Authorization", "Bearer "+os.Getenv("WA_TOKEN"))
	req.Header.Set("Content-Type", "application/json")
	http.DefaultClient.Do(req)
}

func webhook(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		if r.URL.Query().Get("hub.verify_token") == os.Getenv("VERIFY_TOKEN") {
			fmt.Fprint(w, r.URL.Query().Get("hub.challenge"))
			return
		}
		http.Error(w, "Forbidden", 403)
		return
	}

	var data map[string]interface{}
	json.NewDecoder(r.Body).Decode(&data)

	defer func() {
		if recover() != nil {
			fmt.Fprint(w, "ok")
		}
	}()

	entry := data["entry"].([]interface{})[0].(map[string]interface{})
	changes := entry["changes"].([]interface{})[0].(map[string]interface{})
	value := changes["value"].(map[string]interface{})
	messages := value["messages"].([]interface{})
	msg := messages[0].(map[string]interface{})

	from := msg["from"].(string)
	textObj := msg["text"].(map[string]interface{})
	text := textObj["body"].(string)

	reply := askAI(text)
	sendWA(from, reply)

	fmt.Fprint(w, "ok")
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "10000"
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "WhatsApp AI Bot Online")
	})

	http.HandleFunc("/webhook", webhook)

	fmt.Println("Server running on port", port)
	http.ListenAndServe(":"+port, nil)
}
