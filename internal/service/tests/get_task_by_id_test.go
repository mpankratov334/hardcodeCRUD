package service_test

import (
	"encoding/json"
	"fmt"
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
	fmt.Println("\n", "here")
	realRow, err := app.Test(req)
	body, _ := io.ReadAll(realRow.Body)
	var real dto.Response

	// ТИП ANY ДАМЫ И ГОСПОДА
	// парсинг ответа
	_ = json.Unmarshal(body, &real)
	rData, _ := real.Data.(map[string]interface{})
	rDesc, _ := rData["data"].(string)

	// Проверки
	fmt.Println("\n", rData)
	fmt.Println()
	assert.Equal(t, expectedTask.Data, rDesc)
	assert.Equal(t, nil, real.Error)
	assert.Equal(t, "success", real.Status)
	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
}

//
//func TestService_GetTaskByID_ValidationError(t *testing.T) {
//	mockRepo := new(mocks.Repository)
//	logger := zap.NewNop().Sugar()
//	svc := service.NewService(mockRepo, logger)
//
//	// Создание тестового контекста
//	app := fiber.New()
//	ctx := app.AcquireCtx(&fasthttp.RequestCtx{})
//	defer app.ReleaseCtx(ctx)
//
//	// Пустой ID (вызовет ошибку валидации)
//	ctx.Params("id", "")
//
//	// Вызов метода
//	err := svc.GetTaskByID(ctx)
//
//	// Проверки
//	assert.NoError(t, err)
//	assert.Equal(t, fiber.StatusBadRequest, ctx.Response().StatusCode())
//
//	var errorResp dto.Response
//	_ = json.Unmarshal(ctx.Response().Body(), &errorResp)
//	assert.Equal(t, "error", errorResp.Status)
//	assert.Equal(t, dto.FieldIncorrect, errorResp.Error.Code)
//	assert.Contains(t, errorResp.Error.Desc, "RequestWithId.ID")
//
//	mockRepo.AssertNotCalled(t, "GetTaskByID")
//}
//
//func TestService_GetTaskByID_RepositoryError(t *testing.T) {
//	mockRepo := new(mocks.Repository)
//	logger := zap.NewNop().Sugar()
//	svc := service.NewService(mockRepo, logger)
//
//	// Создание тестового контекста
//	app := fiber.New()
//	ctx := app.AcquireCtx(&fasthttp.RequestCtx{})
//	defer app.ReleaseCtx(ctx)
//
//	// Корректный ID
//	ctx.Params("id", "123")
//
//	// Настройка мока на возврат ошибки
//	mockRepo.On("GetTaskByID", mock.Anything, "123").Return(nil, errors.New("database error"))
//
//	// Вызов метода
//	err := svc.GetTaskByID(ctx)
//
//	// Проверки
//	assert.NoError(t, err)
//	assert.Equal(t, fiber.StatusInternalServerError, ctx.Response().StatusCode())
//
//	var errorResp dto.Response
//	_ = json.Unmarshal(ctx.Response().Body(), &errorResp)
//	assert.Equal(t, "error", errorResp.Status)
//	assert.Equal(t, dto.ServiceUnavailable, errorResp.Error.Code)
//	assert.Equal(t, dto.InternalError, errorResp.Error.Desc)
//
//	mockRepo.AssertExpectations(t)
//}
//
//func TestService_GetTaskByID_NotFound(t *testing.T) {
//	mockRepo := new(mocks.Repository)
//	logger := zap.NewNop().Sugar()
//	svc := service.NewService(mockRepo, logger)
//
//	// Создание тестового контекста
//	app := fiber.New()
//	ctx := app.AcquireCtx(&fasthttp.RequestCtx{})
//	defer app.ReleaseCtx(ctx)
//
//	// Корректный ID
//	ctx.Params("id", "999")
//
//	// Настройка мока на возврат nil (объект не найден)
//	mockRepo.On("GetTaskByID", mock.Anything, "999").Return(nil, nil)
//
//	// Вызов метода
//	err := svc.GetTaskByID(ctx)
//
//	// Проверки
//	assert.NoError(t, err)
//	assert.Equal(t, fiber.StatusOK, ctx.Response().StatusCode())
//
//	var response dto.Response
//	_ = json.Unmarshal(ctx.Response().Body(), &response)
//	assert.Equal(t, "success", response.Status)
//	assert.Nil(t, response.Data)
//
//	mockRepo.AssertExpectations(t)
//}
