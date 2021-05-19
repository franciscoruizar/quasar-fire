package topsecret

import (
	"errors"
	"net/http"

	domain "github.com/franciscoruizar/quasar-fire/internal/domain"
	"github.com/franciscoruizar/quasar-fire/internal/infrastructure/server/handler"
	"github.com/franciscoruizar/quasar-fire/internal/usecases"
	"github.com/franciscoruizar/quasar-fire/internal/usecases/dto"
	"github.com/gin-gonic/gin"
)

type CreateRequests struct {
	Satellites []CreateRequest `json:"satellites" binding:"required"`
}

type CreateRequest struct {
	Name     string   `json:"name" binding:"required"`
	Distance float64  `json:"distance" binding:"required"`
	Message  []string `json:"message" binding:"required"`
}

type CreateResponse struct {
	Position dto.PositionResponse `json:"position" binding:"required"`
	Message  string               `json:"message" binding:"required"`
}

func TopSecretCreateHandler(service usecases.TopSecretCreator) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req CreateRequests
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
		}

		var satellites []usecases.TopSecretCreatorRequest
		for i := 0; i < len(req.Satellites); i++ {
			satellites = append(satellites, usecases.TopSecretCreatorRequest{
				Name:      req.Satellites[i].Name,
				Dinstance: req.Satellites[i].Distance,
				Message:   req.Satellites[i].Message,
			})
		}

		response, err := service.Create(satellites)

		if err != nil {
			errorMessage := handler.Error{
				Message: err.Error(),
			}
			switch {
			case errors.Is(err, domain.ErrInvalidSateliteID):
				c.JSON(http.StatusBadRequest, errorMessage)
				return
			default:
				c.JSON(http.StatusInternalServerError, errorMessage)
				return
			}
		}

		responseMap := CreateResponse{
			Message:  response.Message,
			Position: response.Position,
		}

		c.Status(http.StatusCreated)

		c.JSON(http.StatusCreated, responseMap)
	}
}