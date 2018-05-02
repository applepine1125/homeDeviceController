package main

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/zabawaba99/firego"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func main() {
	d, err := ioutil.ReadFile("actions-smarthome-firebase-adminsdk.json")
	if err != nil {
		log.Fatal(err)
	}

	conf, err := google.JWTConfigFromJSON(d, "https://www.googleapis.com/auth/userinfo.email",
		"https://www.googleapis.com/auth/firebase.database")
	if err != nil {
		log.Fatal(err)
	}

	fb := firego.New("https://actions-smarthome.firebaseio.com/actions", conf.Client(oauth2.NoContext))

	notifications := make(chan firego.Event)
	if err := fb.Watch(notifications); err != nil {
		log.Fatal(err)
	}

	defer fb.StopWatching()

	// var ss []string
	for event := range notifications {
		fmt.Printf("%s\n", event.Data)

		// 	ss = append([]string{"SEND_ONCE"}, strings.Split(event.Data.(string), " ")...)
		//
		// 	err := exec.Command("irsend", ss...).Run()
		// 	if err != nil {
		// 		log.Fatal(err)
		// 	}
	}

	fmt.Printf("Notifications have stopped")
}
