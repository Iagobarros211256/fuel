// Comando api sobe o backend do Abastece.
package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Iagobarros211256/abastece/backend/internal/fake"
	"github.com/Iagobarros211256/abastece/backend/internal/httpapi"
	"github.com/Iagobarros211256/abastece/backend/internal/places"
	"github.com/Iagobarros211256/abastece/backend/internal/station"
)

func main() {
	log := slog.New(slog.NewTextHandler(os.Stdout, nil))
	if err := run(log); err != nil {
		log.Error("encerrando com erro", "err", err)
		os.Exit(1)
	}
}

func run(log *slog.Logger) error {
	cfg := loadConfig()

	var provider station.Provider
	switch cfg.provider {
	case "google":
		if cfg.googleKey == "" {
			return errors.New("PROVIDER=google exige GOOGLE_MAPS_API_KEY")
		}
		provider = places.New(cfg.googleKey)
	default:
		provider = fake.Provider{}
	}
	log.Info("provider de postos", "provider", cfg.provider, "cache_ttl", cfg.cacheTTL)

	svc := station.NewService(provider, cfg.cacheTTL)

	srv := &http.Server{
		Addr:              ":" + cfg.port,
		Handler:           httpapi.New(svc, log),
		ReadHeaderTimeout: 5 * time.Second,
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// varre o cache periodicamente para não acumular entradas expiradas
	go func() {
		t := time.NewTicker(cfg.cacheTTL)
		defer t.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-t.C:
				if n := svc.SweepCache(); n > 0 {
					log.Info("cache varrido", "removidos", n)
				}
			}
		}
	}()

	errCh := make(chan error, 1)
	go func() {
		log.Info("ouvindo", "addr", srv.Addr)
		errCh <- srv.ListenAndServe()
	}()

	select {
	case err := <-errCh:
		if !errors.Is(err, http.ErrServerClosed) {
			return err
		}
	case <-ctx.Done():
		log.Info("desligando...")
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return srv.Shutdown(shutdownCtx)
}

type config struct {
	port      string
	provider  string
	googleKey string
	cacheTTL  time.Duration
}

func loadConfig() config {
	cfg := config{
		port:      getenv("PORT", "8080"),
		googleKey: os.Getenv("GOOGLE_MAPS_API_KEY"),
		cacheTTL:  10 * time.Minute,
	}
	// sem chave, cai no fake automaticamente
	defaultProvider := "fake"
	if cfg.googleKey != "" {
		defaultProvider = "google"
	}
	cfg.provider = getenv("PROVIDER", defaultProvider)

	if s := os.Getenv("CACHE_TTL"); s != "" {
		if d, err := time.ParseDuration(s); err == nil && d > 0 {
			cfg.cacheTTL = d
		}
	}
	return cfg
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
