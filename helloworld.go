package main

import (
	"fmt"
	"os"

   "io/ioutil"
   "log"
   "net/http"
   "bytes"
   "encoding/json"

	"github.com/dghubble/go-twitter/twitter"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

func getStatus(statusNum int64, client *twitter.Client){
   // status show
	statusShowParams := &twitter.StatusShowParams{}
	tweet, _, _ := client.Statuses.Show(statusNum, statusShowParams)
   getTextFromTweet(tweet)
}

func getRichStatus(statusNum int64, client *twitter.Client){
   // HTML status show
	statusOembedParams := &twitter.StatusOEmbedParams{ID: statusNum, MaxWidth: 500}
	oembed, _, _ := client.Statuses.OEmbed(statusOembedParams)
	fmt.Printf("OEMBED TWEET:\n%+v\n", oembed)
}

func getUserTimeline (screenName string, client *twitter.Client){
   // user timeline
	userTimelineParams := &twitter.UserTimelineParams{ScreenName: screenName, Count: 2}
	tweets, _, _ := client.Timelines.UserTimeline(userTimelineParams)
	fmt.Printf("USER TIMELINE:\n%+v\n", tweets)
}
func getTextFromTweet(tweet *twitter.Tweet){
   fmt.Println(tweet.Text)
}

func spongeBobText(text string) (string, string) {
   text0 := "hello"
   text1 := "bye"
   fmt.Println(text0,text1)
   return text0, text1
}

func getMeme(text0, text1 string){
   postBody, _ := json.Marshal(map[string]string{})
   responseBody := bytes.NewBuffer(postBody)

   url := fmt.Sprintf("https://api.imgflip.com/caption_image?template_id=102156234&text0=%s&text1=%s&username=clarara&password=%s",
            text0, text1, os.Getenv("IMGFLIP_PASSWORD"))
   resp, err := http.Post(url,"application/json", responseBody)
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
}

func main() {
	
   consumerKey:= os.Getenv("TWITTER_CONSUMER_KEY")
   consumerSecret :=  os.Getenv("TWITTER_CONSUMER_SECRET")
	
	// oauth2 configures a client that uses app credentials to keep a fresh token
	config := &clientcredentials.Config{
		ClientID:     consumerKey,
		ClientSecret: consumerSecret,
		TokenURL:     "https://api.twitter.com/oauth2/token",
	}
	// http.Client will automatically authorize Requests
	httpClient := config.Client(oauth2.NoContext)

	// Twitter client
	client := twitter.NewClient(httpClient)

	// status show
	getStatus(1353403903690354688, client)

	// oEmbed status
	getRichStatus(1353403903690354688, client)

   // user timeline
   getUserTimeline("hacktheurban", client)

   //spongebobText
   text0,text1 := spongeBobText("hello")

   getMeme(text0, text1)
}