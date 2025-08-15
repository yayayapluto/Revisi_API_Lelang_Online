package village

import (
	"encoding/json"
	"errors"
	villageDto "github.com/API_Lelang_Online_Go/internal/modules/village/dto"
	shared "github.com/API_Lelang_Online_Go/internal/shared/helper"
	shared2 "github.com/API_Lelang_Online_Go/internal/shared/responses"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"mime/multipart"
	"path/filepath"
	"strings"
)

type Handler struct {
	useCase UseCase
}

func NewHandler(uc UseCase) *Handler {
	return &Handler{useCase: uc}
}

func (h *Handler) List(c *fiber.Ctx) error {
	page := shared.GetQueryInt(c, "page", 1)
	size := shared.GetQueryInt(c, "size", 10)
	offset := (page - 1) * size

	sortBy := c.Query("sortBy")
	sortDir := c.Query("sortDir")

	var search *string
	searchParam := c.Query("search")

	if searchParam != "" {
		search = &searchParam
	}

	villages, total, err := h.useCase.List(c.UserContext(), size, offset, &sortBy, &sortDir, search)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(shared2.ErrorResponse("Failed to retrieve villages", err.Error()))
	}

	currentPageURL := shared.BuildPageURL(c, page)
	firstPageURL := shared.BuildPageURL(c, 1)

	var nextPageURL, prevPageURL *string

	if offset+len(*villages) < int(total) {
		next := shared.BuildPageURL(c, page+1)
		nextPageURL = &next
	}

	if page > 1 {
		prev := shared.BuildPageURL(c, page-1)
		prevPageURL = &prev
	}

	pagination := shared2.NewPaginationResponse[Village](page, currentPageURL, size, *villages, &firstPageURL, nextPageURL, prevPageURL)
	return c.JSON(shared2.SuccessResponse("Categories retrieved successfully", &pagination))
}

func (h *Handler) Create(c *fiber.Ctx) error {
	var req villageDto.CreateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(shared2.ErrorResponse("Invalid payload", err.Error()))
	}
	village := &Village{
		Code:     req.Code,
		Nama:     req.Nama,
		FullCode: req.FullCode,
	}
	if err := h.useCase.Create(c.UserContext(), village); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(shared2.ErrorResponse("Failed to create village", err.Error()))
	}
	res := &Village{
		ID:        village.ID,
		CreatedAt: village.CreatedAt,
		UpdatedAt: village.UpdatedAt,
		Code:      village.Code,
		Nama:      village.Nama,
		FullCode:  village.FullCode,
	}
	return c.Status(fiber.StatusOK).JSON(shared2.SuccessResponse[Village]("Successfully create village", res))
}

func (h *Handler) Batch(c *fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(shared2.ErrorResponse("invalid file", err.Error()))
	}

	if strings.ToLower(filepath.Ext(file.Filename)) != ".json" {
		return c.Status(fiber.StatusBadRequest).JSON(shared2.ErrorResponse("invalid file type, only .json allowed", nil))
	}

	f, err := file.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(shared2.ErrorResponse("unable to open file", nil))
	}
	defer func(f multipart.File) {
		err := f.Close()
		if err != nil {

		}
	}(f)

	var villages []Village
	if err := json.NewDecoder(f).Decode(&villages); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(shared2.ErrorResponse("invalid json format", err.Error()))
	}

	if err := h.useCase.Batch(c.UserContext(), &villages); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(shared2.ErrorResponse("failed to save batch villages", err.Error()))
	}

	return c.Status(fiber.StatusOK).JSON(shared2.SuccessResponse[[]Village]("successfully uploaded batch villages", &villages))
}

func (h *Handler) GetByID(c *fiber.Ctx) error {
	id, err := shared.ParseIdParam(c, "id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(shared2.ErrorResponse("invalid id", err.Error()))
	}

	village, err := h.useCase.GetByID(c.UserContext(), id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(shared2.ErrorResponse("Cannot find village with provided id", err.Error()))
		}
		return c.Status(fiber.StatusInternalServerError).JSON(shared2.ErrorResponse("Failed to retrieve village", err.Error()))
	}

	return c.Status(fiber.StatusOK).JSON(shared2.SuccessResponse[Village]("Successfully retrieve village", village))
}

func (h *Handler) Update(c *fiber.Ctx) error {
	id, err := shared.ParseIdParam(c, "id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(shared2.ErrorResponse("invalid id", err.Error()))
	}

	village, err := h.useCase.GetByID(c.UserContext(), id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(shared2.ErrorResponse("Cannot find village with provided id", err.Error()))
		}
		return c.Status(fiber.StatusInternalServerError).JSON(shared2.ErrorResponse("Failed to retrieve village", err.Error()))
	}

	var req villageDto.UpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(shared2.ErrorResponse("Invalid payload", err.Error()))
	}

	if req.Code != nil {
		village.Code = *req.Code
	}

	if req.Nama != nil {
		village.Nama = *req.Nama
	}

	if req.FullCode != nil {
		village.FullCode = *req.FullCode
	}

	if err := h.useCase.Update(c.UserContext(), village); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(shared2.ErrorResponse("Cannot find village with provided id", nil))
		}
		return c.Status(fiber.StatusInternalServerError).JSON(shared2.ErrorResponse("Failed to update village", err.Error()))
	}

	return c.Status(fiber.StatusOK).JSON(shared2.SuccessResponse[Village]("Successfully update village", village))
}

func (h *Handler) Delete(c *fiber.Ctx) error {
	id, err := shared.ParseIdParam(c, "id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(shared2.ErrorResponse("invalid id", err.Error()))
	}

	if err := h.useCase.Delete(c.UserContext(), id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(shared2.ErrorResponse("Cannot find village with provided id", nil))
		}
		return c.Status(fiber.StatusInternalServerError).JSON(shared2.ErrorResponse("Failed to delete village", err.Error()))
	}

	return c.Status(fiber.StatusOK).JSON(shared2.SuccessResponse[any]("Successfully delete village", nil))
}
