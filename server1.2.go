package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

var tokens = make(map[string]time.Time)
var refreshToken = make(map[string]string)
var timeLifeRefresh = make(map[string]time.Time)

type dataJson struct {
	Login string
	Password string
	Token string
	RefToken string
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
			tok := time.Now()
			resp, _ := json.Marshal("Ваш новый токен: " + tok.String() + " | Ваш рефрешь токен: " + tok.String())
			_, _ = w.Write(resp)
			// Записываем время на этого человека
			tokens[tok.String()] = tok
			refreshToken[tok.String()] = tok.String()
			timeLifeRefresh[tok.String()] = tok
			fmt.Println("Новый токен: " + tok.String() + " | Рефрешь токен:" + tok.String())
		}
	} else if indexer.Token != "" {
		if (time.Now()).Sub(tokens[indexer.Token]) <= time.Minute {
			// Вошел по токену
			_, _ = w.Write([]byte("Вы вошли по токену"))
		} else {
			// Проверка времени жизни рефрешь токена
			if (time.Now()).Sub(timeLifeRefresh[indexer.RefToken]) <= 2 * time.Minute {
				if refreshToken[indexer.RefToken] == indexer.Token {
					tok := time.Now()
					delete(tokens, indexer.Token)
					delete(refreshToken, indexer.RefToken)
					tokens[tok.String()] = tok
					refreshToken[indexer.RefToken] = tok.String()
					_, _ = w.Write([]byte("Новый токен:" + tok.String()))
				}
			} else {
					_, _ = w.Write([]byte("Рефрешь токен закончился"))
			}
			
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
