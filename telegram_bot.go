package main

import (
	"context"
	"fmt"
	"telegram_bot/config"
	"telegram_bot/internal/infra/cache"
	"telegram_bot/pkg/logger"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type User struct {
	ID           int
	TelegramID   int
	SearchParams string
}

type Knife struct {
	ID             int
	Model          string
	Price          float64
	URL            string
	MarketplaceURL string
}

type Config struct {
	Database struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Name     string `yaml:"name"`
		SSLMode  string `yaml:"sslmode"`
	} `yaml:"database"`
}

type ScraperManager struct {
	scrapers map[string]Scraper
}
type SearchService struct {
	cache *cache.RedisCache
	/*https://forest-home.ru/
	https://www.lamnia.com/ru
	https://www.nozhikov.ru
	https://www.euro-knife.com/
	*/
}

type RedisCache struct {
	client *redis.Client
	ttl    time.Duration
}

type Postgres struct {
	Pool *pgxpool.Pool
}

type UserRepository struct {
	db *Postgres
}

type UserService struct {
	userRepo *database.UserRepository
}

type UserRepository interface {
	SaveUser(ctx context.Context, user *User) error
	GetUser(ctx context.Context, telegramID int) (*User, error)
}

type Scraper interface {
	Scrape(query string) ([]model.Knife, error)
}

func NewRedisCache(addr string, password string, db int, ttl time.Duration) (*RedisCache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
}
func HandleSearch(update tgbotapi.Update) {
	// Валидация запроса
	// Запуск поиска через service.Search
	// Форматирование ответа
}

func main() {
	// Загрузка конфигурации
	cfg := config.Load()

	// Инициализация логгера
	logger.Init(cfg.Debug)

	// Подключение к Redis
	redisClient, err := cache.NewRedisClient(
		cfg.Redis.Address,
		cfg.Redis.Password,
		cfg.Redis.DB,
	)
	if err != nil {
		logger.Fatal("Failed to connect to Redis", err)
	}
	defer redisClient.Close()

	// Инициализация и запуск бота
	bot := initBot(redisClient, cfg)
	bot.Start()
}

func (c *Config) GetDBConnectionString() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Database.Host,
		c.Database.Port,
		c.Database.User,
		c.Database.Password,
		c.Database.Name,
		c.Database.SSLMode,
	)
}

func NewPostgres(connString string) (*Postgres, error) {
	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("failed to parse connection string: %v", err)
	}

	// Настройка пула соединений
	config.MaxConns = 20
	config.MinConns = 5
	config.MaxConnLifetime = time.Hour
	config.MaxConnIdleTime = time.Minute * 30
	config.HealthCheckPeriod = time.Minute

	// Создание пула соединений
	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %v", err)
	}

	// Проверка подключения
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	logger.Info("Successfully connected to PostgreSQL")

	return &Postgres{Pool: pool}, nil
}

func (p *Postgres) Close() {
	p.Pool.Close()
}
func NewUserRepository(db *Postgres) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateOrUpdate(ctx context.Context, user *model.User) error {
	query := `
		INSERT INTO users (telegram_id, username, first_name, last_name, language_code, is_premium)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (telegram_id) 
		DO UPDATE SET 
			username = EXCLUDED.username,
			first_name = EXCLUDED.first_name,
			last_name = EXCLUDED.last_name,
			language_code = EXCLUDED.language_code,
			is_premium = EXCLUDED.is_premium,
			updated_at = CURRENT_TIMESTAMP
		RETURNING id, created_at, updated_at
	`

	err := r.db.Pool.QueryRow(ctx, query,
		user.TelegramID,
		user.Username,
		user.FirstName,
		user.LastName,
		user.LanguageCode,
		user.IsPremium,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create or update user: %v", err)
	}

	return nil
}

func (r *UserRepository) GetByTelegramID(ctx context.Context, telegramID int64) (*model.User, error) {
	query := `
		SELECT id, telegram_id, username, first_name, last_name, language_code, is_premium, created_at, updated_at
		FROM users
		WHERE telegram_id = $1
	`

	user := &model.User{}
	err := r.db.Pool.QueryRow(ctx, query, telegramID).Scan(
		&user.ID,
		&user.TelegramID,
		&user.Username,
		&user.FirstName,
		&user.LastName,
		&user.LanguageCode,
		&user.IsPremium,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user: %v", err)
	}

	return user, nil
}
func NewUserService(userRepo *database.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) ProcessUserFromTelegram(ctx context.Context, telegramUser *tgbotapi.User) (*model.User, error) {
	user := &model.User{
		TelegramID:   int64(telegramUser.ID),
		Username:     telegramUser.UserName,
		FirstName:    telegramUser.FirstName,
		LastName:     telegramUser.LastName,
		LanguageCode: telegramUser.LanguageCode,
		IsPremium:    telegramUser.IsPremium,
	}

	if err := s.userRepo.CreateOrUpdate(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}
