package notify

import (
	"firebase.google.com/go/messaging"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"google.golang.org/api/option"
	"os"

	firebase "firebase.google.com/go"
)

var (
	App *firebase.App
)

func InitializeFirebase() {
	logrus.Info("connecting to firebase...")
	opt := option.WithCredentialsJSON([]byte(os.Getenv("FIREBASE_PRIVATE_KEY")))
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		logrus.Panic("Could not establish firebase connection: ", err)
		return
	}
	App = app
	logrus.Info("Firebase connection established!")
}

func SendNotification(registrationTokens []string, data map[string]string) ([]string, error) {
	ctx := context.Background()
	client, err := App.Messaging(ctx)
	if err != nil {
		logrus.Error("error getting Messaging client: ", err)
		return []string{}, err
	}
	message := &messaging.MulticastMessage{
		Data:   data,
		Tokens: registrationTokens,
	}

	br, err := client.SendMulticast(context.Background(), message)

	if err != nil {
		logrus.Error("error sending notifications: ", err)
		return []string{}, err
	}
	logrus.Info(br, err)
	var failedTokens []string
	if br != nil && br.FailureCount > 0 {
		for idx, resp := range br.Responses {
			if !resp.Success {
				failedTokens = append(failedTokens, registrationTokens[idx])
			}
		}
		logrus.Info("Failed Tokens: ", failedTokens)
	}
	return failedTokens, nil
}
