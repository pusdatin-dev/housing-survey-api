package config

import (
	"housing-survey-api/shared"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBConfig
	DBSeed   bool
	Roles    RolesConfig
	Token    string
	AppRole  string
	Resource ResourceConfig
}

type DBConfig struct {
	DBHost string
	DBPort string
	DBUser string
	DBPass string
	DBName string
	AppEnv string
}

type RolesConfig struct {
	SuperAdmin         string
	AdminEselon1       string
	VerificatorEselon1 string
	AdminBalai         string
	VerificatorBalai   string
	Surveyor           string
}

type ResourceConfig struct {
	TagNegara       string
	TagPengembang   string
	TagSwadaya      string
	TagGotongRoyong string
}

func LoadConfig() *Config {
	// Load .env if exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	dbConfig := DBConfig{
		DBHost: getEnv("DB_HOST", "localhost"),
		DBPort: getEnv("DB_PORT", "5432"),
		DBUser: getEnv("DB_USER", "survey_user"),
		DBPass: getEnv("DB_PASS", "survey_pass"),
		DBName: getEnv("DB_NAME", "survey_db"),
		AppEnv: getEnv("APP_ENV", "development"),
	}

	rolesConfig := RolesConfig{
		SuperAdmin:         getEnv("ROLE_SUPER_ADMIN", "Super Admin"),
		AdminEselon1:       getEnv("ROLE_ADMIN_ESELON_1", "Admin Eselon 1"),
		VerificatorEselon1: getEnv("ROLE_VERIFICATOR_ESELON_1", "Verificator Eselon 1"),
		AdminBalai:         getEnv("ROLE_ADMIN_BALAI", "Admin Balai"),
		VerificatorBalai:   getEnv("ROLE_VERIFICATOR_BALAI", "Verificator Balai"),
		Surveyor:           getEnv("ROLE_SURVEYOR", "Surveyor"),
	}

	resConfig := ResourceConfig{
		TagNegara:       getEnv("RESOURCE_NEGARA", shared.TagNegara),
		TagPengembang:   getEnv("RESOURCE_PENGEMBANG", shared.TagPengembang),
		TagSwadaya:      getEnv("RESOURCE_SWADAYA", shared.TagSwadaya),
		TagGotongRoyong: getEnv("RESOURCE_GOTONGROYONG", shared.TagGotongRoyong),
	}

	return &Config{
		DBConfig: dbConfig,
		DBSeed:   getEnv("DB_SEED", "false") == "true",
		Roles:    rolesConfig,
		Token:    getEnv("JWT_SECRET", ""),
		AppRole:  getEnv("APP_ROLE", "api"),
		Resource: resConfig,
	}
}

func getEnv(key, fallback string) string {
	if val, exists := os.LookupEnv(key); exists && val != "" {
		return val
	}
	return fallback
}
