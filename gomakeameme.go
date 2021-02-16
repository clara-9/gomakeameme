package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"encoding/json"
	"io/ioutil"

	"github.com/sfreiberg/gotwilio"
)

type message struct {
	Content  string
	MediaURL []string
	To       string
	Medium   string
}

type meme struct {
	Success bool
	Data    map[string]interface{}
}

//SendWhatsapp - Sends Whatsapp text message
func SendWhatsapp(m message, credentials *gotwilio.Twilio) (*gotwilio.SmsResponse, error) {
	resp, _, err := credentials.SendWhatsAppMedia("+14155238886", m.To, m.Content, m.MediaURL, "", "")
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func getMeme(text0, text1 string) string {
	postBody, _ := json.Marshal(map[string]string{})
	responseBody := bytes.NewBuffer(postBody)

	url := fmt.Sprintf("https://api.imgflip.com/caption_image?template_id=102156234&text0=%s&text1=%s&username=clarara&password=%s",
		text0, text1, os.Getenv("IMGFLIP_PASSWORD"))
	resp, err := http.Post(url, "application/json", responseBody)
	//Handle Error
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}
	defer resp.Body.Close()
	//Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	sb := string(body)
	log.Printf(sb)

	var m meme
	err = json.Unmarshal([]byte(sb), &m)

	fmt.Println(m.Data["url"].(string))

	memeUrl := m.Data["url"].(string)

	return memeUrl
}

func parseText(text string) (string, string) {
	text = strings.TrimSpace(text)
	half := len(text) / 2
	last := strings.LastIndex(text[:half], " ")
	if last == -1 {
		return "", text
	}
	return text[:last], text[last:]
}

func ExampleHandler(w http.ResponseWriter, r *http.Request) {

	// Double check it's a post request being made
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "invalid_http_method")
		return
	}

	// Must call ParseForm() before working with data
	r.ParseForm()

	// Log all data. Form is a map[]
	log.Println(r.Form)

	// Print the data back. We can use Form.Get() or Form["name"][0]

	text := r.Form.Get("Body")

	to := r.Form.Get("From")[len(r.Form.Get("From"))-11:]
	fmt.Println(to)

	text0, text1 := parseText(text)

	memeURL := []string{getMeme(text0, text1)}

	newSms := message{
		To:       to,
		Content:  "",
		MediaURL: memeURL,
		Medium:   "whatsapp",
	}
	twilio := gotwilio.NewTwilioClient(os.Getenv("TwilioAccountSid"), os.Getenv("TwilioAuthToken"))
	SendWhatsapp(newSms, twilio)
}

func main() {
	http.HandleFunc("/", ExampleHandler)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
