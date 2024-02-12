package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"
)

type FileSave struct {
	Id    int64 `json:"response_saved"`
	Event Event `json:"events"`
}

type FormResponse struct {
	FormID      string            `json:"form_id"`
	Token       string            `json:"token"`
	SubmittedAt time.Time         `json:"submitted_at"`
	LandedAt    time.Time         `json:"landed_at"`
	Calculated  Calculated        `json:"calculated"`
	Variables   []Variable        `json:"variables"`
	Hidden      map[string]string `json:"hidden"`
	Definition  Definition        `json:"definition"`
	Answers     []Answer          `json:"answers"`
	Ending      Ending            `json:"ending"`
}

type Calculated struct {
	Score int `json:"score"`
}

type Variable struct {
	Key    string `json:"key"`
	Type   string `json:"type"`
	Number int    `json:"number"`
	Text   string `json:"text"`
}

type Field struct {
	ID                      string   `json:"id"`
	Title                   string   `json:"title"`
	Type                    string   `json:"type"`
	Ref                     string   `json:"ref"`
	AllowMultipleSelections bool     `json:"allow_multiple_selections"`
	AllowOtherChoice        bool     `json:"allow_other_choice"`
	Choices                 []Choice `json:"choices"`
}

type Choice struct {
	ID    string `json:"id"`
	Ref   string `json:"ref"`
	Label string `json:"label"`
}

type Definition struct {
	ID      string   `json:"id"`
	Title   string   `json:"title"`
	Fields  []Field  `json:"fields"`
	Endings []Ending `json:"endings"`
}

type Answer struct {
	Type    string  `json:"type"`
	Text    string  `json:"text"`
	Email   string  `json:"email"`
	Date    string  `json:"date"`
	Number  int     `json:"number"`
	Choices Choices `json:"choices"`
	Boolean bool    `json:"boolean"`
	URL     string  `json:"url"`
	Field   Field   `json:"field"`
}

type Choices struct {
	IDs    []string `json:"ids"`
	Labels []string `json:"labels"`
	Refs   []string `json:"refs"`
}

type Ending struct {
	ID         string     `json:"id"`
	Ref        string     `json:"ref"`
	Title      string     `json:"title"`
	Type       string     `json:"type"`
	Properties Properties `json:"properties"`
}

type Properties struct {
	ButtonText string `json:"button_text"`
	ShowButton bool   `json:"show_button"`
	ShareIcons bool   `json:"share_icons"`
	ButtonMode string `json:"button_mode"`
}

type Event struct {
	EventID      string       `json:"event_id"`
	EventType    string       `json:"event_type"`
	FormResponse FormResponse `json:"form_response"`
}

func eventHandler(w http.ResponseWriter, r *http.Request) {
	var event Event
	var count int64

	err := json.NewDecoder(r.Body).Decode(&event)
	if err != nil {
		http.Error(w, "Failed to decode JSON", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(event)

	fData := &FileSave{
		Id:    count,
		Event: event,
	}
	res, err := json.Marshal(fData)
	if err != nil {
		http.Error(w, "Failed to marshal data", http.StatusInternalServerError)
		return
	}
	if err := os.WriteFile("responses.json", res, 0644); err != nil {
		http.Error(w, "Failed to save to file", http.StatusInternalServerError)
		return
	}
	log.Println("Saved to file")
	count += 1
}

func allWebhooks(w http.ResponseWriter, r *http.Request) {
	var all FileSave
	data, err := os.ReadFile("responses.json")
	if err != nil {
		http.Error(w, "Failed to fetch all webhook responses", http.StatusInternalServerError)
		return
	}
	err = json.Unmarshal(data, &all)
	if err != nil {
		http.Error(w, "Failed to load all webhooks", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(all)
}

func main() {
	http.HandleFunc("/", eventHandler)
	http.HandleFunc("/all", allWebhooks)
	log.Println("Server is running at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
