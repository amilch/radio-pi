package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"time"

	// "syscall"
	// "time"

	"log"

	"github.com/gorilla/mux"
)

type Station struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Url      string `json:"url"`
	ImageUrl string `json:"image_url"`
}

type Response struct {
	Data  string `json:"data"`
	Error string `json:"error"`
}

type MpvJsonResponse struct {
	Data      string `json:"data"`
	RequestId int    `json:"request_id"`
	Error     string `json:"error"`
}

type MpvJsonIntResponse struct {
	Data      float64 `json:"data"`
	RequestId int     `json:"request_id"`
	Error     string  `json:"error"`
}

func SendCommand(cmd ...string) string {
	c, err := net.Dial("unix", "/tmp/mpvsocket")
	if err != nil {
		log.Fatal("Could not connect to socket:", err)
	}
	defer c.Close()

	cmdString := `{ "command": ["` + strings.Join(cmd, `", "`) + `"] }` + "\n"
	log.Print(cmdString)

	_, err = c.Write([]byte(cmdString))
	if err != nil {
		log.Fatal("Could not write to socket:", err)
	}

	res, err := bufio.NewReader(c).ReadString('\n')
	if err != nil {
		log.Fatal("Could not read from socket:", err)
	}
	return res
}

func StationsHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(stations)
}
func PlayHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := r.URL.Query()["id"]
	if !ok || len(id[0]) < 1 {
		json.NewEncoder(w).Encode(Response{Error: "id paramter missing"})
		return
	}

	url := ""
	for i := 0; i < len(stations); i++ {
		if stations[i].ID == id[0] {
			url = stations[i].Url
			break
		}
	}
	if url == "" {
		json.NewEncoder(w).Encode(Response{Error: "station not found"})
		return
	}

	SendCommand("loadfile", url)
	json.NewEncoder(w).Encode(Response{Data: "ok"})
}
func StopHandler(w http.ResponseWriter, r *http.Request) {
	SendCommand("stop")
	json.NewEncoder(w).Encode(Response{Data: "ok"})
}
func PlayingHandler(w http.ResponseWriter, r *http.Request) {
	resJson := SendCommand("get_property", "path")
	var res MpvJsonResponse
	err := json.Unmarshal([]byte(resJson), &res)
	if err != nil {
		log.Fatal(err)
	}

	if res.Error == "property unavailable" {
		json.NewEncoder(w).Encode(Response{Data: ""})
		return
	}

	fmt.Print(res.Data)

	for i := 0; i < len(stations); i++ {
		if stations[i].Url == res.Data {
			json.NewEncoder(w).Encode(Response{Data: stations[i].ID})
			return
		}
	}
	json.NewEncoder(w).Encode(Response{Error: "id of current station not found"})
}
func TrackTitleHandler(w http.ResponseWriter, r *http.Request) {
	resJson := SendCommand("get_property", "media-title")
	var res MpvJsonResponse
	err := json.Unmarshal([]byte(resJson), &res)
	if err != nil {
		log.Fatal(err)
	}
	json.NewEncoder(w).Encode(Response{Data: res.Data})
}
func VolumeHandler(w http.ResponseWriter, r *http.Request) {
	set, ok := r.URL.Query()["set"]
	if !ok || len(set[0]) < 1 {
		resJson := SendCommand("get_property", "volume")
		var res MpvJsonIntResponse
		err := json.Unmarshal([]byte(resJson), &res)
		if err != nil {
			log.Fatal(err)
		}
		json.NewEncoder(w).Encode(Response{Data: fmt.Sprintf("%.0f", res.Data)})
		return
	}

	SendCommand("set_property", "volume", set[0])
	json.NewEncoder(w).Encode(Response{Data: "ok"})
}

var stations []Station
var currentlyPlaying string
var play chan string
var stop chan bool

func main() {
	play = make(chan string)
	stop = make(chan bool)
	stations = []Station{
		Station{"0", "BBC Radio 1", "http://bbcmedia.ic.llnwd.net/stream/bbcmedia_radio1_mf_p", "https://d3kle7qwymxpcy.cloudfront.net/images/broadcasts/38/81/3243/2/c175.png"},
		Station{"1", "FluxFM", "http://streams.fluxfm.de/live/mp3-128/audio/", "https://d3kle7qwymxpcy.cloudfront.net/images/broadcasts/01/ee/2226/1/c175.png"},
		Station{"2", "FluxFM", "http://streams.fluxfm.de/live/mp3-128/audio/", "https://d3kle7qwymxpcy.cloudfront.net/images/broadcasts/01/ee/2226/1/c175.png"},
		Station{"3", "FluxFM", "http://streams.fluxfm.de/live/mp3-128/audio/", "https://d3kle7qwymxpcy.cloudfront.net/images/broadcasts/01/ee/2226/1/c175.png"},
		Station{"4", "FluxFM", "http://streams.fluxfm.de/live/mp3-128/audio/", "https://d3kle7qwymxpcy.cloudfront.net/images/broadcasts/01/ee/2226/1/c175.png"},
		Station{"5", "FluxFM", "http://streams.fluxfm.de/live/mp3-128/audio/", "https://d3kle7qwymxpcy.cloudfront.net/images/broadcasts/01/ee/2226/1/c175.png"},
		Station{"6", "FluxFM", "http://streams.fluxfm.de/live/mp3-128/audio/", "https://d3kle7qwymxpcy.cloudfront.net/images/broadcasts/01/ee/2226/1/c175.png"},
		Station{"7", "FluxFM", "http://streams.fluxfm.de/live/mp3-128/audio/", "https://d3kle7qwymxpcy.cloudfront.net/images/broadcasts/01/ee/2226/1/c175.png"},
		Station{"8", "FluxFM", "http://streams.fluxfm.de/live/mp3-128/audio/", "https://d3kle7qwymxpcy.cloudfront.net/images/broadcasts/01/ee/2226/1/c175.png"},
		Station{"9", "FluxFM", "http://streams.fluxfm.de/live/mp3-128/audio/", "https://d3kle7qwymxpcy.cloudfront.net/images/broadcasts/01/ee/2226/1/c175.png"},
		Station{"10", "FluxFM", "http://streams.fluxfm.de/live/mp3-128/audio/", "https://d3kle7qwymxpcy.cloudfront.net/images/broadcasts/01/ee/2226/1/c175.png"},
		Station{"11", "FluxFM", "http://streams.fluxfm.de/live/mp3-128/audio/", "https://d3kle7qwymxpcy.cloudfront.net/images/broadcasts/01/ee/2226/1/c175.png"},
	}

	router := mux.NewRouter()
	router.HandleFunc("/stations", StationsHandler)
	router.HandleFunc("/play", PlayHandler)
	router.HandleFunc("/stop", StopHandler)
	router.HandleFunc("/playing", PlayingHandler)
	router.HandleFunc("/title", TrackTitleHandler)
	router.HandleFunc("/volume", VolumeHandler)

	srv := &http.Server{Addr: ":5051", Handler: router}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	go func() {
		cmd := exec.Command("mpv", "http://streams.fluxfm.de/live/mp3-128/audio/", "--idle=yes", "--input-ipc-server=/tmp/mpvsocket")
		if err := cmd.Run(); err != nil {
			log.Printf("Mpv error: %s", err)
			done <- syscall.SIGINT
		}
	}()
	log.Print("Server started")

	<-done
	log.Print("Server stopped")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}
	log.Print("Server exited properly")
}
