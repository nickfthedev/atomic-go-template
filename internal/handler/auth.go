package handler

import "net/http"

func (h *Handler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Login"))
}

func (h *Handler) HandleSignup(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Signup"))
}

func (h *Handler) HandleForgetPassword(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Forget Password"))
}
