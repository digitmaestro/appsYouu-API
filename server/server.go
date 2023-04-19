package server

import (
	"context"
	"net/http"
	"strings"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"github.com/caarlos0/env"
	cron "github.com/robfig/cron/v3"
	"github.com/rs/cors"

	"github.com/steebchen/keskin-api/handlers/template/index"
	"github.com/steebchen/keskin-api/handlers/template/webmanifest"
	"github.com/steebchen/keskin-api/i18n"
	"github.com/steebchen/keskin-api/lib/file"
	"github.com/steebchen/keskin-api/lib/share"
)

type Config struct {
	Port                    string   `env:"PORT" envDefault:"4000"`
	DataFolder              string   `env:"DATA_FOLDER" envDefault:"/data/images/"`
	AllowedOrigins          []string `env:"KESKIN_ORIGIN" envDefault:"http://localhost:4000" envSeparator:";"`
	EnableCORS              bool     `env:"ENABLE_CORS"`
	DefaultLanguage         string   `env:"DEFAULT_LANGUAGE" envDefault:"DE"`
	DefaultShareRedirectUrl string   `env:"DEFAULT_SHARE_REDIRECT_URL" envDefault:"https://keskin.luca-steeb.com/"`
	HostName                string   `env:"HOST_NAME" envDefault:"keskin.luca-steeb.com"`
	TemplateFolder          string   `env:"TEMPLATE_FOLDER" envDefault:"/templates/"`
	GcmSenderId             string   `env:"GCM_SENDER_ID"`
}

type Server struct {
	Mux      *http.ServeMux
	Config   *Config
	CronJobs *cron.Cron
}

type MuxHandler struct {
	Next *http.ServeMux
}

func (h *MuxHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.Next.ServeHTTP(w, r)
}

func (s *Server) Listen() error {
	file.BaseDir = s.Config.DataFolder
	share.DefaultShareRedirectUrl = s.Config.DefaultShareRedirectUrl
	share.BaseHostName = s.Config.HostName
	index.TemplateFolder = s.Config.TemplateFolder
	webmanifest.GcmSenderId = s.Config.GcmSenderId

	envDefaultLanguage := strings.ToUpper(s.Config.DefaultLanguage)
	if i18n.IsAllowedLanguage(envDefaultLanguage) {
		i18n.DefaultLanguage = envDefaultLanguage
	}

	addr := ":" + s.Config.Port
	if s.Config.EnableCORS {
		c := cors.New(cors.Options{
			AllowCredentials: true,
			AllowedOrigins:   s.Config.AllowedOrigins,
		})
		handler := c.Handler(&MuxHandler{s.Mux})
		return http.ListenAndServe(addr, handler)
	} else {
		return http.ListenAndServe(addr, &MuxHandler{s.Mux})
	}
}

func NewConfig() (*Config, error) {
	config := &Config{}
	err := env.Parse(config)
	return config, err
}

func NewServer(mux *http.ServeMux, config *Config, cronJobs *cron.Cron) *Server {
	return &Server{
		Mux:      mux,
		Config:   config,
		CronJobs: cronJobs,
	}
}

func NewFirebaseApp() (*firebase.App, error) {
	return firebase.NewApp(context.Background(), nil)
}

func NewFirebaseMessagingClient(app *firebase.App) (*messaging.Client, error) {
	return app.Messaging(context.Background())
}
