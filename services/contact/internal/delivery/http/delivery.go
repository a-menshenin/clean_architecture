package http

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"

	"architecture_go_2/services/contact/internal/useCase"
)

// @title slurm contact service on clean architecture
// @version 1.0
// @description contact service on clean architecture
// @license.name mensheninao

// @contact.name API Support
// @contact.email mensheninao@gmail.com

// @BasePath /

func init() {
	viper.SetConfigName(".env")
	viper.SetConfigType("dotenv")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	viper.SetDefault("HTTP_PORT", 80)
}

type Delivery struct {
	ucContact useCase.Contact
	ucGroup   useCase.Group
	router    *gin.Engine

	options options
}

type options struct {
	Timeout time.Duration
}

type Option func(*options)

func WithTimeout(timeout time.Duration) Option {
	return func(args *options) {
		args.Timeout = timeout
	}
}

func (d *Delivery) SetOptions(setters ...Option) {
	args := &options{
		Timeout: time.Second * 30,
	}

	for _, setter := range setters {
		setter(args)
	}
	d.options = *args
}

func New(ucContact useCase.Contact, ucGroup useCase.Group, setters ...Option) *Delivery {
	var d = &Delivery{
		ucContact: ucContact,
		ucGroup:   ucGroup,
	}

	d.SetOptions(setters...)

	d.router = d.initRouter()
	return d
}

func (d *Delivery) Run() error {
	return d.router.Run(fmt.Sprintf(":%d", uint16(viper.GetUint("HTTP_PORT"))))
}

func checkAuth(c *gin.Context) {
	c.Next()
}
