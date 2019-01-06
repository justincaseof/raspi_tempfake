package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"syscall"
	"time"
)

type StateInfo struct {
	LastQuery   time.Time `json:"lastQuery"`
	LastRandVal float32   `json:"lastRandVal"`
}

var infoFileName = "info.json"

func LoadInfo(stateInfoFileName string) (StateInfo, error) {
	var info StateInfo
	infoFile, err := os.Open(stateInfoFileName)
	defer infoFile.Close()
	if err != nil {
		//fmt.Println("No StateInfo available. Creating new.")
		// no file? try to create one!
		newFile, err := os.Create(stateInfoFileName)

		// still no success? fail hard!
		if err != nil {
			fmt.Println("No StateInfo found and neither able to create one. Giving up.")
			fmt.Println(err.Error())
			return info, nil
		}
		defer newFile.Close()

		// create random initial number with 15 <= number <= 85
		info.LastRandVal = rand.Float32()*70 + 15
	} else {
		// read from existing file
		jsonParser := json.NewDecoder(infoFile)
		jsonParser.Decode(&info)
	}

	return info, nil
}

func WriteInfo(info StateInfo, stateInfoFileName string) error {
	var err error = nil
	var infoFile *os.File

	infoFile, err = os.OpenFile(stateInfoFileName, os.O_RDWR, syscall.FILE_MAP_WRITE)
	defer infoFile.Close()
	if err != nil {
		return err
	}

	encoder := json.NewEncoder(infoFile)
	err = encoder.Encode(&info)

	return err
}

func main() {
	// seed in order to get "actual random" numbers
	rand.Seed(time.Now().UnixNano())

	randBool := (rand.Int()%2 == 0)
	randOffset := rand.Float32() * 2

	var err error = nil
	var info StateInfo

	// LOAD
	info, err = LoadInfo(infoFileName)
	if err != nil {
		fmt.Println("Cannot write file!", err.Error())
		os.Exit(1)
	}

	// EXEC
	// remember current time.
	info.LastQuery = time.Now()
	if randBool {
		info.LastRandVal += randOffset
	} else {
		info.LastRandVal -= randOffset
	}

	// BOUNDARIES! (randomly chosen)
	if info.LastRandVal < -10 {
		info.LastRandVal = -10
	}
	if info.LastRandVal > 120 {
		info.LastRandVal = 120
	}

	// WRITE
	err = WriteInfo(info, infoFileName)
	if err != nil {
		fmt.Println("Cannot write file!", err.Error())
		os.Exit(1)
	}

	// PRINT
	fmt.Printf("temp=%.1f'C\n", info.LastRandVal)
}
