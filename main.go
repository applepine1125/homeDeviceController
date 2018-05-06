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
		fmt.Println("execute command once")
		if err := exec.Command("./cli/broadlink_cli", "--device", "@./cli/"+room+".device", "--send", "@./cli/"+d[0]).Run(); err != nil {
			return err
		}
	} else {
		rep, _ := strconv.Atoi(d[1])
		fmt.Printf("execute command %d times\n", rep)
		for i := 0; i < rep*2; i++ {
			if err := exec.Command("./cli/broadlink_cli", "--device", "@./cli/"+room+".device", "--send", "@./cli/"+d[0]).Run(); err != nil {
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
	for {
		fb, _, err := initialize()
		if err != nil {
			log.Fatal(err)
		}

		notifications := make(chan firego.Event)
		if err := fb.Watch(notifications); err != nil {
			fmt.Println("channel had something error")
			log.Fatal(err)
		}

		for event := range notifications {
			fmt.Printf("Type:%s Data:%s\n", event.Type, event.Data)

			//check event error
			if event.Type == firego.EventTypeError || event.Type == firego.EventTypeAuthRevoked {
				break
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
		fb.StopWatching()
		fmt.Println("reconnect firebase socket")
	}
}
