package storage

import (
	"log"
	"os"
	"path"
	"runtime"
)

const (
	FOLDER_NAME = "galactic-clicker"
	FILE_NAME   = "data.json"
)

func fullConfigPath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		log.Println(err)
		return "", err
	}

	return path.Join(configDir, FOLDER_NAME, FILE_NAME), nil
}

func makeConfigDir() error {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return err
	}

	dir := path.Join(configDir, FOLDER_NAME)
	err = os.MkdirAll(dir, 0755)
	if err != nil {
		return err
	}

	return nil
}

func LoadBytes() ([]byte, error) {
	if runtime.GOARCH == "wasm" {
		log.Fatal("WASM not supported")
	}

	path, err := fullConfigPath()
	if err != nil {
		return nil, err
	}

	return os.ReadFile(path)
}

func SaveBytes(bytes []byte) {
	if runtime.GOARCH == "wasm" {
		log.Fatal("WASM not supported")
	}

	path, err := fullConfigPath()
	if err != nil {
		log.Fatal(err)
	}

	err = makeConfigDir()
	if err != nil {
		log.Fatal(err)
	}

	err = os.WriteFile(path, bytes, 0666)
	if err != nil {
		log.Fatal(err)
	}
}
