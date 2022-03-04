package main

import (
	goContext "context"
	"fmt"
	"log"
	"os"
	"strings"

	// "time"
	// "strconv"

	"github.com/adrg/xdg"
	_ "github.com/mattn/go-sqlite3"
	"github.com/nanu-c/qml-go"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
)

type Client struct {
	Root       qml.Object
	LoginToken string
	Status     string
	Message    string
	wmClient   *whatsmeow.Client
}

func run() error {
	engine := qml.NewEngine()
	component, err := engine.LoadFile("qml/Main.qml")
	if err != nil {
		return err
	}

	qmlBridge := createClient()
	context := engine.Context()
	context.SetVar("qmlBridge", qmlBridge)
	qmlBridge.connect()

	win := component.CreateWindow(nil)
	qmlBridge.Root = win.Root()
	win.Show()
	win.Wait()

	return nil
}

func createClient() *Client {

	var dbPath, err = xdg.ConfigFile("whats-ut.walking-octopus/userStore.db")
	if err != nil {
		panic(err)
	}
	os.Mkdir(strings.Replace(dbPath, "/userStore.db", "", 1), 0755)

	container, err := sqlstore.New("sqlite3", dbPath+"?_foreign_keys=on", waLog.Stdout("Database", "DEBUG", true))
	if err != nil {
		panic(err)
	}

	deviceStore, err := container.GetFirstDevice()
	if err != nil {
		panic(err)
	}

	c := &Client{LoginToken: ""}
	c.wmClient = whatsmeow.NewClient(deviceStore, waLog.Stdout("Client", "DEBUG", true))
	c.wmClient.AddEventHandler(c.handler)
	return c
}

func (c *Client) isConnected() bool {
	return c.wmClient.Store.ID != nil
}

func (c *Client) connect() {
	if !c.isConnected() {
		fmt.Println("main.go: WhatsApp(): No Client Store ID")
		qrChan, _ := c.wmClient.GetQRChannel(goContext.Background())
		err := c.wmClient.Connect()
		if err != nil {
			panic(err)
		}

		go func() {
			for evt := range qrChan {
				if evt.Event == "code" {
					fmt.Println("QR code:", evt.Code)

					// FixMe: The LoginToken can't be redefined from here or the handler usually, but not all of the time
					c.setLoginToken(evt.Code)
					c.setStatus("QR")
				} else {
					fmt.Println("Login event:", evt.Event)
				}
			}
		}()
	} else {
		fmt.Println("main.go: WhatsApp(): Client Store ID:" + c.wmClient.Store.ID.User)
		err := c.wmClient.Connect()
		if err != nil {
			panic(err)
		}
	}
}

func (c *Client) setLoginToken(token string) {
	c.LoginToken = token
	qml.Changed(c, &c.LoginToken)
}

func (c *Client) setStatus(status string) {
	c.Status = status
	qml.Changed(c, &c.Status)
}

func (c *Client) setMessage(message string) {
	c.Message = message
	qml.Changed(c, &c.Message)
}

func (c *Client) handler(rawEvt interface{}) {
	switch evt := rawEvt.(type) {
	case *events.PairSuccess:
		fmt.Println("Pair Success!")
		//c.setLoginToken("DONE")
		c.setStatus("Connected")
	case *events.Connected:
		fmt.Println("Resuming session")
		//c.setLoginToken("DONE")
		c.setStatus("Connected")
	case *events.Message:
		msg := fmt.Sprintf("%s said %s to %s at %s\n", evt.Info.PushName, evt.Message.GetConversation(), evt.Info.Chat, evt.Info.Timestamp)
		fmt.Printf(msg)
		c.setMessage(msg)
		c.setStatus("Messaged")

	}
}

func main() {
	err := qml.Run(run)
	if err != nil {
		log.Fatal(err)
	}
}
