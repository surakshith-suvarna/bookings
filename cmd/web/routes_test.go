package main

import (
	"fmt"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/surakshith-suvarna/bookings/internal/config"
)

func TestRoutes(t *testing.T) {
	var app config.AppConfig

	mux := routes(&app)

	switch v := mux.(type) {
	case *chi.Mux:
		//return nothing. Test passed
	default:
		t.Error(fmt.Sprintf("The type is not chi.mux but %T", v))
	}
}
