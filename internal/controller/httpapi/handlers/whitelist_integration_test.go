package handlers

import (
	"Anti-bruteforce-service/internal/domain/entity"
	"Anti-bruteforce-service/internal/domain/service"
	mock_service "Anti-bruteforce-service/internal/store/postgressql/adapters/mocks"
	"bytes"
	"encoding/json"
	"errors"
	"github.com/golang/mock/gomock"
	jsoniter "github.com/json-iterator/go"
	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWhiteList_AddIP(t *testing.T) {
	// Создаем экземпляр объекта WhiteList
	logger := zap.NewExample().Sugar()

	controller := gomock.NewController(t)
	defer controller.Finish()
	mockStore := mock_service.NewMockWhiteListStore(controller)

	cases := []struct {
		name    string
		network entity.IpNetwork
	}{
		{name: "valid ip and mask", network: entity.IpNetwork{
			Ip:   "192.168.1.1",
			Mask: "255.255.255.0",
		}},
		{name: "invalid ip", network: entity.IpNetwork{
			Ip:   "192.12.256.1",
			Mask: "255.255.255.0",
		}},
	}

	for _, testCase := range cases {
		prefix, err := service.GetPrefix(testCase.network.Ip, testCase.network.Mask)
		require.NoError(t, err)
		mockStore.EXPECT().Add(prefix, testCase.network.Mask).Return(nil).MaxTimes(1)
		mockStore.EXPECT().Add(prefix, testCase.network.Mask).Return(errors.New("this ip network already exist")).AnyTimes()
	}

	whiteListService := service.NewWhiteList(mockStore, logger)
	whitelist := NewWhiteList(whiteListService, logger)

	// Создаем тестовый HTTP-сервер
	router := httprouter.New()
	router.POST("/auth/whitelist", whitelist.AddIP)

	// Создаем тестовый IP для добавления в белый список
	ip := cases[0].network

	// Кодируем IP в формат JSON
	body, err := json.Marshal(ip)
	require.NoError(t, err)

	// Отправляем POST-запрос на сервер для добавления IP в белый список
	req, err := http.NewRequest("POST", "/auth/whitelist", bytes.NewReader(body))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	// Проверяем, что сервер вернул статус-код 204 No Content
	require.Equal(t, http.StatusNoContent, rr.Code)

	// Попытка добавления уже существующего IP в белый список
	req, err = http.NewRequest("POST", "/auth/whitelist", bytes.NewReader(body))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	// Проверяем, что сервер вернул статус-код 400 Bad Request и сообщение об ошибке
	require.Equal(t, http.StatusBadRequest, rr.Code)
	expectedBody := "this ip network already exist"
	require.Equal(t, expectedBody, rr.Body.String())

	// Попытка добавления невалидного IP в черный список
	invalidIP := cases[1].network
	body, err = json.Marshal(invalidIP)
	require.NoError(t, err)

	req, err = http.NewRequest("POST", "/auth/whitelist", bytes.NewReader(body))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	// Проверяем, что сервер вернул статус-код 400 Bad Request
	require.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestWhiteList_RemoveIP(t *testing.T) {
	// Создаем экземпляр объекта WhiteList
	logger := zap.NewExample().Sugar()

	controller := gomock.NewController(t)
	defer controller.Finish()
	mockStore := mock_service.NewMockWhiteListStore(controller)

	cases := []struct {
		name    string
		network entity.IpNetwork
	}{
		{name: "valid ip and mask", network: entity.IpNetwork{
			Ip:   "192.168.1.1",
			Mask: "255.255.255.0",
		}},
		{name: "invalid ip", network: entity.IpNetwork{
			Ip:   "192.12.256.1",
			Mask: "255.255.255.0",
		}},
	}

	for _, testCase := range cases {
		prefix, err := service.GetPrefix(testCase.network.Ip, testCase.network.Mask)
		require.NoError(t, err)
		mockStore.EXPECT().Remove(prefix, testCase.network.Mask).Return(nil).AnyTimes()
	}

	whiteListService := service.NewWhiteList(mockStore, logger)
	whitelist := NewWhiteList(whiteListService, logger)

	// Создаем тестовый HTTP-сервер
	router := httprouter.New()
	router.DELETE("/auth/whitelist", whitelist.RemoveIP)

	// Создаем тестовый IP для удаления из белого списка
	ip := cases[0].network

	// Кодируем IP в формат JSON
	body, err := json.Marshal(ip)
	require.NoError(t, err)

	// Отправляем POST-запрос на сервер для добавления IP в белый список
	req, err := http.NewRequest("DELETE", "/auth/whitelist", bytes.NewReader(body))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	// Проверяем, что сервер вернул статус-код 204 No Content
	require.Equal(t, http.StatusNoContent, rr.Code)

	// Попытка удаления невалидного IP из черного списка
	invalidIP := cases[1].network
	body, err = json.Marshal(invalidIP)
	require.NoError(t, err)

	req, err = http.NewRequest("DELETE", "/auth/whitelist", bytes.NewReader(body))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	// Проверяем, что сервер вернул статус-код 400 Bad Request
	require.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestWhiteList_ShowIPList(t *testing.T) {
	// Создаем экземпляр объекта BlackList
	logger := zap.NewExample().Sugar()

	controller := gomock.NewController(t)
	defer controller.Finish()
	mockStore := mock_service.NewMockWhiteListStore(controller)

	cases := []entity.IpNetwork{{
		Ip:   "192.168.1.1",
		Mask: "255.255.255.0",
	}, {
		Ip:   "192.168.2.1",
		Mask: "255.255.255.0",
	},
	}

	mockStore.EXPECT().Get().Return(cases, nil).AnyTimes()

	whiteListService := service.NewWhiteList(mockStore, logger)
	whitelist := NewWhiteList(whiteListService, logger)

	// Создаем тестовый HTTP-сервер
	r := httprouter.New()
	r.GET("/auth/whitelist", whitelist.ShowIPList)
	ts := httptest.NewServer(r)
	defer ts.Close()

	// Отправляем GET-запрос к тестовому серверу
	res, err := http.Get(ts.URL + "/auth/whitelist")
	require.NoError(t, err)
	defer res.Body.Close()

	// Проверяем код ответа
	assert.Equal(t, http.StatusOK, res.StatusCode)

	// Проверяем содержимое ответа
	var ipList []entity.IpNetwork
	err = jsoniter.NewDecoder(res.Body).Decode(&ipList)
	require.NoError(t, err)
	assert.Equal(t, cases, ipList)
}
