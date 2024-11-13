package engine

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type SEngine struct {
	Gin *gin.Engine
	Srv *http.Server
}

type FHandlers func(*SEngine) error

func (e *SEngine) WithHandlers(handlers ...FHandlers) (*SEngine, error) {
	for _, handler := range handlers {
		if handler == nil {
			continue
		}

		if err := handler(e); err != nil {
			return nil, err
		}
	}
	return e, nil
}

func BuildEngineByGIN() *SEngine {
	engine := &SEngine{
		Gin: gin.New(),
	}
	return engine
}
