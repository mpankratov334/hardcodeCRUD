package service_test

import (
	"encoding/json"
	"github.com/jackc/pgx/v5"
	"io"
	"net/http/httptest"
	"testing"

	"TemplatestPGSQL/internal/dto"
	"TemplatestPGSQL/internal/repo"
	"TemplatestPGSQL/internal/repo/mocks"
	"TemplatestPGSQL/internal/service"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

func TestService_GetTaskByID_Success(t *testing.T) {
	// Инициализация моков
	mockRepo := new(mocks.Repository)
	svc := service.NewService(mockRepo, zap.NewNop().Sugar())

	// Создание тестового контекста Fiber
	app := fiber.New()
	app.Get("/tasks/:id", svc.GetTaskByID)

	// Ожидаемый результат репозитория
	expectedTask := &repo.Task{
		DataObject: repo.DataObject{
			ID:    "123",
			Title: "Test Task",
			Data:  "Task content",
		},
		UserID: "1",
	}
	mockRepo.On("GetTaskByID", mock.Anything, "123").Return(expectedTask, nil)

	// Установка параметров пути
	req := httptest.NewRequest("GET", "/tasks/123", nil)
	// Вызов тестируемого метода
	// получение реального ответа
	realRow, err := app.Test(req)
	body, _ := io.ReadAll(realRow.Body)
	var real dto.Response

	// парсинг ответа
	_ = json.Unmarshal(body, &real)
	rData, _ := real.Data.(map[string]interface{})
	rDesc, _ := rData["data"].(string)
	var expErr *dto.Error

	// Проверки
	assert.Equal(t, expectedTask.Data, rDesc)
	assert.Equal(t, expErr, real.Error)
	assert.Equal(t, "success", real.Status)
	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
}

func TestService_GetTaskByID_InvalidID(t *testing.T) {
	mockRepo := new(mocks.Repository)
	svc := service.NewService(mockRepo, zap.NewNop().Sugar())

	app := fiber.New()
	app.Get("/tasks/:id", svc.GetTaskByID)

	// Пустой ID в пути вызовет ошибку валидации
	req := httptest.NewRequest("GET", "/tasks/1,2", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var response dto.Response
	_ = json.Unmarshal(body, &response)

	// Проверяем структуру ошибки
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, "error", response.Status)
	assert.Nil(t, response.Data)
	assert.NotNil(t, response.Error)
	assert.Equal(t, dto.FieldIncorrect, response.Error.Code)

	// Репозиторий не должен вызываться
	mockRepo.AssertNotCalled(t, "GetTaskByID")
}

func TestService_GetTaskByID_NotFound(t *testing.T) {
	mockRepo := new(mocks.Repository)
	svc := service.NewService(mockRepo, zap.NewNop().Sugar())

	app := fiber.New()
	app.Get("/tasks/:id", svc.GetTaskByID)

	// Настраиваем мок для возврата ошибки
	mockRepo.On("GetTaskByID", mock.Anything, "999").Return(nil, pgx.ErrNoRows)

	req := httptest.NewRequest("GET", "/tasks/999", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var response dto.Response
	_ = json.Unmarshal(body, &response)

	// Проверяем
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
	assert.Equal(t, "error", response.Status)
	assert.Nil(t, response.Data)
	assert.NotNil(t, response.Error)
	assert.Equal(t, dto.ServiceUnavailable, response.Error.Code)
	assert.Equal(t, dto.InternalError, response.Error.Desc)

	mockRepo.AssertExpectations(t)
}
