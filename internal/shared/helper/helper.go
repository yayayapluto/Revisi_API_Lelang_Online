package shared

import (
	"github.com/gofiber/fiber/v2"
	"strconv"
)

func ParseIdParam(c *fiber.Ctx, param string) (uint, error) {
	idStr := c.Params(param)
	id64, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(id64), nil
}

func GetQueryInt(c *fiber.Ctx, key string, defaultValue int) int {
	if v := c.Query(key); v != "" {
		if parsed, err := strconv.Atoi(v); err == nil {
			return parsed
		}
	}
	return defaultValue
}

func BuildPageURL(c *fiber.Ctx, page int) string {
	url := c.BaseURL() + c.Path()
	if page != 0 {
		url = url + "?page=" + strconv.Itoa(page)
	}
	return url
}
