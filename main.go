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

var client *whatsmeow.Client
var status = make(chan string)
var qmlBridge = Window{LoginToken: ""}

type Window struct {
	Root       qml.Object
	LoginToken string
}

func run() error {
	engine := qml.NewEngine()
	component, err := engine.LoadFile("qml/Main.qml")
	if err != nil { return err }

	qmlBridge = Window{
		LoginToken: "",
	}
	context := engine.Context()
	context.SetVar("qmlBridge", &qmlBridge)

	win := component.CreateWindow(nil)
	qmlBridge.Root = win.Root()

    // FixMe: Gets stuck waiting for status most of the time
    // FixMe: Doesn't listen for variable changes most of the time
	go whatsApp()
    var _ = <- status

    win.Show()
	win.Wait()

	return nil
}

func handler(rawEvt interface{}) {
	switch evt := rawEvt.(type) {
        case *events.PairSuccess:
            fmt.Println("Pair Success!")
            qmlBridge.PushToQML("DONE")
            status <- "DONE"
        case *events.Connected:
            fmt.Println("Resuming session")
            qmlBridge.PushToQML("DONE")
            status <- "DONE"
	case *events.Message:
		fmt.Printf("%s said %s to %s at %s\n", evt.Info.PushName, evt.Message.GetConversation(), evt.Info.Chat, evt.Info.Timestamp)
	}
}

func whatsApp() {
    var dbPath, err = xdg.ConfigFile("whatsut.johndoe/userStore.db")
    if err != nil { panic(err) }
    os.Mkdir(strings.Replace(dbPath, "/userStore.db", "", 1), 0755)


	container, err := sqlstore.New("sqlite3", dbPath+"?_foreign_keys=on", waLog.Stdout("Database", "DEBUG", true))
	if err != nil { panic(err) }

	deviceStore, err := container.GetFirstDevice()
	if err != nil { panic(err) }

	client = whatsmeow.NewClient(deviceStore, waLog.Stdout("Client", "DEBUG", true))
	client.AddEventHandler(handler)

	if client.Store.ID == nil {
		qrChan, _ := client.GetQRChannel(goContext.Background())
		err = client.Connect()
		if err != nil { panic(err) }

		for evt := range qrChan {
			if evt.Event == "code" {
				fmt.Println("QR code:", evt.Code)

                // FixMe: The LoginToken can't be redefined from here or the handler usually, but not all of the time
                qmlBridge.PushToQML(evt.Code)
                status <- "QR"
			} else {
				fmt.Println("Login event:", evt.Event)
			}
		}
	} else {
		err = client.Connect()
		if err != nil { panic(err) }
	}
}

func (qmlBridge *Window) PushToQML(token string) {
	fmt.Printf("[whatsut] pushtoqml %s\n", token)
	qmlBridge.LoginToken = token
	qml.Changed(qmlBridge, &qmlBridge.LoginToken)
}

func main() {
	err := qml.Run(run)
	if err != nil { log.Fatal(err) }
}
