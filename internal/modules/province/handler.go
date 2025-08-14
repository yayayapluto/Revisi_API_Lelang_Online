package province

import (
	"encoding/json"
	"errors"
	"github.com/API_Lelang_Online_Go/internal/modules/province/dto"
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

	provinces, total, err := h.useCase.List(c.UserContext(), size, offset, &sortBy, &sortDir, search)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(shared2.ErrorResponse("Failed to retrieve provinces", err.Error()))
	}

	currentPageURL := shared.BuildPageURL(c, page)
	firstPageURL := shared.BuildPageURL(c, 1)

	var nextPageURL, prevPageURL *string

	if offset+len(*provinces) < int(total) {
		next := shared.BuildPageURL(c, page+1)
		nextPageURL = &next
	}

	if page > 1 {
		prev := shared.BuildPageURL(c, page-1)
		prevPageURL = &prev
	}

	pagination := shared2.NewPaginationResponse[Province](page, currentPageURL, size, *provinces, &firstPageURL, nextPageURL, prevPageURL)
	return c.JSON(shared2.SuccessResponse("Categories retrieved successfully", &pagination))
}

func (h *Handler) Create(c *fiber.Ctx) error {
	var req provinceDto.CreateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(shared2.ErrorResponse("Invalid payload", err.Error()))
	}
	province := &Province{
		Code:     req.Code,
		Nama:     req.Nama,
		FullCode: req.FullCode,
	}
	if err := h.useCase.Create(c.UserContext(), province); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(shared2.ErrorResponse("Failed to create province", err.Error()))
	}
	res := &Province{
		ID:        province.ID,
		CreatedAt: province.CreatedAt,
		UpdatedAt: province.UpdatedAt,
		Code:      province.Code,
		Nama:      province.Nama,
		FullCode:  province.FullCode,
	}
	return c.Status(fiber.StatusOK).JSON(shared2.SuccessResponse[Province]("Successfully create province", res))
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

	var provinces []Province
	if err := json.NewDecoder(f).Decode(&provinces); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(shared2.ErrorResponse("invalid json format", nil))
	}

	if err := h.useCase.Batch(c.UserContext(), &provinces); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(shared2.ErrorResponse("failed to save batch provinces", err.Error()))
	}

	return c.Status(fiber.StatusOK).JSON(shared2.SuccessResponse[[]Province]("successfully uploaded batch provinces", &provinces))
}

func (h *Handler) GetByID(c *fiber.Ctx) error {
	id, err := shared.ParseIdParam(c, "id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(shared2.ErrorResponse("invalid id", err.Error()))
	}

	province, err := h.useCase.GetByID(c.UserContext(), id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(shared2.ErrorResponse("Cannot find province with provided id", err.Error()))
		}
		return c.Status(fiber.StatusInternalServerError).JSON(shared2.ErrorResponse("Failed to retrieve province", err.Error()))
	}

	return c.Status(fiber.StatusOK).JSON(shared2.SuccessResponse[Province]("Successfully retrieve province", province))
}

func (h *Handler) Update(c *fiber.Ctx) error {
	id, err := shared.ParseIdParam(c, "id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(shared2.ErrorResponse("invalid id", err.Error()))
	}

	province, err := h.useCase.GetByID(c.UserContext(), id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(shared2.ErrorResponse("Cannot find province with provided id", err.Error()))
		}
		return c.Status(fiber.StatusInternalServerError).JSON(shared2.ErrorResponse("Failed to retrieve province", err.Error()))
	}

	var req provinceDto.UpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(shared2.ErrorResponse("Invalid payload", err.Error()))
	}

	if req.Code != nil {
		province.Code = *req.Code
	}

	if req.Nama != nil {
		province.Nama = *req.Nama
	}

	if req.FullCode != nil {
		province.FullCode = *req.FullCode
	}

	if err := h.useCase.Update(c.UserContext(), province); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(shared2.ErrorResponse("Cannot find province with provided id", nil))
		}
		return c.Status(fiber.StatusInternalServerError).JSON(shared2.ErrorResponse("Failed to update province", err.Error()))
	}

	return c.Status(fiber.StatusOK).JSON(shared2.SuccessResponse[Province]("Successfully update province", province))
}

func (h *Handler) Delete(c *fiber.Ctx) error {
	id, err := shared.ParseIdParam(c, "id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(shared2.ErrorResponse("invalid id", err.Error()))
	}

	if err := h.useCase.Delete(c.UserContext(), id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(shared2.ErrorResponse("Cannot find province with provided id", nil))
		}
		return c.Status(fiber.StatusInternalServerError).JSON(shared2.ErrorResponse("Failed to delete province", err.Error()))
	}

	return c.Status(fiber.StatusOK).JSON(shared2.SuccessResponse[any]("Successfully delete province", nil))
}
