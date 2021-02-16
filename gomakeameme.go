package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/sfreiberg/gotwilio"
)

type message struct {
	Content string
	To      string
	Medium  string
}

//SendWhatsapp - Sends Whatsapp text message
func SendWhatsapp(m message, credentials *gotwilio.Twilio) (*gotwilio.SmsResponse, error) {
	resp, _, err := credentials.SendWhatsApp("+14155238886", m.To, m.Content, "", "")
	if err != nil {
		return nil, err
	}
	return resp, nil
}

//UseTextWhatsapp Sets message text and forwards
func UseTextWhatsapp(text string) {
	newSms := message{
		To:      "+phnoe",
		Content: text,
		Medium:  "whatsapp",
	}
	twilio := gotwilio.NewTwilioClient(os.Getenv("TwilioAccountSid"), os.Getenv("TwilioAuthToken"))
	SendWhatsapp(newSms, twilio)
	fmt.Println("Sent!")
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
	fmt.Println(r.Form.Get("name"))

	text := r.Form.Get("message")

	UseTextWhatsapp(text)

}

func main() {
	http.HandleFunc("/", ExampleHandler)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}

	newSms := message{
		To:      "+34695848183",
		Content: "This is still sandbox",
		Medium:  "whatsapp",
	}
	twilio := gotwilio.NewTwilioClient(os.Getenv("TwilioAccountSid"), os.Getenv("TwilioAuthToken"))
	SendWhatsapp(newSms, twilio)
}
