package postgreSQL

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"os"
	_ "github.com/jackc/pgx/v4/stdlib"
)

const (
	USER     = "natali"
	PASSWORD = "12345"
	DBNAME   = "postgres"
)

func SetupTestDatabase() (testcontainers.Container, *sql.DB) {
	containerReq := testcontainers.ContainerRequest{
		Image:        "postgres:13.3",
		ExposedPorts: []string{"5432/tcp"},
		WaitingFor:   wait.ForListeningPort("5432/tcp"),
		Env: map[string]string{
			"POSTGRES_DB":       DBNAME,
			"POSTGRES_PASSWORD": PASSWORD,
			"POSTGRES_USER":     USER,
		},
	}

	dbContainer, _ := testcontainers.GenericContainer(
		context.Background(),
		testcontainers.GenericContainerRequest{
			ContainerRequest: containerReq,
			Started:          true,
		})

	host, _ := dbContainer.Host(context.Background())
	port, _ := dbContainer.MappedPort(context.Background(), "5432")

	dsnPGConn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port.Int(), USER, PASSWORD, DBNAME)
	db, err := sql.Open("pgx", dsnPGConn)
	if err != nil {
		fmt.Println(err)
		return dbContainer, nil
	}

	err = db.Ping()
	if err != nil {
		fmt.Println(err)
		return dbContainer, nil
	}
	db.SetMaxOpenConns(10)

	text, err := os.ReadFile("../../database/postgreSQL/init.sql")
	if err != nil {
		fmt.Println(err)
		return dbContainer, nil
	}

	if _, err := db.Exec(string(text)); err != nil {
		fmt.Println(err)
		return dbContainer, nil
	}

	return dbContainer, db
}
