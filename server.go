package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

var tokens = make(map[string]bool)

type dataJson struct {
	Login string
	Password string
	Token string
}

func answer(w http.ResponseWriter, r *http.Request) {
	var indexer dataJson

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&indexer)
	_ = r.Body.Close()
	if err != nil {
		fmt.Println("error:", err)
	}

	pwd := map[string]string{
		"kamazer": "1234",
		"dimax": "qwer",
	}

	if indexer.Login != "" && indexer.Password != "" {
		if pwd[indexer.Login] == indexer.Password {
			// Вошел по логину и паролю
			tok := time.Now().String()
			resp, _ := json.Marshal("Ваш токен для входа: " + tok)
			_, _ = w.Write(resp)
			// Записываем время на этого человека
			tokens[tok] = true
			fmt.Println("Новый токен: " + tok)
		}
	} else if indexer.Token != "" {
		if tokens[indexer.Token] {
			// Вошел по токену
			_, _ = w.Write([]byte("Вы вошли по токену"))
		}
	} else {
		_, _ = w.Write([]byte("Не удалось войти"))
	}

}

func main () {
	mux := http.NewServeMux()
	mux.HandleFunc("/", answer)
	err := http.ListenAndServe(":3000", mux)
	println(err)
}
