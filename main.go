package main

import (
	"bytes"
	"crypto/sha256"
	"crypto/tls"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Auth struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Token struct {
	Token string `json:"token"`
}

func main() {
	godotenv.Load()

	addr := os.Getenv("SERVER_ADDR")
	client := &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	token, err := getToken(client, addr)
	if err != nil {
		fmt.Printf("Ошибка: %s", err)
		return
	}

	err = getFile(client, addr, token)
	if err != nil {
		fmt.Printf("Ошибка: %s", err)
		return
	}

}

func getToken(client *http.Client, addr string) (string, error) {

	var auth Auth

	fmt.Print("Введите логин: ")
	fmt.Scan(&auth.Username)
	fmt.Print("Введите пароль: ")
	fmt.Scan(&auth.Password)

	body, err := json.Marshal(auth)
	if err != nil {
		return "", err
	}

	url := fmt.Sprintf("%s/auth", addr)

	resp, err := client.Post(url, "application/json", bytes.NewBuffer(body))

	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	var t Token

	if err := json.NewDecoder(resp.Body).Decode(&t); err != nil {
		return "", err
	}

	return t.Token, nil
}

func getFile(client *http.Client, addr string, token string) error {

	var version string
	var fileID string
	fmt.Print("Введите file id: ")
	fmt.Scan(&fileID)
	fmt.Print("Введите версию в формате X.Y.Z: ")
	fmt.Scan(&version)

	url := fmt.Sprintf("%s/files/%s/version/%s/download", addr, fileID, version)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	filename := "downloaded_file"
	if cd := resp.Header.Get("Content-Disposition"); cd != "" {

		idx := strings.Index(cd, "filename=")
		name := cd[idx+len("filename="):]
		if name != "" {
			filename = name
		}
	}

	expectedHash := resp.Header.Get("Digest")
	destPath := filepath.Join("downloads", filename)

	if err := os.MkdirAll("downloads", 0755); err != nil {
		return err
	}

	file, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return err
	}

	data, err := os.ReadFile(destPath)
	if err != nil {
		return err
	}

	actualHash := computeSHA256(data)
	if expectedHash != "" && actualHash != expectedHash {
		return fmt.Errorf("хэш не совпадает: ожидался %s, получен %s", expectedHash, actualHash)
	}

	fmt.Println("Файл прошел проверку целостности")

	return nil
}

func computeSHA256(data []byte) string {
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}
