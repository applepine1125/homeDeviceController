package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"regexp"

	"github.com/zabawaba99/firego"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func main() {
	bbcPath := "./cli/broadlink_cli"

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
		// err := exec.Command(bbcPath, "--device @LIVINGROOM.device --send", "@"+event.Data.(string)).Run()
		if check_regexp("ALL+", event.Data.(string)) {
		} else {
			err := exec.Command(bbcPath, "--device", "@./cli/LIVINGROOM.device", "--send", "@./cli/"+event.Data.(string)).Run()
			if err != nil {
				log.Fatal(err)
			}
		}
	}
	fmt.Printf("Notifications have stopped")
}

func check_regexp(reg, str string) bool {
	return regexp.MustCompile(reg).Match([]byte(str))
}
