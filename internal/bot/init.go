package bot

import (
	"context"
	"fmt"
	"os"

	qrterminal "github.com/mdp/qrterminal/v3"
	"go.mau.fi/whatsmeow"
)

type Bot struct {
	Client *whatsmeow.Client
}

func Start() *Bot {
	c, err := Connection()
	if err != nil {
		panic(err)
	}

	b := &Bot{Client: c}

	if b.Client.Store.ID == nil {
		qrChan, _ := b.Client.GetQRChannel(context.Background())

		err = b.Client.Connect()
		if err != nil {
			panic(err)
		}

		for evt := range qrChan {
			if evt.Event == "code" {
				qrterminal.Generate(evt.Code, qrterminal.L, os.Stdout)
			} else {
				fmt.Println("Login Event:", evt.Event)
			}
		}
	} else {
		err = b.Client.Connect()
		if err != nil {
			panic(err)
		}
	}

	go b.StartServer()

	return b
}
