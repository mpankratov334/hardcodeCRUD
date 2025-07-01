package service_test

import (
	"TemplatestPGSQL/pkg/validator"
	"encoding/json"
	"errors"
	"github.com/valyala/fasthttp"
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

func TestService_CreateTask_Success(t *testing.T) {
	// Инициализация моков
	mockRepo := new(mocks.Repository)
	logger := zap.NewNop().Sugar()
	svc := service.NewService(mockRepo, logger)

	// Создание тестового контекста Fiber
	app := fiber.New()
	ctx := app.AcquireCtx(&fasthttp.RequestCtx{})
	defer app.ReleaseCtx(ctx)

	// Установка тела запроса
	body := []byte(`{"title": "Test Task", "data": "Task content", "user_id": "1"}`)
	ctx.Request().SetBody(body)
	ctx.Request().Header.SetContentType("application/json")
	ctx.Request().Header.SetMethod("POST")

	// Ожидаемый вызов репозитория
	expectedTask := repo.Task{
		DataObject: repo.DataObject{
			Title: "Test Task",
			Data:  "Task content",
		},
		UserID: "1",
	}
	mockRepo.On("CreateTask", mock.Anything, expectedTask).Return(nil)

	// Вызов тестируемого метода
	err := svc.CreateTask(ctx)

	// Проверки
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, ctx.Response().StatusCode())

	var response dto.Response
	_ = json.Unmarshal(ctx.Response().Body(), &response)
	assert.Equal(t, "success", response.Status)
	assert.Equal(t, nil, response.Data)

	mockRepo.AssertExpectations(t)
}

func TestService_CreateTask_InvalidJSON(t *testing.T) {
	mockRepo := new(mocks.Repository)
	logger := zap.NewNop().Sugar()
	svc := service.NewService(mockRepo, logger)

	// Создание тестового контекста
	app := fiber.New()
	ctx := app.AcquireCtx(&fasthttp.RequestCtx{})
	defer app.ReleaseCtx(ctx)

	// Некорректное JSON тело
	body := []byte(`{"title": notString, "data": "Task content", "user_id": "1"}`)
	ctx.Request().SetBody(body)
	ctx.Request().Header.SetContentType("application/json")
	ctx.Request().Header.SetMethod("POST")

	// Вызов метода
	err := svc.CreateTask(ctx)

	// Проверки
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, ctx.Response().StatusCode())

	var errorResp dto.Response
	_ = json.Unmarshal(ctx.Response().Body(), &errorResp)
	assert.Equal(t, "error", errorResp.Status)

	mockRepo.AssertNotCalled(t, "CreateTask")
}

func TestService_CreateTask_ValidationError(t *testing.T) {
	mockRepo := new(mocks.Repository)
	logger := zap.NewNop().Sugar()
	svc := service.NewService(mockRepo, logger)

	// Создание тестового контекста
	app := fiber.New()
	ctx := app.AcquireCtx(&fasthttp.RequestCtx{})
	defer app.ReleaseCtx(ctx)

	// Тело с ошибкой валидации (отсутствует title)
	body := []byte(`{"data": "Task content", "user_id": "1"}`)
	ctx.Request().SetBody(body)
	ctx.Request().Header.SetContentType("application/json")
	ctx.Request().Header.SetMethod("POST")

	// Вызов метода
	err := svc.CreateTask(ctx)

	// Проверки
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, ctx.Response().StatusCode())

	var errorResp dto.Response
	_ = json.Unmarshal(ctx.Response().Body(), &errorResp)
	assert.Equal(t, "error", errorResp.Status)
	assert.Equal(t, dto.FieldIncorrect, errorResp.Error.Code)
	assert.Equal(t, validator.ErrFieldRequired+": PostRequest.Title", errorResp.Error.Desc)

	mockRepo.AssertNotCalled(t, "CreateTask")
}

func TestService_CreateTask_RepositoryError(t *testing.T) {
	mockRepo := new(mocks.Repository)
	logger := zap.NewNop().Sugar()
	svc := service.NewService(mockRepo, logger)

	// Создание тестового контекста
	app := fiber.New()
	ctx := app.AcquireCtx(&fasthttp.RequestCtx{})
	defer app.ReleaseCtx(ctx)

	// Корректное тело запроса
	body := []byte(`{"title": "Test Task", "data": "Task content", "user_id": "1"}`)
	ctx.Request().SetBody(body)
	ctx.Request().Header.SetContentType("application/json")
	ctx.Request().Header.SetMethod("POST")

	// Настройка мока на возврат ошибки
	mockRepo.On("CreateTask", mock.Anything, mock.Anything).Return(errors.New("database error"))

	// Вызов метода
	err := svc.CreateTask(ctx)

	// Проверки
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, ctx.Response().StatusCode())

	var errorResp dto.Response
	_ = json.Unmarshal(ctx.Response().Body(), &errorResp)
	assert.Equal(t, "error", errorResp.Status)
	assert.Equal(t, dto.ServiceUnavailable, errorResp.Error.Code)
	assert.Equal(t, dto.InternalError, errorResp.Error.Desc)

	mockRepo.AssertExpectations(t)
}
