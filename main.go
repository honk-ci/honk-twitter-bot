package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/url"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"syscall"
	"time"

	"github.com/ChimeraCoder/anaconda"
)

var (
	commandHonkMatch  = regexp.MustCompile(`(?mi)^*/(?:honk)(?: +(.+?))?\s*$`)
	commandMeowMatch  = regexp.MustCompile(`(?mi)^*/(?:meow)(?: +(.+?))?\s*$`)
	commandPonyMatch  = regexp.MustCompile(`(?mi)^*/(?:pony)(?: +(.+?))?\s*$`)
	commandWoofMatch  = regexp.MustCompile(`(?mi)^*/(?:woof)(?: +(.+?))?\s*$`)
	commandOinkMatch  = regexp.MustCompile(`(?mi)^*/(?:oink)(?: +(.+?))?\s*$`)
	commandQuackMatch = regexp.MustCompile(`(?mi)^*/(?:quack)(?: +(.+?))?\s*$`)
	commandMooMatch   = regexp.MustCompile(`(?mi)^*/(?:moo)(?: +(.+?))?\s*$`)
	commandBaaMatch   = regexp.MustCompile(`(?mi)^*/(?:baa)(?: +(.+?))?\s*$`)
)

func main() {
	rand.Seed(time.Now().UnixNano())

	var flagConfigFile string
	flag.StringVar(&flagConfigFile, "config", "config-honk.json", "Configuration file for the honk bot.")
	flag.Parse()

	err := loadConfig(flagConfigFile)
	if err != nil {
		log.Fatalf("Error starting the job: %v", err)
	}

	log.Println("***starting Honk Bot***")

	api := anaconda.NewTwitterApiWithCredentials(Config.TwitterAccessToken, Config.TwitterAccessSecret, Config.TwitterConsumerKey, Config.TwitterConsumerSecretKey)

	streamValues := url.Values{}
	streamValues.Set("track", "/honk,/meow,/pony,/woof,/oink,/quack,/moo,/baa")
	streamValues.Set("stall_warnings", "true")
	log.Println("Starting Honk Stream...")
	s := api.PublicStreamFilter(streamValues)

	go func() {
		for t := range s.C {
			switch v := t.(type) {
			case anaconda.Tweet:
				log.Printf("Got one message from @%s", v.User.ScreenName)
				processHonk(api, v)
			}
		}
	}()

	// Wait for SIGINT and SIGTERM
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	log.Println(<-ch)

	fmt.Println("Stopping Stream...")
	s.Stop()

}

func checkHonkReply(twitterAPI *anaconda.TwitterApi, tweet anaconda.Tweet) bool {
	searchReplyParams := url.Values{}
	searchReplyParams.Set("to", fmt.Sprintf("@%s", tweet.User.ScreenName))
	searchReplyParams.Set("count", Config.TwitterSearchCounts)
	searchResultReply, err := twitterAPI.GetSearch("", searchReplyParams)
	if err != nil {
		log.Printf("Error getting the search: %v", err.Error())
		return false
	}
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

func processHonk(twitterAPI *anaconda.TwitterApi, tweet anaconda.Tweet) {
	tweetTime, _ := tweet.CreatedAtTime()
	log.Printf("Checking Tweet from @%s ID = %s Text = %s TweetTime = %s\n", tweet.User.ScreenName, tweet.IdStr, tweet.Text, tweetTime.UTC())
	if commandHonkMatch.MatchString(tweet.Text) || commandMeowMatch.MatchString(tweet.Text) ||
		commandPonyMatch.MatchString(tweet.Text) || commandWoofMatch.MatchString(tweet.Text) ||
		commandOinkMatch.MatchString(tweet.Text) || commandQuackMatch.MatchString(tweet.Text) ||
		commandMooMatch.MatchString(tweet.Text) || commandBaaMatch.MatchString(tweet.Text) {
		go func() {
			n := randomInt(30, 120)
			time.Sleep(time.Duration(n) * time.Second)
			_, err := twitterAPI.Favorite(tweet.Id)
			if err != nil {
				log.Printf("Error while trying to favorite the tweet. Err=%s\n", err.Error())
			}
		}()
	}

	if checkHonkReply(twitterAPI, tweet) {
		return
	}

	if commandHonkMatch.MatchString(tweet.Text) {
		if strings.Contains(tweet.Text, "capybara") {
			matchTweet(twitterAPI, tweet, "honk", "capabara", "Honkbara the planet")
		} else {
			matchTweet(twitterAPI, tweet, "honk", "goose", "Honk the planet")
		}
	}

	if commandMeowMatch.MatchString(tweet.Text) {
		matchTweet(twitterAPI, tweet, "meow", "cat", "Meow the planet")
	}

	if commandPonyMatch.MatchString(tweet.Text) {
		matchTweet(twitterAPI, tweet, "pony", "pony", "#pony")
	}

	if commandWoofMatch.MatchString(tweet.Text) {
		matchTweet(twitterAPI, tweet, "woof", "dog", "woof woof woof")
	}

	if commandOinkMatch.MatchString(tweet.Text) {
		matchTweet(twitterAPI, tweet, "oink", "pig", "oink! oink!")
	}

	if commandQuackMatch.MatchString(tweet.Text) {
		matchTweet(twitterAPI, tweet, "quack", "duck", "quack! quack!")
	}

	if commandMooMatch.MatchString(tweet.Text) {
		matchTweet(twitterAPI, tweet, "moo", "cow", "Mooooo!")
	}

	if commandBaaMatch.MatchString(tweet.Text) {
		matchTweet(twitterAPI, tweet, "baa", "goat", "baa! baa!!")
	}
}

func matchTweet(twitterAPI *anaconda.TwitterApi, tweet anaconda.Tweet, matched, animal, message string) {
	log.Printf("Tweet matched %v", matched)

	if strings.Contains(tweet.Text, "RT ") {
		log.Println("This is a RT dont reply to not flood")
	}

	mention := tweet.InReplyToScreenName
	if mention == "" {
		mention = tweet.User.ScreenName
	}

	var image []byte
	switch animal {
	case "pony":
		image = getPony()
	default:
		image = getImage(animal)
	}
	if image == nil {
		image = getDefaultGoose()
	}
	msg := fmt.Sprintf("@%s %v #honkbot", mention, message)
	sendTweet(twitterAPI, tweet.IdStr, msg, image)
}

func sendTweet(twitterAPI *anaconda.TwitterApi, originalTweetID, message string, image []byte) {
	n := randomInt(30, 120)
	log.Printf("Sleeping for %s", time.Duration(n)*time.Second)
	time.Sleep(time.Duration(n) * time.Second)

	replyParams := url.Values{}
	msg := message
	mediaResponse, err := twitterAPI.UploadMedia(base64.StdEncoding.EncodeToString(image))
	if err != nil {
		log.Printf("Error uploading the image Err=%s\n", err.Error())
		msg = fmt.Sprintf("No image :(")
	} else {
		replyParams.Set("media_ids", mediaResponse.MediaIDString)
	}

	replyParams.Set("in_reply_to_status_id", originalTweetID)
	replyParams.Set("auto_populate_reply_metadata", "true")
	replyParams.Set("display_coordinates", "false")
	result, err := twitterAPI.PostTweet(msg, replyParams)
	if err != nil {
		log.Printf("Error while posting the tweet. Err=%s\n", err.Error())
		return
	}
	log.Printf("Tweet posted. TweetID = %s\n", result.IdStr)
}

// Returns an int >= min, < max
func randomInt(min, max int) int {
	return min + rand.Intn(max-min)
}
