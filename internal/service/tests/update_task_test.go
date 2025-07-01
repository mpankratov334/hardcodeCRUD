package service

import (
	"TemplatestPGSQL/internal/dto"
	"TemplatestPGSQL/internal/repo/mocks"
	"TemplatestPGSQL/internal/service"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
	"io"
	"net/http/httptest"
	"testing"
)

func TestService_UpdateStatusByID_InvalidInput(t *testing.T) {
	// Инициализация моков
	mockRepo := new(mocks.Repository)
	svc := service.NewService(mockRepo, zap.NewNop().Sugar())

	// Создание тестового контекста Fiber
	app := fiber.New()
	app.Put("/tasks/:id/:status", svc.UpdateStatusByID)

	// Создание запроса с невалидными параметрами (пустой ID)
	req := httptest.NewRequest("PUT", "/tasks/niceId/awesomeStatus", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	defer resp.Body.Close()

	// Чтение и парсинг ответа
	body, _ := io.ReadAll(resp.Body)
	var response dto.Response
	_ = json.Unmarshal(body, &response)

	// Проверки
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, "error", response.Status)
	assert.Nil(t, response.Data)
	assert.NotNil(t, response.Error)
	assert.Equal(t, dto.FieldIncorrect, response.Error.Code)
	assert.Contains(t, response.Error.Desc, "Invalid int")

	// Репозиторий не должен вызываться
	mockRepo.AssertNotCalled(t, "UpdateStatusByID")
}

func TestService_UpdateStatusByID_NotFound(t *testing.T) {
	// Инициализация моков
	mockRepo := new(mocks.Repository)
	svc := service.NewService(mockRepo, zap.NewNop().Sugar())

	// Создание тестового контекста Fiber
	app := fiber.New()
	app.Put("/tasks/:id/:status", svc.UpdateStatusByID)

	// Настройка мока для возврата ошибки "не найдено"
	mockRepo.On("UpdateStatusByID", mock.Anything, "999", "completed").Return(dto.ErrNotFound)

	// Создание запроса
	req := httptest.NewRequest("PUT", "/tasks/999/completed", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	defer resp.Body.Close()

	// Чтение и парсинг ответа
	body, _ := io.ReadAll(resp.Body)
	var response dto.Response
	err = json.Unmarshal(body, &response)
	assert.NoError(t, err)
	// Проверки
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
	assert.Equal(t, "error", response.Status)
	assert.Nil(t, response.Data)
	assert.NotNil(t, response.Error)
	assert.Equal(t, dto.NotFound, response.Error.Code)
	assert.Equal(t, "not found", response.Error.Desc)

	// Проверка вызова репозитория
	mockRepo.AssertExpectations(t)
}
