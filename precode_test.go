package gotestify

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// проверяем, что приходит корректный код ответа, а именно status code 200 и тело ответа не пустое
func TestMainHandlerWhenStatusCodeOk(t *testing.T) {
	req := httptest.NewRequest("GET", "/cafe?count=1&city=moscow", nil)

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	require.Equal(t, http.StatusOK, responseRecorder.Code, "Expected status code OK")

	body := responseRecorder.Body.String()
	require.NotEmpty(t, body, "body response must not be empty")
}

// проверяем, что город, который передается в параметре city, не поддерживается
// сервис возвращает код ответа 400 и ошибку wrong city value в теле ответа
func TestMainHandlerWhenWrongCityPointed(t *testing.T) {
	req := httptest.NewRequest("GET", "/cafe?count=1&city=someCity", nil)

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	assert.Equal(t, http.StatusBadRequest, responseRecorder.Code, "Expected status Bad Request")

	body := responseRecorder.Body.String()
	assert.Equal(t, "wrong city value", body, "message error expected")
}

// проверяем, что в параметре count указано больше, чем есть всего и должны вернуться все доступные кафе
func TestMainHandlerWhenCountMoreThanTotal(t *testing.T) {
	totalCount := 4

	req := httptest.NewRequest("GET", "/cafe?count=5&city=moscow", nil)

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	body := responseRecorder.Body.String()
	list := strings.Split(body, ",")
	assert.Len(t, list, totalCount, "Expected cafe count %d", totalCount)
}
