package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

var task string // Глобальная переменная для хранения значения "task"

// HelloHandler отправляет сообщение с содержимым переменной task
func GetHandler(w http.ResponseWriter, r *http.Request) {
	var tasks []Message

	// Извлекаем все записи из таблицы
	if err := DB.Find(&tasks).Error; err != nil {
		http.Error(w, "Не удалось получить задачи из базы данных.", http.StatusInternalServerError)
		return
	}

	// Отправляем список задач в формате JSON
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(tasks); err != nil {
		http.Error(w, "Ошибка при отправке данных.", http.StatusInternalServerError)
	}
}

// PostHandler обрабатывает POST-запрос, декодирует JSON и сохраняет значение в переменную task и в БД
func PostHandler(w http.ResponseWriter, r *http.Request) {
	var reqBody struct {
		Task   string `json:"task"`
		IsDone bool   `json:"is_done"`
	}

	// Декодируем JSON из тела запроса в структуру reqBody
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		http.Error(w, "Неверная нагрузка запроса.", http.StatusBadRequest)
		return
	}

	// Сохраняем задачу в базе данных
	message := Message{Task: reqBody.Task, IsDone: reqBody.IsDone}
	if err := DB.Create(&message).Error; err != nil {
		http.Error(w, "Не удалось сохранить задачу в базу данных.", http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(w, "Task updated successfully")
}

func UpdateHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Неверный формат ID", http.StatusBadRequest)
		return
	}
	var reqBody struct {
		Task   string `json:"task"`
		IsDone bool   `json:"is_done"`
	}
	err = json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		http.Error(w, "Неверная нагрузка запроса.", http.StatusBadRequest)
		return
	}
	var message Message
	if err := DB.First(&message, id).Error; err != nil {
		http.Error(w, "Задача не найдена.", http.StatusBadRequest)
		return
	}
	message.Task = reqBody.Task
	message.IsDone = reqBody.IsDone
	if err := DB.Save(&message).Error; err != nil {
		http.Error(w, "Не удалось обновить задачу.", http.StatusInternalServerError)
		return
	}
	fmt.Fprintln(w, "Task updated successfully")

}

func DeleteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Неверный формат ID.", http.StatusBadRequest)
		return
	}
	var message Message
	if err := DB.First(&message, id).Error; err != nil {
		http.Error(w, "Задача не найдена.", http.StatusBadRequest)
		return
	}
	if err := DB.Delete(&message).Error; err != nil {
		http.Error(w, "Не удалось удалить задачу.", http.StatusBadRequest)
		return
	}
	fmt.Fprintln(w, "Task deleted successfully")
}

func main() {
	// Инициализируем базу данных
	InitDB()

	// Выполняем миграцию структуры Message
	if err := DB.AutoMigrate(&Message{}); err != nil {
		log.Fatal("Не удалось применить миграцию базы данных:", err)
	}

	router := mux.NewRouter()

	// Регистрируем обработчики для POST и GET
	router.HandleFunc("/api/task", GetHandler).Methods("GET")
	router.HandleFunc("/api/task", PostHandler).Methods("POST")
	router.HandleFunc("/api/task/{id:[0-9]+}", UpdateHandler).Methods("PUT")    // PUT для обновления задачи
	router.HandleFunc("/api/task/{id:[0-9]+}", DeleteHandler).Methods("DELETE") // DELETE для удаления задачи

	// Запускаем сервер на localhost:8080
	fmt.Println("Server is listening on port 8080...")
	http.ListenAndServe(":8080", router)
}
