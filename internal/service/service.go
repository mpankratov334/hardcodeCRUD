package service

import (
	"TemplatestPGSQL/internal/dto"
	repo2 "TemplatestPGSQL/internal/repo"
	"TemplatestPGSQL/pkg/validator"
	"encoding/json"
	"errors"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type Service interface {
	CreateUser(ctx *fiber.Ctx) error
	CreateTask(ctx *fiber.Ctx) error
	GetAllTasks(ctx *fiber.Ctx) error
	GetTaskByID(ctx *fiber.Ctx) error
	GetLastTaskByUserID(ctx *fiber.Ctx) error
	GetAllTasksByUserID(ctx *fiber.Ctx) error
	GetTasksByUserName(ctx *fiber.Ctx) error
	UpdateStatusByID(ctx *fiber.Ctx) error
	DeleteTaskByID(ctx *fiber.Ctx) error
}

type service struct {
	repo repo2.Repository
	log  *zap.SugaredLogger
}

func NewService(repo repo2.Repository, logger *zap.SugaredLogger) Service {
	return &service{
		repo: repo,
		log:  logger,
	}
}

func (s *service) CreateTask(ctx *fiber.Ctx) error {
	var obj PostRequest

	// deserialize  JSON-request
	if err := json.Unmarshal(ctx.Body(), &obj); err != nil {
		s.log.Error("Invalid request body", zap.Error(err))
		return dto.BadResponseError(ctx, dto.FieldBadFormat, "Invalid request body")
	}

	// validation
	if vErr := validator.Validate(ctx.Context(), obj); vErr != nil {
		s.log.Error("Invalid request data", zap.Error(vErr))
		return dto.BadResponseError(ctx, dto.FieldIncorrect, vErr.Error())
	}

	// adds to memory
	dataObj := repo2.Task{
		DataObject: repo2.DataObject{
			Title: obj.Title,
			Data:  obj.Data,
		},
		UserID: obj.UserID,
	}
	err := s.repo.CreateTask(ctx.Context(), dataObj)
	if err != nil {
		s.log.Error("Failed to insert object", zap.Error(err))
		return dto.InternalServerError(ctx)
	}
	s.log.Infof("object was appended %s", dataObj.Title)

	// forms the answer
	response := dto.Response{
		Status: "success",
		Data:   map[string]int{"task_id": 1},
	}

	return ctx.Status(fiber.StatusOK).JSON(response)
}

func (s *service) GetTasksByUserName(ctx *fiber.Ctx) error {
	req := RequestWithUserName{Name: ctx.Params("username")}

	// Validation
	if vErr := validator.Validate(ctx.Context(), req); vErr != nil {
		s.log.Error("Invalid request", zap.Error(vErr))
		return dto.BadResponseError(ctx, dto.FieldIncorrect, vErr.Error())
	}

	// Gets from memory
	objPtr, err := s.repo.GetTasksByUserName(ctx.Context(), req.Name)
	if err != nil {
		s.log.Error("Failed to get task", zap.Error(err))
		if errors.Is(err, dto.ErrNotFound) {
			return dto.NotFoundError(ctx, dto.NotFound, err.Error())
		}
		return dto.InternalServerError(ctx)
	}

	// Forms answer
	response := dto.Response{
		Status: "success",
		Data:   objPtr,
	}
	return ctx.Status(fiber.StatusOK).JSON(response)
}

func (s *service) GetAllTasks(ctx *fiber.Ctx) error {
	// Gets from memory
	obj, err := s.repo.GetAllTasks(ctx.Context())
	if err != nil {
		s.log.Error("Failed to get task", zap.Error(err))
		return dto.InternalServerError(ctx)
	}

	// Forms answer
	jsonData, err := json.Marshal(obj)
	if err != nil {
		s.log.Error("Failed to marshal response", zap.Error(err))
		return dto.InternalServerError(ctx)
	}
	response := dto.Response{
		Status: "success",
		Data:   jsonData,
	}
	s.log.Info("all tasks was read and sent")
	return ctx.Status(fiber.StatusOK).JSON(response)
}

func (s *service) GetTaskByID(ctx *fiber.Ctx) error {
	req := RequestWithId{ID: ctx.Params("id")}

	// Validation
	if vErr := validator.Validate(ctx.Context(), req); vErr != nil {
		s.log.Error("Invalid request id", zap.Error(vErr))
		return dto.BadResponseError(ctx, dto.FieldIncorrect, vErr.Error())
	}

	// Gets from memory
	objPtr, err := s.repo.GetTaskByID(ctx.Context(), req.ID)
	if err != nil {
		s.log.Error("Failed to get task", zap.Error(err))
		return dto.InternalServerError(ctx)
	}

	// Forms answer
	response := dto.Response{
		Status: "success",
		Data:   objPtr,
	}
	return ctx.Status(fiber.StatusOK).JSON(response)
}

func (s *service) GetLastTaskByUserID(ctx *fiber.Ctx) error {
	req := RequestWithId{ID: ctx.Params("id")}

	// Validation
	if vErr := validator.Validate(ctx.Context(), req); vErr != nil {
		s.log.Error("Invalid request id", zap.Error(vErr))
		return dto.BadResponseError(ctx, dto.FieldIncorrect, vErr.Error())
	}

	// Gets from memory
	objPtr, err := s.repo.GetLastTaskByUserID(ctx.Context(), req.ID)
	if err != nil {
		s.log.Error("Failed to get task", zap.Error(err))
		if errors.Is(err, dto.ErrNotFound) {
			return dto.NotFoundError(ctx, dto.NotFound, err.Error())
		}
		return dto.InternalServerError(ctx)
	}

	// Forms answer
	response := dto.Response{
		Status: "success",
		Data:   objPtr,
	}
	return ctx.Status(fiber.StatusOK).JSON(response)
}

func (s *service) GetAllTasksByUserID(ctx *fiber.Ctx) error {
	req := RequestWithId{ID: ctx.Params("id")}

	// Validation
	if vErr := validator.Validate(ctx.Context(), req); vErr != nil {
		s.log.Error("Invalid request id", zap.Error(vErr))
		return dto.BadResponseError(ctx, dto.FieldIncorrect, vErr.Error())
	}

	obj, err := s.repo.GetAllTasksByUserID(ctx.Context(), req.ID)
	if err != nil {
		s.log.Error("Failed to get task", zap.Error(err))
		if errors.Is(err, dto.ErrNotFound) {
			return dto.NotFoundError(ctx, dto.NotFound, err.Error())
		}
		return dto.InternalServerError(ctx)
	}

	jsonData, err := json.Marshal(obj)
	if err != nil {
		s.log.Error("Failed to marshal response", zap.Error(err))
		return dto.InternalServerError(ctx)
	}
	response := dto.Response{
		Status: "success",
		Data:   jsonData,
	}
	s.log.Info("whole memory was read and sent")
	return ctx.Status(fiber.StatusOK).JSON(response)
}

func (s *service) UpdateStatusByID(ctx *fiber.Ctx) error {
	req := UpdateRequest{ID: ctx.Params("id"), Status: ctx.Params("status")}

	// Validation
	if vErr := validator.Validate(ctx.Context(), req); vErr != nil {
		s.log.Error("Invalid request id", zap.Error(vErr))
		return dto.BadResponseError(ctx, dto.FieldIncorrect, vErr.Error())
	}

	// Gets from memory
	err := s.repo.UpdateStatusByID(ctx.Context(), req.ID, req.Status)
	if err != nil {
		s.log.Error("Failed to get task", zap.Error(err))
		if errors.Is(err, dto.ErrNotFound) {
			return dto.NotFoundError(ctx, dto.NotFound, err.Error())
		}
		return dto.InternalServerError(ctx)
	}

	// Forms answer
	response := dto.Response{
		Status: "success",
	}
	return ctx.Status(fiber.StatusOK).JSON(response)
}

func (s *service) DeleteTaskByID(ctx *fiber.Ctx) error {
	req := UpdateRequest{ID: ctx.Params("id"), Status: ctx.Params("status")}

	// Validation
	if vErr := validator.Validate(ctx.Context(), req); vErr != nil {
		s.log.Error("Invalid request id", zap.Error(vErr))
		return dto.BadResponseError(ctx, dto.FieldIncorrect, vErr.Error())
	}

	// Gets from memory
	err := s.repo.DeleteTaskByID(ctx.Context(), req.ID)
	if err != nil {
		s.log.Error("Failed to get task", zap.Error(err))
		if errors.Is(err, dto.ErrNotFound) {
			return dto.NotFoundError(ctx, dto.NotFound, err.Error())
		}
		return dto.InternalServerError(ctx)
	}

	// Forms answer
	response := dto.Response{
		Status: "success",
	}
	return ctx.Status(fiber.StatusOK).JSON(response)
}

func (s *service) CreateUser(ctx *fiber.Ctx) error {
	var obj PostUserRequest

	// deserialize  JSON-request
	if err := json.Unmarshal(ctx.Body(), &obj); err != nil {
		s.log.Error("Invalid request body", zap.Error(err))
		return dto.BadResponseError(ctx, dto.FieldBadFormat, "Invalid request body")
	}

	// validation
	if vErr := validator.Validate(ctx.Context(), obj); vErr != nil {
		s.log.Error("Invalid request data", zap.Error(vErr))
		return dto.BadResponseError(ctx, dto.FieldIncorrect, vErr.Error())
	}

	// adds to memory
	user := repo2.User{
		Name:     obj.Name,
		Password: obj.Password,
	}
	err := s.repo.CreateUser(ctx.Context(), user)
	if err != nil {
		s.log.Error("Failed to insert object", zap.Error(err))
		return dto.InternalServerError(ctx)
	}
	s.log.Infof("object was appended %s", user.Name)

	// forms the answer
	response := dto.Response{
		Status: "success",
		Data:   map[string]int{"task_id": 1},
	}

	return ctx.Status(fiber.StatusOK).JSON(response)
}
