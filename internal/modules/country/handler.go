package country

import (
	"encoding/json"
	"errors"
	countryDto "github.com/API_Lelang_Online_Go/internal/modules/country/dto"
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

	countries, total, err := h.useCase.List(c.UserContext(), size, offset, &sortBy, &sortDir, search)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(shared2.ErrorResponse("Failed to retrieve countries", err.Error()))
	}

	currentPageURL := shared.BuildPageURL(c, page)
	firstPageURL := shared.BuildPageURL(c, 1)

	var nextPageURL, prevPageURL *string

	if offset+len(*countries) < int(total) {
		next := shared.BuildPageURL(c, page+1)
		nextPageURL = &next
	}

	if page > 1 {
		prev := shared.BuildPageURL(c, page-1)
		prevPageURL = &prev
	}

	pagination := shared2.NewPaginationResponse[Country](page, currentPageURL, size, *countries, &firstPageURL, nextPageURL, prevPageURL)
	return c.JSON(shared2.SuccessResponse("Categories retrieved successfully", &pagination))
}

func (h *Handler) Create(c *fiber.Ctx) error {
	var req countryDto.CreateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(shared2.ErrorResponse("Invalid payload", err.Error()))
	}
	country := &Country{
		Kode:  req.Kode,
		Nama:  req.Nama,
		Nomor: req.Nomor,
	}
	if err := h.useCase.Create(c.UserContext(), country); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(shared2.ErrorResponse("Failed to create country", err.Error()))
	}
	res := &Country{
		ID:        country.ID,
		CreatedAt: country.CreatedAt,
		UpdatedAt: country.UpdatedAt,
		Kode:      country.Kode,
		Nama:      country.Nama,
		Nomor:     country.Nomor,
	}
	return c.Status(fiber.StatusOK).JSON(shared2.SuccessResponse[Country]("Successfully create country", res))
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

	var countries []Country
	if err := json.NewDecoder(f).Decode(&countries); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(shared2.ErrorResponse("invalid json format", nil))
	}

	if err := h.useCase.Batch(c.UserContext(), &countries); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(shared2.ErrorResponse("failed to save batch countries", err.Error()))
	}

	return c.Status(fiber.StatusOK).JSON(shared2.SuccessResponse[[]Country]("successfully uploaded batch countries", &countries))
}

func (h *Handler) GetByID(c *fiber.Ctx) error {
	id, err := shared.ParseIdParam(c, "id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(shared2.ErrorResponse("invalid id", err.Error()))
	}

	country, err := h.useCase.GetByID(c.UserContext(), id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(shared2.ErrorResponse("Cannot find country with provided id", err.Error()))
		}
		return c.Status(fiber.StatusInternalServerError).JSON(shared2.ErrorResponse("Failed to retrieve country", err.Error()))
	}

	return c.Status(fiber.StatusOK).JSON(shared2.SuccessResponse[Country]("Successfully retrieve country", country))
}

func (h *Handler) Update(c *fiber.Ctx) error {
	id, err := shared.ParseIdParam(c, "id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(shared2.ErrorResponse("invalid id", err.Error()))
	}

	country, err := h.useCase.GetByID(c.UserContext(), id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(shared2.ErrorResponse("Cannot find country with provided id", err.Error()))
		}
		return c.Status(fiber.StatusInternalServerError).JSON(shared2.ErrorResponse("Failed to retrieve country", err.Error()))
	}

	var req countryDto.UpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(shared2.ErrorResponse("Invalid payload", err.Error()))
	}

	if req.Kode != nil {
		country.Kode = *req.Kode
	}

	if req.Nama != nil {
		country.Nama = *req.Nama
	}

	if req.Nomor != nil {
		country.Nomor = *req.Nomor
	}

	if err := h.useCase.Update(c.UserContext(), country); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(shared2.ErrorResponse("Cannot find country with provided id", nil))
		}
		return c.Status(fiber.StatusInternalServerError).JSON(shared2.ErrorResponse("Failed to update country", err.Error()))
	}

	return c.Status(fiber.StatusOK).JSON(shared2.SuccessResponse[Country]("Successfully update country", country))
}

func (h *Handler) Delete(c *fiber.Ctx) error {
	id, err := shared.ParseIdParam(c, "id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(shared2.ErrorResponse("invalid id", err.Error()))
	}

	if err := h.useCase.Delete(c.UserContext(), id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(shared2.ErrorResponse("Cannot find country with provided id", nil))
		}
		return c.Status(fiber.StatusInternalServerError).JSON(shared2.ErrorResponse("Failed to delete country", err.Error()))
	}

	return c.Status(fiber.StatusOK).JSON(shared2.SuccessResponse[any]("Successfully delete country", nil))
}
