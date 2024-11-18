package httphandlers

import (
	"GophKeeper/internal/app/HTTP/middlewares"
	"io"
	"net/http"
	"strconv"
)

func (h *handlerHTTP) PasswordGet(w http.ResponseWriter, r *http.Request) {
	//get userID from ctx
	userID := r.Context().Value(middlewares.UserIDContextKey)
	userIDInt, ok := userID.(int)
	if !ok || userIDInt <= 0 {
		h.Logger.Warnf("unauthenticated request")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	//parse data
	login, err := io.ReadAll(r.Body)
	if err != nil {
		h.Logger.Errorf("cannot read body: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	loginStr := string(login)

	if len(loginStr) == 0 {
		h.Logger.Debug("login is empty")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//get login and encryptedPassword
	encryptedPassword, dataID, err := h.Storage.GetPasswordByLogin(r.Context(), userIDInt, loginStr)
	if err != nil {
		h.Logger.Errorf("cannot get login and encryptedPassword")
		h.Logger.Debugf("cannot get login and encryptedPassword, err: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//read encryption key
	key, err := h.KeyKeeper.GetLoginAndPasswordKey(strconv.Itoa(userIDInt), strconv.Itoa(dataID))
	if err != nil {
		h.Logger.Errorf("cant get encryption key from key storage")
		h.Logger.Debugf("cant get encryption key from key storage, err: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//decrypt
	passwordBytes, err := h.Encryptor.DecryptAESGCM([]byte(encryptedPassword), []byte(key))
	if err != nil {
		h.Logger.Errorf("cannot decrypt password")
		h.Logger.Debugf("cannot decrypt password, err: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//return encryptedPassword
	w.WriteHeader(http.StatusOK)
	w.Write(passwordBytes)
}
