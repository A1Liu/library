package web

import (
	"github.com/A1Liu/library/database"
	"github.com/gin-gonic/gin"
)

func AddImagessApi(images *gin.RouterGroup) {

	images.GET("/get", func(c *gin.Context) {
		id, err := QueryParamUint(c, "id")
		if JsonFail(c, err) {
			return
		}

		image, extension, err := database.GetImage(*id)
		if JsonFail(c, err) {
			return
		}

		c.Writer.Header().Add("Content-Type", "image/"+extension)
		c.Writer.Write(image)
	})
}