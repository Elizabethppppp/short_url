package main

import (
	"encoding/json"

	server "github.com/Elizabethppppp/tcp_server"
)

type ErrorResponse struct {
	Error      string `yaml:"error"`
	Message    string `yaml:"message"`
	MessageRUS string `yaml:"message_rus"`
}

var errorResponseMessage = map[int]struct {
	En  string
	Rus string
}{
	200: {"OK", "OK"},
	201: {"Created", "Создано"},
	302: {"Found", "Найдено"},
	400: {"Bad Request", "Неверный запрос"},
	403: {"Forbidden", "Доступ запрещён"},
	404: {"Not Found", "Не найдено"},
	405: {"Method Not Allowed", "Метод не разрешён"},
	418: {"I'm a teapot", "Я чайник"},
	429: {"Too Many Requests", "Слишком много запросов"},
	500: {"Internal Server Error", "Внутренняя ошибка сервера"},
	502: {"Bad Gateway", "Плохой шлюз"},
	503: {"Service Unavailable", "Сервис недоступен"},
}

func ResponseJSON(w server.ResponseWriter, status int, codeError error) {
	w.SetHeader("Content-Type", "application/json")
	w.WriteHeader(status)

	var code string
	if codeError != nil {
		code = codeError.Error()
	}

	msg, ok := errorResponseMessage[status]
	if !ok {
		msg = struct {
			En  string
			Rus string
		}{En: "Unknown", Rus: "Неизвестная ошибка"}
	}

	response := ErrorResponse{
		Error:      code,
		Message:    msg.En,
		MessageRUS: msg.Rus,
	}

	body, err := json.Marshal(response)
	if err != nil {
		w.Write([]byte(`{"code":"","message":"failed to marshal error","message_rus":"не удалось представитьв нужном формате"}`))
		return
	}

	w.Write(body)
}
