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
	commandHissMatch  = regexp.MustCompile(`(?mi)^*/(?:hiss)(?: +(.+?))?\s*$`)
	commandSnekMatch  = regexp.MustCompile(`(?mi)^*/(?:snek)(?: +(.+?))?\s*$`)
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
	streamValues.Set("track", "/honk,/meow,/pony,/woof,/oink,/quack,/moo,/baa,/snek,/hiss")
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
		commandHissMatch.MatchString(tweet.Text) || commandSnekMatch.MatchString(tweet.Text) ||
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
		log.Println("Tweet matched honk")

		if strings.Contains(tweet.Text, "RT ") {
			log.Println("This is a RT dont reply to not flood")
			return
		}

		mention := tweet.InReplyToScreenName
		if mention == "" {
			mention = tweet.User.ScreenName
		}

		var image []byte
		var msg string
		if strings.Contains(tweet.Text, "capybara") {
			image = getImage("capybara")
			msg = fmt.Sprintf("@%s Honkbara the Planet", mention)
		} else {
			image = getImage("goose")
			msg = fmt.Sprintf("@%s Honk the Planet #honkbot", mention)
		}

		if image == nil {
			image = getDefaultGoose()
		}

		sendTweet(twitterAPI, tweet.IdStr, msg, image)
	}

	if commandMeowMatch.MatchString(tweet.Text) {
		log.Println("Tweet matched meow")

		if strings.Contains(tweet.Text, "RT ") {
			log.Println("This is a RT dont reply to not flood")
			return
		}

		mention := tweet.InReplyToScreenName
		if mention == "" {
			mention = tweet.User.ScreenName
		}

		image := getImage("cat")
		if image == nil {
			image = getDefaultGoose()
		}

		msg := fmt.Sprintf("@%s Meow the Planet #honkbot", mention)
		sendTweet(twitterAPI, tweet.IdStr, msg, image)
	}

	if commandPonyMatch.MatchString(tweet.Text) {
		log.Println("Tweet matched pony")

		if strings.Contains(tweet.Text, "RT ") {
			log.Println("This is a RT dont reply to not flood")
			return
		}

		mention := tweet.InReplyToScreenName
		if mention == "" {
			mention = tweet.User.ScreenName
		}

		image := getPony()
		if image == nil {
			image = getDefaultGoose()
		}

		msg := fmt.Sprintf("@%s #honkbot #pony", mention)
		sendTweet(twitterAPI, tweet.IdStr, msg, image)
	}

	if commandWoofMatch.MatchString(tweet.Text) {
		log.Println("Tweet matched woof")

		if strings.Contains(tweet.Text, "RT ") {
			log.Println("This is a RT dont reply to not flood")
			return
		}

		mention := tweet.InReplyToScreenName
		if mention == "" {
			mention = tweet.User.ScreenName
		}

		image := getImage("dog")
		if image == nil {
			image = getDefaultGoose()
		}

		msg := fmt.Sprintf("@%s woof woof woof #honkbot", mention)
		sendTweet(twitterAPI, tweet.IdStr, msg, image)
	}

	if commandOinkMatch.MatchString(tweet.Text) {
		log.Println("Tweet matched oink")

		if strings.Contains(tweet.Text, "RT ") {
			log.Println("This is a RT dont reply to not flood")
			return
		}

		mention := tweet.InReplyToScreenName
		if mention == "" {
			mention = tweet.User.ScreenName
		}

		image := getImage("pig")
		if image == nil {
			image = getDefaultGoose()
		}

		msg := fmt.Sprintf("@%s oink! oink! #honkbot", mention)
		sendTweet(twitterAPI, tweet.IdStr, msg, image)
	}

	if commandQuackMatch.MatchString(tweet.Text) {
		log.Println("Tweet matched quack")

		if strings.Contains(tweet.Text, "RT ") {
			log.Println("This is a RT dont reply to not flood")
			return
		}

		mention := tweet.InReplyToScreenName
		if mention == "" {
			mention = tweet.User.ScreenName
		}

		image := getImage("duck")
		if image == nil {
			image = getDefaultGoose()
		}

		msg := fmt.Sprintf("@%s quack! quack! #honkbot", mention)
		sendTweet(twitterAPI, tweet.IdStr, msg, image)
	}

	if commandMooMatch.MatchString(tweet.Text) {
		log.Println("Tweet matched moo")

		if strings.Contains(tweet.Text, "RT ") {
			log.Println("This is a RT dont reply to not flood")
			return
		}

		mention := tweet.InReplyToScreenName
		if mention == "" {
			mention = tweet.User.ScreenName
		}

		image := getImage("cow")
		if image == nil {
			image = getDefaultGoose()
		}

		msg := fmt.Sprintf("@%s Mooooo! #honkbot", mention)
		sendTweet(twitterAPI, tweet.IdStr, msg, image)
	}

	if commandBaaMatch.MatchString(tweet.Text) {
		log.Println("Tweet matched baa")

		if strings.Contains(tweet.Text, "RT ") {
			log.Println("This is a RT dont reply to not flood")
			return
		}

		mention := tweet.InReplyToScreenName
		if mention == "" {
			mention = tweet.User.ScreenName
		}

		image := getImage("goat")
		if image == nil {
			image = getDefaultGoose()
		}
		msg := fmt.Sprintf("@%s baa! baa!! #honkbot", mention)
		sendTweet(twitterAPI, tweet.IdStr, msg, image)
	}

	if commandHissMatch.MatchString(tweet.Text) {
		log.Println("Tweet matched hiss")

		if strings.Contains(tweet.Text, "RT ") {
			log.Println("This is a RT dont reply to not flood")
			return
		}

		mention := tweet.InReplyToScreenName
		if mention == "" {
			mention = tweet.User.ScreenName
		}

		image := getImage("snake")
		if image == nil {
			image = getDefaultGoose()
		}
		msg := fmt.Sprintf("@%s hiss! hiiiiissssss!! #honkbot", mention)
		sendTweet(twitterAPI, tweet.IdStr, msg, image)
	}

	if commandSnekMatch.MatchString(tweet.Text) {
		log.Println("Tweet matched sneck")

		if strings.Contains(tweet.Text, "RT ") {
			log.Println("This is a RT dont reply to not flood")
			return
		}

		mention := tweet.InReplyToScreenName
		if mention == "" {
			mention = tweet.User.ScreenName
		}

		image := getImage("snake")
		if image == nil {
			image = getDefaultGoose()
		}
		msg := fmt.Sprintf("@%s beware the snek! #honkbot", mention)
		sendTweet(twitterAPI, tweet.IdStr, msg, image)
	}

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
