package controller

import "net/http"

type Controller struct {
	serveMux    *http.ServeMux
	application IApplication
}

func NewController(app IApplication) *Controller {
	controller := &Controller{serveMux: http.NewServeMux(), application: app}
	controller.serveMux.HandleFunc("/get", controller.Get)
	controller.serveMux.HandleFunc("/set", controller.Set)
	return controller
}

func (c *Controller) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c.serveMux.ServeHTTP(w, r)
}

func (c *Controller) Get(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	if key == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	value, err := c.application.Get(key)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write([]byte(value))
}

func (c *Controller) Set(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	value := r.URL.Query().Get("value")
	if key == "" || value == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err := c.application.Set(key, value)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

type IApplication interface {
	Get(key string) (string, error)
	Set(key, value string) error
}
