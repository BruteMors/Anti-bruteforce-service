package handlers

import (
	"Anti-bruteforce-service/internal/config"
	"Anti-bruteforce-service/internal/domain/entity"
	"Anti-bruteforce-service/internal/domain/service"
	mock_service "Anti-bruteforce-service/internal/store/postgressql/adapters/mocks"
	"bytes"
	"encoding/json"
	"github.com/golang/mock/gomock"
	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestTryAuthorization(t *testing.T) {
	logger := zap.NewExample().Sugar()

	controller := gomock.NewController(t)
	defer controller.Finish()

	blackListMockStore := mock_service.NewMockBlackListStore(controller)
	blackList := service.NewBlackList(blackListMockStore, logger)

	whiteListMockStore := mock_service.NewMockWhiteListStore(controller)
	whiteList := service.NewWhiteList(whiteListMockStore, logger)

	cfg, err := config.New()
	if err != nil {
		require.NoError(t, err)
	}

	serviceAuth := service.NewAuthorization(blackList, whiteList, cfg, logger)

	authorization := NewAuthorization(serviceAuth, logger)

	cases := []struct {
		name    string
		request entity.Request
	}{
		{name: "valid request", request: entity.Request{
			Login:    "test",
			Password: "1234",
			Ip:       "192.1.5.1",
		}},
	}

	blackListMockStore.EXPECT().Get().Return([]entity.IpNetwork{}, nil).AnyTimes()
	whiteListMockStore.EXPECT().Get().Return([]entity.IpNetwork{}, nil).AnyTimes()

	router := httprouter.New()
	router.POST("/auth/check", authorization.TryAuthorization)

	request := cases[0].request

	body, err := json.Marshal(request)
	require.NoError(t, err)

	req, err := http.NewRequest("POST", "/auth/check", bytes.NewReader(body))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	s := rr.Body.String()
	require.Equal(t, "ok=true", s)

}
