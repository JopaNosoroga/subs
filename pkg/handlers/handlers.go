package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"subscriptions/pkg/dbwork"
	"subscriptions/pkg/models"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// AddSubscriptions обрабатывает создание новой подписки
// swagger:operation POST /subscriptions subscriptions createSubscription
// ---
// summary: Создать новую подписку
// description: Добавляет новую подписку в систему
// parameters:
//   - name: body
//     in: body
//     description: Данные подписки
//     required: true
//     schema:
//     $ref: "#/definitions/subscription"
//
// responses:
//
//	201:
//	  description: Подписка успешно создана
//	400:
//	  description: Неверный формат запроса
//	  schema:
//	    $ref: "#/definitions/errorResponse"
//	500:
//	  description: Ошибка сервера
//	  schema:
//	    $ref: "#/definitions/errorResponse"
func AddSubscriptions(rw http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(rw)
	data, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		encoder.Encode(models.ErrorResponse{
			Message: "Ошибка сервера",
			Code:    500,
		})
		return
	}

	sub := models.Subscriptions{}
	err = json.Unmarshal(data, &sub)
	if err != nil {
		log.Println(err)
		encoder.Encode(models.ErrorResponse{
			Message: "Ошибка чтеняи запроса",
			Code:    400,
		})
		return
	}

	err = dbwork.DB.AddSubscriptions(sub)
	if err != nil {
		log.Println(err)
		encoder.Encode(models.ErrorResponse{
			Message: "Ошибка добавления подписки",
			Code:    500,
		})
		return
	}

	rw.WriteHeader(http.StatusCreated)
}

// DeleteSubscriptions обрабатывает удаление подписки
// swagger:operation DELETE /subscriptions/{id} subscriptions deleteSubscription
// ---
// summary: Удалить подписку
// description: Удаляет подписку по идентификатору
// parameters:
//   - name: id
//     in: path
//     description: ID подписки
//     required: true
//     type: integer
//     format: int64
//
// responses:
//
//	200:
//	  description: Подписка успешно удалена
//	400:
//	  description: Неверный ID подписки
//	  schema:
//	    $ref: "#/definitions/errorResponse"
//	500:
//	  description: Ошибка сервера
//	  schema:
//	    $ref: "#/definitions/errorResponse"
func DeleteSubscriptions(rw http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(rw)
	vars := mux.Vars(r)
	strID := vars["id"]

	id, err := strconv.Atoi(strID)
	if err != nil {
		log.Println(err)
		encoder.Encode(models.ErrorResponse{
			Message: "Ошибка запроса",
			Code:    400,
		})
		return
	}

	err = dbwork.DB.DeleteSubscriptions(id)
	if err != nil {
		log.Println(err)
		encoder.Encode(models.ErrorResponse{
			Message: "Ошибка удаления",
			Code:    500,
		})
		return
	}

	rw.WriteHeader(http.StatusOK)
}

// UpdateSubscriptions обрабатывает обновление подписки
// swagger:operation PUT /subscriptions/{id} subscriptions updateSubscription
// ---
// summary: Обновить подписку
// description: Обновляет данные существующей подписки
// parameters:
//   - name: id
//     in: path
//     description: ID подписки
//     required: true
//     type: integer
//     format: int64
//   - name: body
//     in: body
//     description: Обновленные данные подписки
//     required: true
//     schema:
//     $ref: "#/definitions/subscription"
//
// responses:
//
//	200:
//	  description: Подписка успешно обновлена
//	400:
//	  description: Неверный запрос
//	  schema:
//	    $ref: "#/definitions/errorResponse"
//	500:
//	  description: Ошибка сервера
//	  schema:
//	    $ref: "#/definitions/errorResponse"
func UpdateSubscriptions(rw http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(rw)
	vars := mux.Vars(r)
	strID := vars["id"]

	id, err := strconv.Atoi(strID)
	if err != nil {
		log.Println(err)
		encoder.Encode(models.ErrorResponse{
			Message: "Ошибка запроса",
			Code:    400,
		})
		return
	}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		encoder.Encode(models.ErrorResponse{
			Message: "Ошибка сервера",
			Code:    500,
		})
		return
	}

	sub := models.Subscriptions{}

	err = json.Unmarshal(data, &sub)
	if err != nil {
		log.Println(err)
		encoder.Encode(models.ErrorResponse{
			Message: "Ошибка чтения запроса",
			Code:    500,
		})
		return
	}

	err = dbwork.DB.UpdateSubscriptions(id, sub)
	if err != nil {
		log.Println(err)
		encoder.Encode(models.ErrorResponse{
			Message: "Ошибка изменения подписки",
			Code:    500,
		})
		return
	}

	rw.WriteHeader(http.StatusOK)
}

// GetSubscriptions обрабатывает получение информации о подписке
// swagger:operation GET /subscriptions/{id} subscriptions getSubscription
// ---
// summary: Получить информацию о подписке
// description: Возвращает информацию о конкретной подписке
// parameters:
//   - name: id
//     in: path
//     description: ID подписки
//     required: true
//     type: integer
//     format: int64
//
// responses:
//
//	200:
//	  description: Успешный запрос
//	  schema:
//	    $ref: "#/definitions/subscription"
//	404:
//	  description: Подписка не найдена
//	  schema:
//	    $ref: "#/definitions/errorResponse"
//	500:
//	  description: Ошибка сервера
//	  schema:
//	    $ref: "#/definitions/errorResponse"
func GetSubscriptions(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	strID := vars["id"]
	encoder := json.NewEncoder(rw)

	id, err := strconv.Atoi(strID)
	if err != nil {
		log.Println(err)
		encoder.Encode(models.ErrorResponse{
			Message: "Ошибка запроса",
			Code:    400,
		})
		return
	}

	sub, err := dbwork.DB.GetSubscriptions(id)
	if err != nil {
		log.Println(err)
		encoder.Encode(models.ErrorResponse{
			Message: "Ошибка получения подписки",
			Code:    500,
		})
		return
	}

	if sub.UserID == uuid.Nil {
		encoder.Encode(models.ErrorResponse{
			Message: "Нет такой записи",
			Code:    400,
		})
		return
	}

	err = encoder.Encode(sub)
	if err != nil {
		log.Println(err)
		encoder.Encode(models.ErrorResponse{
			Message: "Ошибка отправки подписки",
			Code:    500,
		})
		return
	}
}

// GetAllSubscriptions обрабатывает получение списка подписок
// swagger:operation GET /subscriptions subscriptions listSubscriptions
// ---
// summary: Получить список подписок
// description: Возвращает список подписок с возможностью фильтрации
// parameters:
//   - name: body
//     in: body
//     description: Параметры фильтрации (все поля опциональны)
//     schema:
//     $ref: "#/definitions/subscription"
//
// responses:
//
//	200:
//	  description: Успешный запрос
//	  schema:
//	    type: array
//	    items:
//	      $ref: "#/definitions/subscription"
//	500:
//	  description: Ошибка сервера
//	  schema:
//	    $ref: "#/definitions/errorResponse"
func GetAllSubscriptions(rw http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(rw)
	data, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		encoder.Encode(models.ErrorResponse{
			Message: "Ошибка сервера",
			Code:    500,
		})
		return
	}

	sub := models.Subscriptions{}

	err = json.Unmarshal(data, &sub)
	if err != nil {
		log.Println(err)
		encoder.Encode(models.ErrorResponse{
			Message: "Ошибка чтения запросов",
			Code:    500,
		})
		return
	}

	subs, err := dbwork.DB.GetAllSubscriptions(sub)
	if err != nil {
		log.Println(err)
		encoder.Encode(models.ErrorResponse{
			Message: "Ошибка получения подписок",
			Code:    500,
		})
		return
	}

	err = encoder.Encode(subs)
	if err != nil {
		log.Println(err)
		encoder.Encode(models.ErrorResponse{
			Message: "Ошибка отправки подписок",
			Code:    500,
		})
		return
	}
}

// CalculatePrice рассчитывает сумму подписок за период
// swagger:operation POST /sum subscriptions calculatePrice
// ---
// summary: Рассчитать стоимость подписок
// description: Возвращает сумму цен активных подписок за указанный период
// parameters:
//   - name: body
//     in: body
//     description: Параметры расчета (обязательны start_date и end_date)
//     required: true
//     schema:
//     $ref: "#/definitions/subscription"
//
// responses:
//
//	200:
//	  description: Успешный расчет
//	  schema:
//	    type: integer
//	    example: 3500
//	400:
//	  description: Неверные параметры запроса
//	  schema:
//	    $ref: "#/definitions/errorResponse"
//	500:
//	  description: Ошибка сервера
//	  schema:
//	    $ref: "#/definitions/errorResponse"
func CalculatePrice(rw http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(rw)
	data, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		encoder.Encode(models.ErrorResponse{
			Message: "Ошибка сервера",
			Code:    500,
		})
		return
	}

	sub := models.Subscriptions{}

	err = json.Unmarshal(data, &sub)
	if err != nil {
		log.Println(err)
		encoder.Encode(models.ErrorResponse{
			Message: "Ошибка чтения запроса",
			Code:    500,
		})
		return
	}

	startDate := sub.StartDate
	endDate := sub.EndDate
	sub.EndDate = time.Time{}
	sub.StartDate = time.Time{}

	subs, err := dbwork.DB.GetAllSubscriptions(sub)
	if err != nil {
		log.Println(err)
		encoder.Encode(models.ErrorResponse{
			Message: "Ошибка получения подписок",
			Code:    500,
		})
		return
	}

	sum := 0

	for _, val := range subs {
		if val.StartDate.After(startDate) && val.StartDate.Before(endDate) {
			sum += val.Price
		}
	}

	err = encoder.Encode(sum)
	if err != nil {
		log.Println(err)
		encoder.Encode(models.ErrorResponse{
			Message: "Ошибка отправки суммы",
			Code:    500,
		})
		return
	}
}
