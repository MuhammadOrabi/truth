package config

import (
	"fmt"
	"os"
)

func GetDBType() string {
	return os.Getenv("DB_TYPE")
}

func GetPostgresConnectionString() string {
	dataBase := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PASSWORD"))
	return dataBase
}

func GetTestingPostgresConnectionString() string {
	dataBase := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		os.Getenv("DB_TESTING_HOST"),
		os.Getenv("DB_TESTING_PORT"),
		os.Getenv("DB_TESTING_USER"),
		os.Getenv("DB_TESTING_NAME"),
		os.Getenv("DB_TESTING_PASSWORD"))
	return dataBase
}
