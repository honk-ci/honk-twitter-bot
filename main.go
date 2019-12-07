package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"log"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/ChimeraCoder/anaconda"
)

func main() {
	var flagConfigFile string
	flag.StringVar(&flagConfigFile, "config", "config-honk.json", "Configuration file for the honk bot.")
	flag.Parse()

	err := loadConfig(flagConfigFile)
	if err != nil {
		log.Fatalf("Error starting the job: %v", err)
	}

	log.Println("***starting job***")

	api := anaconda.NewTwitterApiWithCredentials(Config.TwitterAccessToken, Config.TwitterAccessSecret, Config.TwitterConsumerKey, Config.TwitterConsumerSecretKey)

	initialSearch := url.Values{}
	initialSearch.Set("count", Config.TwitterSearchCounts)
	searchResult, err := api.GetSearch("/honk", initialSearch)
	if err != nil {
		log.Fatalf("Error searching tweets: %v", err)
	}

	log.Printf("Got %d results, checking...\n", searchResult.Metadata.Count)
	processHonk(api, searchResult.Statuses)

	log.Println("***Job done***")
}

func checkHonkReply(twitterAPI *anaconda.TwitterApi, tweet anaconda.Tweet) bool {
	searchReplyParams := url.Values{}
	searchReplyParams.Set("to", fmt.Sprintf("@%s", tweet.User.ScreenName))
	searchReplyParams.Set("count", Config.TwitterSearchCounts)
	searchResultReply, _ := twitterAPI.GetSearch("", searchReplyParams)
	for _, tweetReply := range searchResultReply.Statuses {
		if tweetReply.InReplyToStatusIdStr == tweet.IdStr {
			if tweetReply.User.ScreenName == "honk_bot" {
				log.Printf("already replied to this tweet ID = %s, skipping...\n", tweetReply.IdStr)
				return true
			}
		}
	}
	return false
}

func processHonk(twitterAPI *anaconda.TwitterApi, tweetsStatuses []anaconda.Tweet) {
	commandMatch := regexp.MustCompile(`(?mi)^/(?:honk)(?: +(.+?))?\s*$`)
	for _, tweet := range tweetsStatuses {
		tweetTime, _ := tweet.CreatedAtTime()
		log.Printf("Checking Tweet from @%s ID = %s Text = %s TweetTime = %s\n", tweet.User.ScreenName, tweet.IdStr, tweet.Text, tweetTime.UTC())
		if commandMatch.MatchString(tweet.Text) {
			log.Println("Tweet matched honk")

			if strings.Contains(tweet.Text, "RT ") {
				log.Println("This is a RT dont reply to not flood")
				continue
			}

			if checkHonkReply(twitterAPI, tweet) {
				continue
			}

			var goose []byte
			goose = getGoose()
			if goose == nil {
				goose = getDefaultGoose()
			}

			replyParams := url.Values{}
			mediaResponse, err := twitterAPI.UploadMedia(base64.StdEncoding.EncodeToString(goose))
			if err == nil {
				replyParams.Set("media_ids", mediaResponse.MediaIDString)
			}

			replyParams.Set("in_reply_to_status_id", tweet.IdStr)
			replyParams.Set("auto_populate_reply_metadata", "true")
			replyParams.Set("display_coordinates", "false")
			msg := fmt.Sprintf("Honk the Planet @%s", tweet.User.ScreenName)
			result, err := twitterAPI.PostTweet(msg, replyParams)
			if err != nil {
				log.Printf("Error while posting the tweet. Err=%s\n", err.Error())
				continue
			}
			// to avoid getting rate from twitter in case there are too much replies
			time.Sleep(1 * time.Second)
			log.Printf("Tweet posted. TweetID = %s\n", result.IdStr)
		}
	}
}
