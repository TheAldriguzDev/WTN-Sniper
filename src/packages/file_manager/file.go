package file_manager

import (
	console_logger "WTN-Sniper/src/packages/logger"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"time"
)

/*
- Settings --> Save Username, email, password, webhook
- Data --> Save logs like order confirmed
- Task --> Save the shoes and input the min price that is acceptable
- Proxies --> Proxy file folder
*/
var folders = []string{"Proxies", "Data", "Settings"}

type SettingsData struct {
	SuccessWebhook string     `json:"success_webhook"`
	ErrorWebhook   string     `json:"error_webhook"`
	WTN            WtnAccount `json:"wtn_account"`
	Delay          int        `json:"delay"`
}

type WtnAccount struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

var l = console_logger.New("Dashboard", true)

func CheckSetup() bool {
	if !getDir() {
		l.Error("An error occured while getting the directories!")
		time.Sleep(20 * time.Second)
		os.Exit(1)
	}

	if !getFiles() {
		l.Error("An error occured while getting the settings.json!")
		time.Sleep(20 * time.Second)
		os.Exit(1)
	}

	return true
}

func getDir() bool {
	for _, folderName := range folders {
		folder, err := os.Open("./" + folderName)
		if err != nil {
			l.Error(fmt.Sprintf("Folder %s not found!", folderName))

			if createFolder(folderName) {
				l.Success(fmt.Sprintf("Folder %s has been created!", folderName))
			} else {
				l.Success(fmt.Sprintf("Error while creating %s folder! Open a ticket!", folderName))
				time.Sleep(20 * time.Second)
				os.Exit(1)
			}
		}
		defer folder.Close()
	}
	return true
}

func createFolder(folderName string) bool {
	err := os.Mkdir(folderName, 0755)
	if err != nil {
		return false
	}

	return true
}

func getFiles() bool {
	if _, err := os.Stat("./settings/settings.json"); err != nil {
		l.Error("Settings.json not found!")
		err := createSettings()
		if err != nil {
			l.Error("Error while creating the settings.json file! Open a ticket")
			time.Sleep(20 * time.Second)
			os.Exit(1)
		}

		l.Success("Settings.json file has been created! Please fill the data")
		time.Sleep(30 * time.Second)
		os.Exit(0)
	}

	return true
}

func createSettings() error {
	settData := SettingsData{
		SuccessWebhook: "",
		ErrorWebhook:   "",
		WTN: WtnAccount{
			Email:    "",
			Password: "",
		},
		Delay: 3,
	}

	jsonData, err := json.MarshalIndent(settData, "", "    ")

	if err != nil {
		return err
	}

	err = ioutil.WriteFile("./settings/settings.json", jsonData, 0644)
	if err != nil {
		return err
	}

	return nil
}

func ReadSettings() (SettingsData, error) {
	settingsData := SettingsData{}

	content, err := ioutil.ReadFile("./settings/settings.json")
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(content, &settingsData)
	if err != nil {
		panic(err)
	}

	return settingsData, nil
}

func ReadProxies(path string) (string, error) {
	folder, err := os.Open("./" + path)

	if err != nil {
		return "", nil
	}

	defer folder.Close()

	fileList, err := folder.ReadDir(0)

	if err != nil {
		return "", err
	}

	var filesName []string

	for _, file := range fileList {
		filesName = append(filesName, file.Name())
	}

	if len(filesName) == 0 {
		l.Error("No files found! Closing the bot...")
		time.Sleep(20 * time.Second)
		os.Exit(0)
	}

	var selection int
	var counter int
	for {
		l.General(">>> Please select a proxy file")
		counter = 0

		for _, elem := range filesName {
			l.General(fmt.Sprintf(">>> %s - %s", strconv.Itoa(counter), elem))
			counter++
		}

		fmt.Printf(console_logger.TextColors.White+"\r[%s] [%s] >>> "+console_logger.TextColors.Reset, "Dashboard", l.GetTime())

		if _, err := fmt.Scan(&selection); err != nil || selection <= -1 || selection > len(filesName)+1 {
			l.Error("Wrong input! Please retry!")
			time.Sleep(3 * time.Second)
			continue
		}
		break
	}

	return filesName[selection], nil
}
