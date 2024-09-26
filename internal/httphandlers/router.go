package httphandlers

import (
	"github.com/gin-gonic/gin"
)

type Router interface {
	Run() error
}

type RouterCfg struct {
	ServerAddress string `arg:"--server-address,env:SERVER_ADDRESS" default:":8080" help:"Address and port of this service"`
	FormFilePath  string `arg:"--form-file-path,env:FORM_FILE_PATH" default:"web/templates/form.html" help:"Path to HTML file with form"`
}

type routerImpl struct {
	cfg RouterCfg
}

func NewRouter(cfg RouterCfg) Router {
	return &routerImpl{
		cfg: cfg,
	}
}

var _ Router = (*routerImpl)(nil)

func (r *routerImpl) Run() error {
	router := gin.Default()

	ph := pageHandler{
		formFilePath: r.cfg.FormFilePath,
	}

	router.GET("/", ph.Web)
	router.POST("/", ph.Process)
	// router.GET("/stat", ph.Stat)

	if err := router.Run(r.cfg.ServerAddress); err != nil {
		return err
	}

	return nil
}
