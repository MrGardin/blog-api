package handler

import (
	"blog-api/repository"
	"database/sql"
	"net/http"
)

func HomeHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userRep := repository.NewUserRepository(db)

		if err := userRep.CheckConnection(); err != nil {
			w.Write([]byte("Ошибка проверки пинга"))
			return
		}

		msg, err := userRep.GetWelcomeMessage()
		if err != nil {
			w.Write([]byte("Ошибка приветствия"))
			return
		}

		w.Write([]byte(msg))
	}
}
