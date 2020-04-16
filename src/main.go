package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"

	"github.com/andybrewer/mack"
)

// Path of erxes repos
const (
	Erxes             = "path/to/erxes"
	ErxesAPI          = "path/to/erxes-api"
	ErxesIntegrations = "path/to/erxes-integrations"
)

func loadInitialData() string {
	return fmt.Sprintf("%s", ErxesAPI)
}

func startErxesScript() string {
	return fmt.Sprintf(`
		tell application "iTerm"
			tell current window
					-- create a tab for erxes/erxes
					create tab with default profile
					tell current session
							write text "cd %s"
							write text "yarn start"
							-- split tab vertically to run scheduler
							split vertically with default profile
					end tell

					-- create tab for erxes/erxes-api
					tell last session of last tab
							write text "cd %s"
							write text "yarn dev"
							-- split tab vertically to run scheduler
							split vertically with default profile
					end tell

					-- create tab for erxes/erxes-integrations
					tell last session of last tab
							write text "cd %s"
							write text "yarn dev"
					end tell

					-- start redis-server
					create tab with default profile
					tell current session
							write text "redis-server"
					end tell

					-- show notification
					display notification "Let's get some something done" with title "Enjoy your day!" subtitle "💪💪💪"

			end tell
		end tell
	`, Erxes, ErxesAPI, ErxesIntegrations)
}

func getCurrentDate() (string, time.Month, int, int) {
	now := time.Now()

	formattedDate := now.Format("2006-01-02 15:04:05")

	currentYear, currentMonth, _ := now.Date()
	currentLocation := now.Location()

	firstOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)
	lastOfMonth := firstOfMonth.AddDate(0, 1, -1).Day()

	return formattedDate, currentMonth, now.Day(), lastOfMonth
}

func getHostname() string {
	name, err := os.Hostname()

	if err != nil {
		panic(err)
	}

	return name
}

func sendNotification(title string, description string, subtitle string) {
	mack.Notify(description, title, subtitle, "Glass")
}

func greetings() {
	hostname := getHostname()

	sendNotification("Hey "+hostname, "", "Let's rock something today 🦾")
}

func checkMonthlyReport() {
	hostname := getHostname()

	_, _, currentDay, lastDayOfMonth := getCurrentDate()

	if currentDay == 14 || currentDay == (lastDayOfMonth-1) {
		sendNotification("Hey 👋"+hostname, "You did your best 🎉🎉🎉", "Today is monthly report day, Goodluck")
	}
}

func checkDailyStandUp() {
	sendNotification("Psst 🙋🙋‍♀️", "Let's write it right away", "Don't forget the daily report")

	_, err := mack.Tell("Notes", createNote())
	if err != nil {
		panic(err)
	}

	browsers := []string{"Google Chrome", "Firefox", "Safari"}
	opened := false

	for _, browser := range browsers {
		_, err := mack.Tell(browser, `open location "https://trello.com/"`)
		if err != nil {
			fmt.Println("Error occured while trying to open browser")
		} else {
			// exit when we found a browser that works
			opened = true
			break
		}
	}

	if !opened {
		panic("Error occured while trying to open browser")
	}
}

func createNote() string {
	now := time.Now()

	currentYear, currentMonth, _ := now.Date()

	today := now.Day()
	parsedMonth := int(currentMonth)

	body := fmt.Sprintf(`
		[%d/%d/%d] What did you do yesterday?
		[%d/%d/%d] What will you do today?
	`, parsedMonth, today-1, currentYear, parsedMonth, today, currentYear)

	noteAppleScript := fmt.Sprintf(`
		tell application "Notes"
			activate
			tell account "iCloud"
				make new note at folder "Notes" with properties {name:"StandUp", body: "%s"}
			end tell
		end tell
	`, body)

	return noteAppleScript
}

func executeSelectedUtil(SelectedRow int) {
	switch SelectedRow {
	case 0:
		_, err := mack.Tell("iTerm", startErxesScript())
		if err != nil {
			panic(err)
		}
	case 1:
		cmd := exec.Command("mongo", "./scripts/clearDb.js")
		err := cmd.Run()
		if err != nil {
			panic(err)
		}

		renderLog("Successfully removed erxes dbs")
	case 2:
		cmd := exec.Command("mongo", "./scripts/removeAllDb.js")
		err := cmd.Run()
		if err != nil {
			panic(err)
		}

		renderLog("Successfully removed all dbs")
	}
}

func renderHeader() {
	currentTime, _, _, _ := getCurrentDate()

	header := widgets.NewParagraph()
	header.Text = "[Press q to quit ](fg:green,mod:bold)" + currentTime
	header.SetRect(0, 0, 70, 3)
	header.BorderStyle.Fg = ui.ColorCyan

	ui.Render(header)
}

func renderBody() {
	hostname := getHostname()

	intro := widgets.NewParagraph()
	intro.Title = "Hello " + hostname
	intro.Text = "[Erxes automation tool for those who have passionate for open source contribution](fg:blue,mod:bold) 💻"
	intro.BorderStyle.Fg = ui.ColorYellow
	intro.SetRect(0, 3, 70, 20)

	ui.Render(intro)
}

func renderFeatures() {
	features := widgets.NewList()
	features.Title = "Features"
	features.Rows = []string{
		"* Notify daily StandUp and automatically create note",
		"* Check monthly report and notify",
	}

	features.TextStyle = ui.NewStyle(ui.ColorYellow)
	features.WrapText = false
	features.Border = false
	features.SetRect(1, 7, 55, 11)

	ui.Render(features)
}

func renderList() {
	list := widgets.NewList()
	list.Title = "Utils [Press enter to execute command]"
	list.Rows = []string{
		"[0] start erxes project",
		"[1] mongo remove erxes, erxes-integrations db",
		"[2] mongo remove all db",
	}

	list.TextStyle = ui.NewStyle(ui.ColorYellow)
	list.WrapText = false
	list.Border = false
	list.SetRect(1, 11, 50, 18)

	ui.Render(list)

	previousKey := ""

	uiEvents := ui.PollEvents()

	for {
		e := <-uiEvents
		switch e.ID {
		case "q", "<C-c>":
			return
		case "<Enter>":
			executeSelectedUtil(list.SelectedRow)
		case "j", "<Down>":
			list.ScrollDown()
		case "k", "<Up>":
			list.ScrollUp()
		case "<C-d>":
			list.ScrollHalfPageDown()
		case "<C-u>":
			list.ScrollHalfPageUp()
		case "<C-f>":
			list.ScrollPageDown()
		case "<C-b>":
			list.ScrollPageUp()
		case "g":
			if previousKey == "g" {
				list.ScrollTop()
			}
		case "<Home>":
			list.ScrollTop()
		case "G", "<End>":
			list.ScrollBottom()
		}

		if previousKey == "g" {
			previousKey = ""
		} else {
			previousKey = e.ID
		}

		ui.Render(list)
	}
}

func renderMessage(text string) {
	message := widgets.NewParagraph()
	message.Text = text
	message.TextStyle = ui.NewStyle(ui.ColorBlack)
	message.TextStyle.Bg = ui.ColorGreen
	message.SetRect(0, 35, 34, 20)
	message.Border = false

	ui.Render(message)
}

func renderLog(text string) {
	log := widgets.NewParagraph()
	log.Text = "[LOG] => " + text
	log.TextStyle = ui.NewStyle(ui.ColorWhite)
	log.BorderStyle.Fg = ui.ColorWhite
	log.SetRect(0, 20, 70, 23)

	ui.Render(log)
}

func initUI() {
	if err := ui.Init(); err != nil {
		log.Fatalf("Failed to initialize termui: %v", err)
	}

	renderHeader()
	renderBody()
	renderFeatures()
	renderList()
}
func main() {
	greetings()

	checkDailyStandUp()
	checkMonthlyReport()

	initUI()

	defer ui.Close()
}
