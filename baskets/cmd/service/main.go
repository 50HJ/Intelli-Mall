package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	_ "github.com/jackc/pgx/v4/stdlib"

	"github.com/50HJ/Intelli-Mall/baskets"
	"github.com/50HJ/Intelli-Mall/baskets/migrations"
	"github.com/50HJ/Intelli-Mall/internal/config"
	"github.com/50HJ/Intelli-Mall/internal/system"
	"github.com/50HJ/Intelli-Mall/internal/web"
)

func main() {
	if err := run(); err != nil {
		fmt.Printf("baskets exitted abnormally: %s\n", err)
		os.Exit(1)
	}
}

func run() (err error) {
	var cfg config.AppConfig
	cfg, err = config.InitConfig()
	if err != nil {
		return err
	}
	s, err := system.NewSystem(cfg)
	if err != nil {
		return err
	}
	defer func(db *sql.DB) {
		if err = db.Close(); err != nil {
			return
		}
	}(s.DB())
	if err = s.MigrateDB(migrations.FS); err != nil {
		return err
	}
	s.Mux().Mount("/", http.FileServer(http.FS(web.WebUI)))
	// call the module composition root
	if err = baskets.Root(s.Waiter().Context(), s); err != nil {
		return err
	}

	fmt.Println("started baskets service")
	defer fmt.Println("stopped baskets service")

	s.Waiter().Add(
		s.WaitForWeb,
		s.WaitForRPC,
		s.WaitForStream,
	)

	return s.Waiter().Wait()
}
