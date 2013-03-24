package main

import (
	"encoding/json"
	"flag"
	"fmt"
	irc "github.com/fluffle/goirc/client"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
)

var (
	channels   = []string{"#collinjester"}
	HttpClient = &http.Client{}
)

type GisResult struct {
	Url string
}

type GisResponse struct {
	ResponseData struct {
		Results []GisResult
	}
	Results []GisResult
}

// returns true if the url ends with some well known suffixes
func isImage(url string) bool {
	suffixes := []string{".jpg", ".jpeg", ".gif", ".png", ".bmp"}
	lowered := strings.ToLower(url)
	for _, suffix := range suffixes {
		if strings.HasSuffix(lowered, suffix) {
			return true
		}
	}
	return false
}

// when an error occurs, calling this method will send the error back to the irc channel
func sendError(conn *irc.Conn, channel string, nick string, err error) {
	log.Print(err)
	conn.Privmsg(channel, fmt.Sprintf("%s: oops: %v", nick, err))
}

func googleImageSearch(conn *irc.Conn, channel string, nick string, search string) {
	url := "http://ajax.googleapis.com/ajax/services/search/images?v=1.0&q=" + url.QueryEscape(search)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		sendError(conn, channel, nick, err)
		return
	}
	resp, err := HttpClient.Do(req)
	if err != nil {
		sendError(conn, channel, nick, err)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		sendError(conn, channel, nick, err)
		return
	}
	var gisResponse GisResponse
	json.Unmarshal(body, &gisResponse)
	if len(gisResponse.ResponseData.Results) > 0 {
		indexes := rand.Perm(len(gisResponse.ResponseData.Results))
		for i := 0; i < len(indexes); i++ {
			imageUrl := gisResponse.ResponseData.Results[indexes[i]].Url
			if isImage(imageUrl) {
				conn.Privmsg(channel, fmt.Sprintf("%s: %s", nick, imageUrl))
				return
			}
		}
		conn.Privmsg(channel, fmt.Sprintf("%s: %s", nick, "Nothing found"))
	}
}

// findLinks returns a slice of strings that look like links. adds a protocol to the beginning of 
// the link if it doesn't already have one
func findLinks(message string) []string {
	prefixes := []string{ "t.co", "cl.ly", "www", "bit.ly", "j.mp", "tcrn.ch"}
	result := make([]string, 0)
	for _, token := range strings.Split(message, " ") {
		if strings.HasPrefix(token, "http") {
			result = append(result, token)
		} else {
			// check to see if it looks like it might be a link
			for _, prefix := range prefixes {
				if strings.HasPrefix(token, prefix) {
					result = append(result, "http://" + token)
					break
				}
			}
		}
	}
	return result
}

// expandLink fully un-shortens a url
func expandLink(url string) (expanded string, err error) {
	log.Printf("Expanding link %s\n", url)
	resp, err := HttpClient.Head(url) // will follow redirects
	if err != nil {
		return expanded, err
	}
	defer resp.Body.Close() // not sure if i have to do this with a head response
	expanded = resp.Request.URL.String()
	if expanded != url {
		return
	}
	return "", nil
}

// takes an input line and rewrite any links that are shortened links into their full representation
func rewriteLine(conn *irc.Conn, line *irc.Line) {
	channel := line.Args[0]
	message := line.Args[1]
	for _, link := range findLinks(message) {
		expandedLink, err := expandLink(link)
		if err != nil {
			log.Printf("Could not expand link %s: %s", link, err)
		} else if expandedLink != "" {
			conn.Privmsg(channel, fmt.Sprintf("%s: %s", line.Nick, expandedLink))
		}
	}
}

// When a message comes in on a channel gesture has joined, this method will be called.
func messageReceived(conn *irc.Conn, line *irc.Line) {
	if len(line.Args) > 1 {
		channel := line.Args[0]
		message := line.Args[1]
		messageSliced := strings.Split(message, " ")
		command := messageSliced[0]
		commandArgs := messageSliced[1:]

		fmt.Printf(">> %s (%s): %s\n", line.Nick, channel, message)

		if command == "gis" && len(commandArgs) >= 1 {
			googleImageSearch(conn, channel, line.Nick, strings.Join(commandArgs, " "))
		} else if command == "echo" {
			response := line.Nick + ": " + message
			conn.Privmsg(channel, response)
		} else {
			rewriteLine(conn, line)
		}
	}
}

func main() {
	flag.Parse()
	c := irc.SimpleClient("gesturebot")
	c.SSL = true
	c.AddHandler(irc.CONNECTED,
		func(conn *irc.Conn, line *irc.Line) {
			for _, channel := range channels {
				conn.Join(channel)
			}
		})
	quit := make(chan bool)
	c.AddHandler(irc.DISCONNECTED, func(conn *irc.Conn, line *irc.Line) { quit <- true })
	c.AddHandler("PRIVMSG", func(conn *irc.Conn, line *irc.Line) {
		messageReceived(conn, line)
	})
	if err := c.Connect("irc.freenode.net"); err != nil {
		fmt.Printf("Connection error: %s\n", err)
	}
	// Wait for disconnect
	<-quit
}
