package subdistrict

import (
	"encoding/json"
	"errors"
	subdistrictDto "github.com/API_Lelang_Online_Go/internal/modules/subdistrict/dto"
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

	subdistricts, total, err := h.useCase.List(c.UserContext(), size, offset, &sortBy, &sortDir, search)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(shared2.ErrorResponse("Failed to retrieve subdistricts", err.Error()))
	}

	currentPageURL := shared.BuildPageURL(c, page)
	firstPageURL := shared.BuildPageURL(c, 1)

	var nextPageURL, prevPageURL *string

	if offset+len(*subdistricts) < int(total) {
		next := shared.BuildPageURL(c, page+1)
		nextPageURL = &next
	}

	if page > 1 {
		prev := shared.BuildPageURL(c, page-1)
		prevPageURL = &prev
	}

	pagination := shared2.NewPaginationResponse[Subdistrict](page, currentPageURL, size, *subdistricts, &firstPageURL, nextPageURL, prevPageURL)
	return c.JSON(shared2.SuccessResponse("Categories retrieved successfully", &pagination))
}

func (h *Handler) Create(c *fiber.Ctx) error {
	var req subdistrictDto.CreateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(shared2.ErrorResponse("Invalid payload", err.Error()))
	}
	subdistrict := &Subdistrict{
		Code:     req.Code,
		Nama:     req.Nama,
		FullCode: req.FullCode,
	}
	if err := h.useCase.Create(c.UserContext(), subdistrict); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(shared2.ErrorResponse("Failed to create subdistrict", err.Error()))
	}
	res := &Subdistrict{
		ID:        subdistrict.ID,
		CreatedAt: subdistrict.CreatedAt,
		UpdatedAt: subdistrict.UpdatedAt,
		Code:      subdistrict.Code,
		Nama:      subdistrict.Nama,
		FullCode:  subdistrict.FullCode,
	}
	return c.Status(fiber.StatusOK).JSON(shared2.SuccessResponse[Subdistrict]("Successfully create subdistrict", res))
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

	var subdistricts []Subdistrict
	if err := json.NewDecoder(f).Decode(&subdistricts); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(shared2.ErrorResponse("invalid json format", nil))
	}

	if err := h.useCase.Batch(c.UserContext(), &subdistricts); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(shared2.ErrorResponse("failed to save batch subdistricts", err.Error()))
	}

	return c.Status(fiber.StatusOK).JSON(shared2.SuccessResponse[[]Subdistrict]("successfully uploaded batch subdistricts", &subdistricts))
}

func (h *Handler) GetByID(c *fiber.Ctx) error {
	id, err := shared.ParseIdParam(c, "id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(shared2.ErrorResponse("invalid id", err.Error()))
	}

	subdistrict, err := h.useCase.GetByID(c.UserContext(), id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(shared2.ErrorResponse("Cannot find subdistrict with provided id", err.Error()))
		}
		return c.Status(fiber.StatusInternalServerError).JSON(shared2.ErrorResponse("Failed to retrieve subdistrict", err.Error()))
	}

	return c.Status(fiber.StatusOK).JSON(shared2.SuccessResponse[Subdistrict]("Successfully retrieve subdistrict", subdistrict))
}

func (h *Handler) Update(c *fiber.Ctx) error {
	id, err := shared.ParseIdParam(c, "id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(shared2.ErrorResponse("invalid id", err.Error()))
	}

	subdistrict, err := h.useCase.GetByID(c.UserContext(), id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(shared2.ErrorResponse("Cannot find subdistrict with provided id", err.Error()))
		}
		return c.Status(fiber.StatusInternalServerError).JSON(shared2.ErrorResponse("Failed to retrieve subdistrict", err.Error()))
	}

	var req subdistrictDto.UpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(shared2.ErrorResponse("Invalid payload", err.Error()))
	}

	if req.Code != nil {
		subdistrict.Code = *req.Code
	}

	if req.Nama != nil {
		subdistrict.Nama = *req.Nama
	}

	if req.FullCode != nil {
		subdistrict.FullCode = *req.FullCode
	}

	if err := h.useCase.Update(c.UserContext(), subdistrict); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(shared2.ErrorResponse("Cannot find subdistrict with provided id", nil))
		}
		return c.Status(fiber.StatusInternalServerError).JSON(shared2.ErrorResponse("Failed to update subdistrict", err.Error()))
	}

	return c.Status(fiber.StatusOK).JSON(shared2.SuccessResponse[Subdistrict]("Successfully update subdistrict", subdistrict))
}

func (h *Handler) Delete(c *fiber.Ctx) error {
	id, err := shared.ParseIdParam(c, "id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(shared2.ErrorResponse("invalid id", err.Error()))
	}

	if err := h.useCase.Delete(c.UserContext(), id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(shared2.ErrorResponse("Cannot find subdistrict with provided id", nil))
		}
		return c.Status(fiber.StatusInternalServerError).JSON(shared2.ErrorResponse("Failed to delete subdistrict", err.Error()))
	}

	return c.Status(fiber.StatusOK).JSON(shared2.SuccessResponse[any]("Successfully delete subdistrict", nil))
}
