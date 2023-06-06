package main

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/dghubble/go-twitter/twitter"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

var t *tweet

type authorize struct {
	Token string
}

type tweet struct {
	client *twitter.Client
}

func (a authorize) Add(req *http.Request) {
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", a.Token))
}

func newTweeter() *tweet {
	config := &clientcredentials.Config{
		ClientID:     TwitterAccessToken,
		ClientSecret: TwitterClientSecret,
		TokenURL:     "https://api.twitter.com/oauth2/token",
	}

	// http.Client will automatically authorize Requests
	httpClient := config.Client(oauth2.NoContext)

	// Twitter client
	cl := twitter.NewClient(httpClient)
	t := &tweet{
		client: cl,
	}

	return t
}

func GetTweet() *tweet {
	if t != nil {
		return t
	}

	t = newTweeter()

	return t
}

var nonce string
var nErr error
var timestamp string

func (tw *tweet) sendTweet(text string) {
	nonce = *generateNonce(24)
	timestamp = strconv.FormatInt(time.Now().Unix(), 10)

	urlStr := "https://api.twitter.com/1.1/statuses/update.json" // Target URL

	// Create a map to hold the OAuth parameters
	params := make(map[string]string)

	// params["status"] = "helloe"

	authSign := generateOAuthSignature(urlStr, params)
	generateOAuthParameters(urlStr, params)
	fmt.Println("authSign", authSign)

	authHeader := generateOAuthAuthorizationHeader(params)

	fmt.Println("authHeaderauthHeader", authHeader)

	// Create the HTTP client
	client := http.DefaultClient

	// Create the form data
	formData := url.Values{}
	formData.Set("status", "Hello, Twitter!")

	// Create the HTTP request
	req, err := http.NewRequest(http.MethodPost, urlStr, strings.NewReader(formData.Encode()))
	if err != nil {
		fmt.Println("Failed to create request:", err)
		return
	}

	req.Header.Set("Authorization", authHeader)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Send the HTTP request
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Failed to send request:", err)
		return
	}
	defer resp.Body.Close()

	fmt.Println("resprespresp", resp)

	// Check the response status
	if resp.StatusCode != http.StatusOK {
		fmt.Println("Request failed with status:", resp.Status)
		return
	}

	fmt.Println("Tweet sent successfully!")

	// tweet, resp, err := tw.client.Statuses.Update("req", nil)
	// if err != nil {
	// 	fmt.Println("create tweet error: %v", err.Error())
	// 	return
	// }

	// fmt.Println("tweettweet", tweet)
	// fmt.Println("resprespresp", resp)

	return
}

func generateNonce(length int) *string {
	// Calculate the required number of bytes to generate a base64 string of the desired length
	numBytes := (length * 6) / 8

	// Generate random bytes
	bytes := make([]byte, numBytes)
	_, err := rand.Read(bytes)
	if err != nil {

		return nil
	}

	// Encode bytes to base64
	base64String := base64.URLEncoding.EncodeToString(bytes)

	// Trim any padding characters from the base64 string
	base64String = base64String[:length]

	return &base64String
}

func generateOAuthSignature(urlStr string, params map[string]string) *string {
	// Create a sorted list of parameter keys
	var keys []string
	for key := range params {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	// Build the parameter string
	var paramStr strings.Builder
	for i, key := range keys {
		if i > 0 {
			paramStr.WriteString("&")
		}
		paramStr.WriteString(url.QueryEscape(key))
		paramStr.WriteString("=")
		paramStr.WriteString(url.QueryEscape(params[key]))
	}

	// Build the base string
	baseString := strings.Join([]string{
		"POST",
		urlStr,
		url.QueryEscape(paramStr.String()),
	}, "&")

	// Generate the signing key
	signingKey := strings.Join([]string{
		url.QueryEscape(TwitterConsumerSecret),
		url.QueryEscape(TwitterAccessSecret),
	}, "&")

	// Generate the signature using HMAC-SHA1
	hmacHash := hmac.New(sha1.New, []byte(signingKey))
	hmacHash.Write([]byte(baseString))
	signature := base64.StdEncoding.EncodeToString(hmacHash.Sum(nil))

	return &signature
}

func generateOAuthParameters(apiURL string, params map[string]string) map[string]string {
	// Add the OAuth parameters to the map
	encodedURL := url.QueryEscape(apiURL)

	params["oauth_consumer_key"] = TwitterConsumerKey
	params["oauth_nonce"] = nonce
	params["oauth_signature"] = *generateOAuthSignature(encodedURL, params)
	params["oauth_signature_method"] = "HMAC-SHA1"
	params["oauth_timestamp"] = timestamp
	params["oauth_token"] = Twitter0AuthToken
	params["oauth_version"] = "1.0"

	return params
}

func generateOAuthAuthorizationHeader(params map[string]string) string {
	var headerParams []string
	for key, value := range params {
		param := fmt.Sprintf("%s=\"%s\"", url.QueryEscape(key), url.QueryEscape(value))
		headerParams = append(headerParams, param)
	}

	authorizationHeader := fmt.Sprintf("OAuth %s", strings.Join(headerParams, ", "))

	return authorizationHeader
}
