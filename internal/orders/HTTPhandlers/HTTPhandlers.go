package HTTPHandlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"test_wb/internal/cache"

	"github.com/labstack/echo/v4"
)

func OrderHandler(c echo.Context, os *cache.OrderCache) error {
	response, err := json.MarshalIndent(os.GetOrderByUid(c.Param("uid")), "", "\t")

	if err != nil {
		fmt.Printf("Error at marshaling: %v", err)
		return err
	}

	return c.JSONBlob(http.StatusOK, response)
}
