package page

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/chrootlogin/go-wiki/src/repo"
	"github.com/chrootlogin/go-wiki/src/common"
	"os"
)

type ApiRequest struct {
	Title   string `json:"title,omitempty"`
	Path	string `json:"path,omitempty"`
	Content string `json:"content,omitempty"`
}

type apiResponse struct {
	ContentType string `json:"contentType,omitempty"`
	Content     string `json:"content,omitempty"`
}

// READ
func GetPageHandler(c *gin.Context) {
	pagePath := c.Param("page")
	if pagePath != "" {
		pagePath = trimLeftChar(pagePath)
	}

	file, err := repo.GetFile(pagePath)
	if err != nil {
		if os.IsNotExist(err) {
			c.JSON(http.StatusNotFound, common.ApiResponse{ Message: "Not found" })
			return
		}
		c.JSON(http.StatusInternalServerError, common.ApiResponse{ Message: err.Error() })
		return
	}

	c.JSON(http.StatusOK, apiResponse{
		Content: file.Content,
	})
}

func trimLeftChar(s string) string {
	for i := range s {
		if i > 0 {
			return s[i:]
		}
	}
	return s[:0]
}