package config

import (
	"log"
	"time"

	"strings"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
	"github.com/ulule/limiter/v3"
	"gorm.io/gorm"
)

// Config menyimpan semua konfigurasi aplikasi
type Config struct {
	AppName  string         `mapstructure:"APP_NAME"`
	Version  string         `mapstructure:"VERSION"`
	AppEnv   string         `mapstructure:"APP_ENV"`
	Server   ServerConfig   `mapstructure:",squash"`
	Database DBConfig       `mapstructure:",squash"`
	Redis    RedisConfig    `mapstructure:",squash"`
	Security SecurityConfig `mapstructure:",squash"`
	Logger   LoggerConfig   `mapstructure:",squash"`
}

// ServerConfig menyimpan konfigurasi server
type ServerConfig struct {
	ServerHost     string   `mapstructure:"APP_HOST"`
	ServerPort     string   `mapstructure:"APP_PORT"`
	ServerEnv      string   `mapstructure:"APP_ENV"`
	AllowedOrigins []string `mapstructure:"ALLOWED_ORIGINS"`
}

// RDBMSConfig menyimpan konfigurasi database
type DBConfig struct {
	Enabled         bool          `mapstructure:"ACTIVATE_RDBMS"`
	Driver          string        `mapstructure:"DBDRIVER"`
	Host            string        `mapstructure:"DBHOST"`
	Port            int           `mapstructure:"DBPORT"`
	DB              string        `mapstructure:"DBNAME"`
	User            string        `mapstructure:"DBUSER"`
	Pass            string        `mapstructure:"DBPASS"`
	TimeZone        string        `mapstructure:"DBTIMEZONE"`
	LogLevel        int           `mapstructure:"DBLOGLEVEL"`
	MaxIdleConns    int           `mapstructure:"DBMAXIDLECONNS"`
	MaxOpenConns    int           `mapstructure:"DBMAXOPENCONNS"`
	ConnMaxLifetime time.Duration `mapstructure:"DBCONNMAXLIFETIME"`
	ConnMaxIdleTime time.Duration `mapstructure:"DBCONNMAXIDLETIME"`
	Ssl             DBSsl         `mapstructure:",squash"`
	Client          *gorm.DB
}

// DBSsl menyimpan konfigurasi SSL untuk database
type DBSsl struct {
	Mode       string `mapstructure:"DBSSLMODE"`
	MinTLS     string `mapstructure:"DBSSL_TLS_MIN"`
	RootCA     string `mapstructure:"DBSSL_ROOT_CA"`
	ServerCert string `mapstructure:"DBSSL_SERVER_CERT"`
	ClientCert string `mapstructure:"DBSSL_CLIENT_CERT"`
	ClientKey  string `mapstructure:"DBSSL_CLIENT_KEY"`
}

// RedisConfig menyimpan konfigurasi Redis
type RedisConfig struct {
	Enabled  bool   `mapstructure:"ACTIVATE_REDIS"`
	Host     string `mapstructure:"REDISHOST"`
	Port     int    `mapstructure:"REDISPORT"`
	Password string `mapstructure:"REDISPASSWORD"`
	DB       int    `mapstructure:"REDIS_DB"`
	PoolSize int    `mapstructure:"POOLSIZE"`
	ConnTTL  int    `mapstructure:"CONNTTL"`
	Client   *redis.Client
}

// SecurityConfig menyimpan konfigurasi keamanan aplikasi
type SecurityConfig struct {
	CheckOrigin       bool   `mapstructure:"ACTIVATE_ORIGIN_VALIDATION"`
	RateLimit         string `mapstructure:"RATE_LIMIT"`
	TrustedPlatform   string `mapstructure:"TRUSTED_PLATFORM"`
	ExpectedHost      string `mapstructure:"EXPECTED_HOST"`
	XFrameOptions     string `mapstructure:"X_FRAME_OPTIONS"`
	ContentSecurity   string `mapstructure:"CONTENT_SECURITY_POLICY"`
	XXSSProtection    string `mapstructure:"X_XSS_PROTECTION"`
	StrictTransport   string `mapstructure:"STRICT_TRANSPORT_SECURITY"`
	ReferrerPolicy    string `mapstructure:"REFERRER_POLICY"`
	XContentTypeOpts  string `mapstructure:"X_CONTENT_TYPE_OPTIONS"`
	PermissionsPolicy string `mapstructure:"PERMISSIONS_POLICY"`
	JWTSecretKey      string `mapstructure:"JWT_SECRET_KEY"`
	JWTExpired        int    `mapstructure:"JWT_EXPIRED" envDefault:"2"`
	LimiterInstance   *limiter.Limiter
}

// LoggerConfig menyimpan konfigurasi logger
type LoggerConfig struct {
	LogLevel string `mapstructure:"LOG_LEVEL"`
	Log      *zerolog.Logger
}

// LoadConfig membaca konfigurasi dari environment
func LoadConfig() (*Config, error) {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	// Jika .env tidak ditemukan, gunakan variabel lingkungan
	if err := viper.ReadInConfig(); err != nil {
		log.Printf("No .env file found, using system environment variables: %v", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	// Memastikan AllowedOrigins tidak kosong
	config.Server.AllowedOrigins = strings.Split(viper.GetString("ALLOWED_ORIGINS"), ",")

	return &config, nil
}
