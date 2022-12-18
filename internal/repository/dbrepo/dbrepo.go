package dbrepo

import (
	"database/sql"

	"github.com/budidak/booking/internal/config"
	"github.com/budidak/booking/internal/repository"
)

// Burada da database connection için repository pattern kullanacağız, bunu yapmamızın sebebi: Codebase içinde bir değişiklik yapmamız gerektiğinde olabildiğince az yerden müdahale ederek bu işlemi yapmak.
// Örneğin postgres yerine başka bir db kullanmak istersek mysql, mariadb vs. O zaman aşağıdaki template gibi onlar için de yazmalıyız.

// template: newrepo oluşturmak için kullanılacak
type postgresDBRepo struct {
	App *config.AppConfig
	DB  *sql.DB
}

// initialization bununla yapılacak mainde
func NewPostgresRepo(conn *sql.DB, a *config.AppConfig) repository.DatabaseRepo {
	return &postgresDBRepo{
		App: a,
		DB:  conn,
	}
}
