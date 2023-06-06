package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

const (
	ClientSecret = ""
	ClientID     = ""

	spotifyBaseURL = "https://accounts.spotify.com"
	cbURL          = "http://localhost:8888/cb"

	TwitterAccessToken  = ""
	TwitterAccessSecret = ""

	TwitterClientSecret = ""
	TwitterClientID     = ""

	Twitter0AuthToken = ""

	TwitterConsumerKey    = ""
	TwitterConsumerSecret = ""
	BearerToken           = ""

	twitterCBURL = "http://localhost:8888/tweets/cb"
)

type SpotifyAuthtokenMeta struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

type SpotifyCurrentlyPlaying struct {
	Timestamp int64 `json:"timestamp"`
	Context   struct {
		ExternalUrls struct {
			Spotify string `json:"spotify"`
		} `json:"external_urls"`
		Href string `json:"href"`
		Type string `json:"type"`
		URI  string `json:"uri"`
	} `json:"context"`
	ProgressMs int `json:"progress_ms"`
	Item       struct {
		Album struct {
			AlbumType string `json:"album_type"`
			Artists   []struct {
				ExternalUrls struct {
					Spotify string `json:"spotify"`
				} `json:"external_urls"`
				Href string `json:"href"`
				ID   string `json:"id"`
				Name string `json:"name"`
				Type string `json:"type"`
				URI  string `json:"uri"`
			} `json:"artists"`
			ExternalUrls struct {
				Spotify string `json:"spotify"`
			} `json:"external_urls"`
			Href   string `json:"href"`
			ID     string `json:"id"`
			Images []struct {
				Height int    `json:"height"`
				URL    string `json:"url"`
				Width  int    `json:"width"`
			} `json:"images"`
			Name                 string `json:"name"`
			ReleaseDate          string `json:"release_date"`
			ReleaseDatePrecision string `json:"release_date_precision"`
			TotalTracks          int    `json:"total_tracks"`
			Type                 string `json:"type"`
			URI                  string `json:"uri"`
		} `json:"album"`
		DiscNumber  int  `json:"disc_number"`
		DurationMs  int  `json:"duration_ms"`
		Explicit    bool `json:"explicit"`
		ExternalIds struct {
			Isrc string `json:"isrc"`
		} `json:"external_ids"`
		ExternalUrls struct {
			Spotify string `json:"spotify"`
		} `json:"external_urls"`
		Href        string `json:"href"`
		ID          string `json:"id"`
		IsLocal     bool   `json:"is_local"`
		Name        string `json:"name"`
		Popularity  int    `json:"popularity"`
		PreviewURL  string `json:"preview_url"`
		TrackNumber int    `json:"track_number"`
		Type        string `json:"type"`
		URI         string `json:"uri"`
	} `json:"item"`
	CurrentlyPlayingType string `json:"currently_playing_type"`
	Actions              struct {
		Disallows struct {
			Resuming bool `json:"resuming"`
		} `json:"disallows"`
	} `json:"actions"`
	IsPlaying bool `json:"is_playing"`
}

var cache map[string]interface{}

func main() {

	cln := Get()
	cln.Configure(
		SetSpotifyClientKey(ClientID),
		SetSpotifySecretKey(ClientSecret),
	)
	fmt.Println(c)

	GetTweet()

	router := mux.NewRouter()
	router.HandleFunc("/token", authKeyGeneratorHandler).Methods("GET")
	router.HandleFunc("/cb", callbackHandler).Methods("GET")

	fmt.Println("Server running")
	log.Fatal(http.ListenAndServe(":8888", router))
}

func authKeyGeneratorHandler(w http.ResponseWriter, r *http.Request) {
	scope := "user-read-currently-playing"

	redirectURL := "https://accounts.spotify.com/authorize?" +
		"&response_type=code" +
		"&client_id=" + ClientID +
		"&scope=" + scope +
		"&redirect_uri=" + cbURL

	http.Redirect(w, r, redirectURL, http.StatusFound)
}

func callbackHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	code := query.Get("code")

	fmt.Println("Authorization Code:", code)
	res, err := getAccessToken(code)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	var rez SpotifyAuthtokenMeta

	jErr := json.Unmarshal(res, &rez)
	if err != nil {
		fmt.Println(jErr.Error())
		return
	}

	cache = make(map[string]interface{})

	cache["access_token"] = rez.AccessToken

	queryPlaylists()

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Callback Received"))
}

func queryPlaylists() {
	cTrack, cErr := c.makeGetRequest(context.Background(), "https://api.spotify.com/v1/me/player/currently-playing")
	if cErr != nil {
		fmt.Println("cTrack:  ", cErr.Error())
		return
	}

	if len(cTrack) > 0 {
		var track SpotifyCurrentlyPlaying

		err := json.Unmarshal(cTrack, &track)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		if cache["current_track"] != track.Item.Album.ExternalUrls.Spotify {
			cache["current_track"] = track.Item.Album.ExternalUrls.Spotify

			t.sendTweet(cache["current_track"].(string))
		}
	}

	fmt.Println(cache)
}

func getAccessToken(code string) ([]byte, error) {
	return c.makePostRequest(context.Background(), spotifyBaseURL+"/api/token", map[string]string{
		"code":         code,
		"redirect_uri": cbURL,
		"grant_type":   "authorization_code",
	})
}
