package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/zabawaba99/firego"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/jwt"
)

func initialize() (*firego.Firebase, *jwt.Config, error) {
	d, err := ioutil.ReadFile("actions-smarthome-firebase-adminsdk.json")
	if err != nil {
		return nil, nil, err
	}

	conf, err := google.JWTConfigFromJSON(d, "https://www.googleapis.com/auth/userinfo.email",
		"https://www.googleapis.com/auth/firebase.database")
	if err != nil {
		return nil, nil, err
	}

	fb := firego.New("https://actions-smarthome.firebaseio.com/actions", conf.Client(oauth2.NoContext))
	return fb, conf, nil
}

func reconnectReference(fb *firego.Firebase, conf *jwt.Config) error {
	notifications := make(chan firego.Event)
	fb = firego.New("https://actions-smarthome.firebaseio.com/actions", conf.Client(oauth2.NoContext))
	return fb.Watch(notifications)
}

func executeCommand(fb *firego.Firebase, eventData string, room string) error {
	d := strings.Split(eventData, " ")
	if len(d) == 1 {
		if err := exec.Command("./cli/broadlink_cli", "--device", "@./cli/"+room+".device", "--send", "@./cli/"+d[0]).Run(); err != nil {
			return err
		}
	} else {
		rep, _ := strconv.Atoi(d[1])
		for i := 0; i < rep; i++ {
			if err := exec.Command("./cli/broadlink_cli", "--device", "@./cli/"+room+".device", "--send", "@./cli/"+eventData).Run(); err != nil {
				return err
			}

		}
	}
	if err := fb.Set("--"); err != nil {
		return err
	}
	return nil
}

func checkRegexp(reg, str string) bool {
	return regexp.MustCompile(reg).Match([]byte(str))
}

func main() {

	fb, conf, err := initialize()
	if err != nil {
		log.Fatal(err)
	}

	notifications := make(chan firego.Event)
	if err := fb.Watch(notifications); err != nil {
		fmt.Println("channel had something error")
		log.Fatal(err)
	}

	// defer fb.StopWatching()
	for {
		for event := range notifications {
			fmt.Printf("Type:%s Data:%s\n", event.Type, event.Data)

			// check event error
			if event.Type == firego.EventTypeAuthRevoked || event.Type == firego.EventTypeError {
				if err := reconnectReference(fb, conf); err != nil {
					fmt.Println("reconnect channel had something error")
					log.Fatal(err)
				}
				fmt.Println("reconnect firebase socket")
				continue
			}

			// check command
			if checkRegexp("ALL+", event.Data.(string)) {
				continue //TODO all off command

			} else if checkRegexp("--", event.Data.(string)) {
				continue

			} else if checkRegexp("BEDROOM+", event.Data.(string)) {
				if err := executeCommand(fb, event.Data.(string), "BEDROOM"); err != nil {
					log.Fatal(err)
				}

			} else {
				if err := executeCommand(fb, event.Data.(string), "LIVINGROOM"); err != nil {
					log.Fatal(err)
				}

			}
		}
	}

}
