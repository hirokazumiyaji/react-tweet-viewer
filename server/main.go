package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/garyburd/go-oauth/oauth"
	"github.com/naoya/go-pit"
)

const endpoint = "https://api.twitter.com"

type Client struct {
	Client oauth.Client
	Token  *oauth.Credentials
}

func NewClient(config pit.Profile) (*Client, error) {
	c := Client{
		Client: oauth.Client{
			TemporaryCredentialRequestURI: endpoint + "/oauth/request/token",
			ResourceOwnerAuthorizationURI: endpoint + "/oauth/authenticate",
			TokenRequestURI:               endpoint + "/oauth/access_token",
		},
	}
	c.Client.Credentials.Token = config["consumer_key"]
	c.Client.Credentials.Secret = config["consumer_secret"]
	accessToken, isFoundAccessToken := config["access_token"]
	accessTokenSecret, isFoundAccessTokenSecret := config["access_token_secret"]
	if isFoundAccessToken && isFoundAccessTokenSecret {
		c.Token = &oauth.Credentials{accessToken, accessTokenSecret}
	} else {
		token, err := c.Client.RequestTemporaryCredentials(http.DefaultClient, "", nil)
		if err != nil {
			fmt.Println("failed to request temporary credentials: ", err)
			return nil, err
		}
		c.Token = token
	}
	return &c, nil
}

func (c *Client) HomeTimeline(q url.Values) ([]Tweet, error) {
	u := endpoint + "/1.1/statuses/home_timeline.json"
	client.Client.SignParam(client.Token, "GET", u, q)
	u = u + "?" + q.Encode()
	res, err := http.Get(u)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, err
	}
	var timeline []Tweet
	if err = json.NewDecoder(res.Body).Decode(&timeline); err != nil {
		return nil, err
	}
	return timeline, nil
}

type Tweet struct {
	Text string `json:"text"`
	ID   string `json:"id_str"`
	User struct {
		ScreenName      string `json:"screen_name"`
		ProfileImageURL string `json:"profile_image_url"`
	} `json:"user"`
}

var client *Client

func homeTimeline(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	timeline, err := client.HomeTimeline(r.URL.Query())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	res, err := json.Marshal(timeline)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(res)
}

func main() {
	config, err := pit.Get("twitter.com")
	if err != nil {
		fmt.Println("pit.Get Error: ", err)
		os.Exit(1)
	}

	client, err = NewClient(config)
	if err != nil {
		fmt.Println("NewClient Error: ", err)
		os.Exit(1)
	}

	http.HandleFunc("/", homeTimeline)
	fmt.Println("listening...")
	http.ListenAndServe(":8080", nil)
}
