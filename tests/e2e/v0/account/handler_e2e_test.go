package account_e2e_test

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/go-testfixtures/testfixtures/v3"
	"github.com/google/uuid"
	appconfig "github.com/mandarine-io/Backend/internal/config"
	"github.com/mandarine-io/Backend/internal/domain/dto"
	"github.com/mandarine-io/Backend/internal/helper/security"
	model2 "github.com/mandarine-io/Backend/internal/persistence/model"
	http2 "github.com/mandarine-io/Backend/internal/transport/http"
	dto2 "github.com/mandarine-io/Backend/pkg/transport/http/dto"
	validator2 "github.com/mandarine-io/Backend/pkg/transport/http/validator"
	"github.com/mandarine-io/Backend/tests/e2e"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"
)

var (
	testEnvironment *e2e.TestEnvironment
	server          *httptest.Server
	fixtures        *testfixtures.Loader
)

type mailhogMessagesResponse struct {
	Total    int           `json:"total"`
	Start    int           `json:"start"`
	Count    int           `json:"count"`
	Messages []interface{} `json:"messages"`
}

func TestMain(m *testing.M) {
	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	cfg := &appconfig.Config{
		Server: appconfig.ServerConfig{
			Name:           "mandarine_test",
			Mode:           "test",
			ExternalOrigin: "http://localhost:8081",
			Port:           8081,
			Version:        "0.0.0",
			RPS:            100,
			MaxRequestSize: 524288000,
		},
		Database: appconfig.DatabaseConfig{
			Type: "postgres",
			Postgres: &appconfig.PostgresDatabaseConfig{
				Username: "mandarine",
				Password: "password",
				DBName:   "mandarine_test",
			},
		},
		Cache: appconfig.CacheConfig{
			TTL:  120,
			Type: "redis",
			Redis: &appconfig.RedisCacheConfig{
				Username: "default",
				Password: "password",
				DBIndex:  0,
			},
		},
		PubSub: appconfig.PubSubConfig{
			Type: "redis",
			Redis: &appconfig.RedisPubSubConfig{
				Username: "default",
				Password: "password",
				DBIndex:  0,
			},
		},
		S3: appconfig.S3Config{
			Type: "minio",
			Minio: &appconfig.MinioS3Config{
				AccessKey: "admin",
				SecretKey: "Password_10",
				Bucket:    "mandarine-test",
			},
		},
		SMTP: appconfig.SmtpConfig{
			Host:     "127.0.0.1",
			Port:     25,
			Username: "admin",
			Password: "password",
			From:     "Mandarine <admin@localhost>",
			SSL:      false,
		},
		Locale: appconfig.LocaleConfig{
			Path:     pwd + "/../../../../locales",
			Language: "ru",
		},
		Template: appconfig.TemplateConfig{
			Path: pwd + "/../../../../templates",
		},
		Migrations: appconfig.MigrationConfig{
			Path: pwd + "/../../../../migrations",
		},
		Logger: appconfig.LoggerConfig{
			Level: "debug",
			Console: appconfig.ConsoleLoggerConfig{
				Enable:   true,
				Encoding: "text",
			},
			File: appconfig.FileLoggerConfig{
				Enable: false,
			},
		},
		OAuthClients: map[string]appconfig.OauthClientConfig{
			"google": {
				ClientID:     "",
				ClientSecret: "",
			},
			"yandex": {
				ClientID:     "",
				ClientSecret: "",
			},
			"mailru": {
				ClientID:     "",
				ClientSecret: "",
			},
		},
		Security: appconfig.SecurityConfig{
			JWT: appconfig.JWTConfig{
				Secret:          "",
				AccessTokenTTL:  3600,
				RefreshTokenTTL: 86400,
			},
			OTP: appconfig.OTPConfig{
				Length: 6,
				TTL:    300,
			},
		},
		Websocket: appconfig.WebsocketConfig{
			PoolSize: 1024,
		},
	}

	// Initialize test environment
	testEnvironment = e2e.NewTestContainer()
	defer testEnvironment.Close()
	testEnvironment.MustInitialize(cfg)

	// Setup routes
	router := http2.SetupRouter(testEnvironment.Container)
	// Setup validators
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("pastdate", validator2.PastDateValidator)
		_ = v.RegisterValidation("zxcvbn", validator2.ZxcvbnPasswordValidator)
		_ = v.RegisterValidation("username", validator2.UsernameValidator)
		_ = v.RegisterValidation("point", validator2.PointValidator)
	}

	// Create server
	server = httptest.NewServer(router)
	defer server.Close()

	fixtures = e2e.MustNewFixtures(testEnvironment.Container.DB, pwd+"/fixtures/users.yml")

	os.Exit(m.Run())
}

func Test_AccountHandler_GetAccount(t *testing.T) {
	e2e.MustLoadFixtures(fixtures)

	userEntity := &model2.UserEntity{
		ID:        uuid.MustParse("a02fc7e1-c19a-4c1a-b66e-29fed1ed452f"),
		Username:  "user1",
		Email:     "user1@example.com",
		Password:  "$2a$12$4XWfvkfvvLxLlLyPQ9CA7eNhkUIFSj7sF3768lAMJi9G2kl4XjGve",
		Role:      model2.RoleEntity{Name: model2.RoleUser},
		IsEnabled: true,
		DeletedAt: nil,
	}
	accessToken, _, _ := security.GenerateTokens(testEnvironment.Container.Config.Security.JWT, userEntity)

	url := server.URL + "/v0/account"
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)

	t.Run("Unauthorized", func(t *testing.T) {
		req, _ := http.NewRequest("GET", url, nil)
		resp, err := server.Client().Do(req)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("Banned user", func(t *testing.T) {
		anotherUserEntity := model2.UserEntity{
			ID:        uuid.New(),
			Username:  "user",
			Email:     "user@example.com",
			Role:      model2.RoleEntity{Name: model2.RoleUser},
			IsEnabled: false,
			DeletedAt: nil,
		}
		anotherAccessToken, _, _ := security.GenerateTokens(testEnvironment.Container.Config.Security.JWT, &anotherUserEntity)
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("Authorization", "Bearer "+anotherAccessToken)

		resp, err := server.Client().Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	})

	t.Run("Deleted user", func(t *testing.T) {
		deletedTime := time.Now().UTC()
		anotherUserEntity := model2.UserEntity{
			ID:        uuid.New(),
			Username:  "user",
			Email:     "user@example.com",
			Role:      model2.RoleEntity{Name: model2.RoleUser},
			IsEnabled: true,
			DeletedAt: &deletedTime,
		}
		anotherAccessToken, _, _ := security.GenerateTokens(testEnvironment.Container.Config.Security.JWT, &anotherUserEntity)
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("Authorization", "Bearer "+anotherAccessToken)

		resp, err := server.Client().Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	})

	t.Run("User not found", func(t *testing.T) {
		accessToken, _, _ := security.GenerateTokens(
			testEnvironment.Container.Config.Security.JWT,
			&model2.UserEntity{
				ID:        uuid.New(),
				IsEnabled: true,
				DeletedAt: nil,
			})
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("Authorization", "Bearer "+accessToken)

		resp, err := server.Client().Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("Success", func(t *testing.T) {
		resp, err := server.Client().Do(req)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var body dto.AccountOutput
		err = e2e.ReadResponseBody(resp, &body)
		assert.NoError(t, err)
		assert.Equal(t, userEntity.Email, body.Email)
		assert.Equal(t, userEntity.Username, body.Username)
		assert.True(t, body.IsEnabled)
		assert.True(t, body.IsEmailVerified)
		assert.False(t, body.IsPasswordTemp)
		assert.False(t, body.IsDeleted)
	})
}

func Test_AccountHandler_UpdateUsername(t *testing.T) {
	e2e.MustLoadFixtures(fixtures)

	url := server.URL + "/v0/account/username"

	userEntity := &model2.UserEntity{
		ID:        uuid.MustParse("c51013eb-d179-4f14-90da-0e9ac732ae86"),
		Username:  "user_for_update_username_1",
		Email:     "user_for_update_username_1@example.com",
		Password:  "$2a$12$jVO1hn15BIlyXNvm5sgUoOnGpjLMsUR654fv5qibD7AeW1XmZ7XNq",
		Role:      model2.RoleEntity{Name: model2.RoleUser},
		IsEnabled: true,
		DeletedAt: nil,
	}
	accessToken, _, _ := security.GenerateTokens(testEnvironment.Container.Config.Security.JWT, userEntity)

	t.Run("Invalid body", func(t *testing.T) {
		bodies := []*dto.UpdateUsernameInput{
			nil,

			// Bad username
			{
				Username: "User",
			},
			{
				Username: "1user",
			},
			{
				Username: "user!",
			},
			{
				Username: "",
			},
			{
				Username: strings.Repeat("user", 256),
			},
		}

		for i, body := range bodies {
			t.Run(fmt.Sprintf("Bad body %d", i), func(t *testing.T) {
				var reqBodyReader io.Reader = nil
				if body != nil {
					reqBodyReader, _ = e2e.NewJSONReader(body)
				}

				req, _ := http.NewRequest("PATCH", url, reqBodyReader)
				req.Header.Set("Authorization", "Bearer "+accessToken)

				resp, err := server.Client().Do(req)
				assert.NoError(t, err)
				assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

				var body dto2.ErrorResponse
				err = e2e.ReadResponseBody(resp, &body)
				assert.NoError(t, err)
			})
		}

	})

	t.Run("Unauthorized", func(t *testing.T) {
		reqBody := dto.UpdateUsernameInput{
			Username: "new_user_for_update_username",
		}
		reqBodyReader, _ := e2e.NewJSONReader(reqBody)

		req, _ := http.NewRequest("PATCH", url, reqBodyReader)
		resp, err := server.Client().Do(req)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

		var body dto2.ErrorResponse
		err = e2e.ReadResponseBody(resp, &body)
		assert.NoError(t, err)
	})

	t.Run("Banned user", func(t *testing.T) {
		anotherUserEntity := model2.UserEntity{
			ID:        uuid.New(),
			Username:  "user",
			Email:     "user@example.com",
			Role:      model2.RoleEntity{Name: model2.RoleUser},
			IsEnabled: false,
			DeletedAt: nil,
		}
		anotherAccessToken, _, _ := security.GenerateTokens(testEnvironment.Container.Config.Security.JWT, &anotherUserEntity)
		req, _ := http.NewRequest("PATCH", url, nil)
		req.Header.Set("Authorization", "Bearer "+anotherAccessToken)

		resp, err := server.Client().Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	})

	t.Run("Deleted user", func(t *testing.T) {
		deletedTime := time.Now().UTC()
		anotherUserEntity := model2.UserEntity{
			ID:        uuid.New(),
			Username:  "user",
			Email:     "user@example.com",
			Role:      model2.RoleEntity{Name: model2.RoleUser},
			IsEnabled: true,
			DeletedAt: &deletedTime,
		}
		anotherAccessToken, _, _ := security.GenerateTokens(testEnvironment.Container.Config.Security.JWT, &anotherUserEntity)
		req, _ := http.NewRequest("PATCH", url, nil)
		req.Header.Set("Authorization", "Bearer "+anotherAccessToken)

		resp, err := server.Client().Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	})

	t.Run("User not found", func(t *testing.T) {
		accessToken, _, _ := security.GenerateTokens(
			testEnvironment.Container.Config.Security.JWT,
			&model2.UserEntity{
				ID:        uuid.New(),
				IsEnabled: true,
				DeletedAt: nil,
			},
		)
		reqBody := dto.UpdateUsernameInput{
			Username: "new_user_for_update_username",
		}
		reqBodyReader, _ := e2e.NewJSONReader(reqBody)

		req, _ := http.NewRequest("PATCH", url, reqBodyReader)
		req.Header.Set("Authorization", "Bearer "+accessToken)

		resp, err := server.Client().Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)

		var body dto2.ErrorResponse
		err = e2e.ReadResponseBody(resp, &body)
		assert.NoError(t, err)
	})

	t.Run("Username not changed", func(t *testing.T) {
		reqBody := dto.UpdateUsernameInput{
			Username: "user_for_update_username_1",
		}
		reqBodyReader, _ := e2e.NewJSONReader(reqBody)

		req, _ := http.NewRequest("PATCH", url, reqBodyReader)
		req.Header.Set("Authorization", "Bearer "+accessToken)

		resp, err := server.Client().Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var body dto.AccountOutput
		err = e2e.ReadResponseBody(resp, &body)
		assert.NoError(t, err)
		assert.Equal(t, reqBody.Username, body.Username)
	})

	t.Run("Username already in use", func(t *testing.T) {
		reqBody := dto.UpdateUsernameInput{
			Username: "user_for_update_username_2",
		}
		reqBodyReader, _ := e2e.NewJSONReader(reqBody)

		req, _ := http.NewRequest("PATCH", url, reqBodyReader)
		req.Header.Set("Authorization", "Bearer "+accessToken)

		resp, err := server.Client().Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusConflict, resp.StatusCode)

		var body dto2.ErrorResponse
		err = e2e.ReadResponseBody(resp, &body)
		assert.NoError(t, err)
	})

	t.Run("Success", func(t *testing.T) {
		reqBody := dto.UpdateUsernameInput{
			Username: "new_user_for_update_username",
		}
		reqBodyReader, _ := e2e.NewJSONReader(reqBody)

		req, _ := http.NewRequest("PATCH", url, reqBodyReader)
		req.Header.Set("Authorization", "Bearer "+accessToken)

		resp, err := server.Client().Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var body dto.AccountOutput
		err = e2e.ReadResponseBody(resp, &body)
		assert.NoError(t, err)
		assert.Equal(t, reqBody.Username, body.Username)
	})
}

func Test_AccountHandler_UpdateEmail(t *testing.T) {
	e2e.MustLoadFixtures(fixtures)

	url := server.URL + "/v0/account/email"
	userEntity := &model2.UserEntity{
		ID:        uuid.MustParse("d76d9da5-ff66-4397-99cb-8b0298e887bd"),
		Username:  "user_for_update_email_1",
		Email:     "user_for_update_email_1@example.com",
		Password:  "$2a$12$cfjXzgolp1b2sdoP7RNX4ui/cLtoHrTUF.c7JIphuPNWCPVB1s3.2",
		Role:      model2.RoleEntity{Name: model2.RoleUser},
		IsEnabled: true,
		DeletedAt: nil,
	}
	accessToken, _, _ := security.GenerateTokens(testEnvironment.Container.Config.Security.JWT, userEntity)

	t.Run("Bad body", func(t *testing.T) {
		bodies := []*dto.UpdateEmailInput{
			nil,

			// Bad email
			{
				Email: "",
			},
			{
				Email: "user_for_update_email_1example.com",
			},
		}

		for i, body := range bodies {
			t.Run(fmt.Sprintf("Bad body %d", i), func(t *testing.T) {
				var reqBodyReader io.Reader = nil
				if body != nil {
					reqBodyReader, _ = e2e.NewJSONReader(body)
				}

				req, _ := http.NewRequest("PATCH", url, reqBodyReader)
				req.Header.Set("Authorization", "Bearer "+accessToken)

				resp, err := server.Client().Do(req)
				assert.NoError(t, err)
				assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

				var body dto2.ErrorResponse
				err = e2e.ReadResponseBody(resp, &body)
				assert.NoError(t, err)
			})
		}
	})

	t.Run("Unauthorized", func(t *testing.T) {
		reqBody := dto.UpdateEmailInput{
			Email: "new_user_for_update_email_1@example.com",
		}
		reqBodyReader, _ := e2e.NewJSONReader(reqBody)

		req, _ := http.NewRequest("PATCH", url, reqBodyReader)

		resp, err := server.Client().Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

		var body dto2.ErrorResponse
		err = e2e.ReadResponseBody(resp, &body)
		assert.NoError(t, err)
	})

	t.Run("Banned user", func(t *testing.T) {
		anotherUserEntity := model2.UserEntity{
			ID:        uuid.New(),
			Username:  "user",
			Email:     "user@example.com",
			Role:      model2.RoleEntity{Name: model2.RoleUser},
			IsEnabled: false,
			DeletedAt: nil,
		}
		anotherAccessToken, _, _ := security.GenerateTokens(testEnvironment.Container.Config.Security.JWT, &anotherUserEntity)
		req, _ := http.NewRequest("PATCH", url, nil)
		req.Header.Set("Authorization", "Bearer "+anotherAccessToken)

		resp, err := server.Client().Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	})

	t.Run("Deleted user", func(t *testing.T) {
		deletedTime := time.Now().UTC()
		anotherUserEntity := model2.UserEntity{
			ID:        uuid.New(),
			Username:  "user",
			Email:     "user@example.com",
			Role:      model2.RoleEntity{Name: model2.RoleUser},
			IsEnabled: true,
			DeletedAt: &deletedTime,
		}
		anotherAccessToken, _, _ := security.GenerateTokens(testEnvironment.Container.Config.Security.JWT, &anotherUserEntity)
		req, _ := http.NewRequest("PATCH", url, nil)
		req.Header.Set("Authorization", "Bearer "+anotherAccessToken)

		resp, err := server.Client().Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	})

	t.Run("User not found", func(t *testing.T) {
		accessToken, _, _ := security.GenerateTokens(
			testEnvironment.Container.Config.Security.JWT,
			&model2.UserEntity{
				ID:        uuid.New(),
				IsEnabled: true,
				DeletedAt: nil,
			},
		)
		reqBody := dto.UpdateEmailInput{
			Email: "new_user_for_update_email_1@example.com",
		}
		reqBodyReader, _ := e2e.NewJSONReader(reqBody)

		req, _ := http.NewRequest("PATCH", url, reqBodyReader)
		req.Header.Set("Authorization", "Bearer "+accessToken)

		resp, err := server.Client().Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)

		var body dto2.ErrorResponse
		err = e2e.ReadResponseBody(resp, &body)
		assert.NoError(t, err)
	})

	t.Run("Email not changed", func(t *testing.T) {
		reqBody := dto.UpdateEmailInput{
			Email: "user_for_update_email_1@example.com",
		}
		reqBodyReader, _ := e2e.NewJSONReader(reqBody)

		req, _ := http.NewRequest("PATCH", url, reqBodyReader)
		req.Header.Set("Authorization", "Bearer "+accessToken)

		resp, err := server.Client().Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var body dto.AccountOutput
		err = e2e.ReadResponseBody(resp, &body)
		assert.NoError(t, err)
		assert.Equal(t, userEntity.Email, body.Email)
	})

	t.Run("Email already in use", func(t *testing.T) {
		reqBody := dto.UpdateEmailInput{
			Email: "user_for_update_email_2@example.com",
		}
		reqBodyReader, _ := e2e.NewJSONReader(reqBody)

		req, _ := http.NewRequest("PATCH", url, reqBodyReader)
		req.Header.Set("Authorization", "Bearer "+accessToken)

		resp, err := server.Client().Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusConflict, resp.StatusCode)

		var body dto2.ErrorResponse
		err = e2e.ReadResponseBody(resp, &body)
		assert.NoError(t, err)
	})

	t.Run("Success", func(t *testing.T) {
		reqBody := dto.UpdateEmailInput{
			Email: "new_user_for_update_email_1@example.com",
		}
		reqBodyReader, _ := e2e.NewJSONReader(reqBody)

		req, _ := http.NewRequest("PATCH", url, reqBodyReader)
		req.Header.Set("Authorization", "Bearer "+accessToken)

		resp, err := server.Client().Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusAccepted, resp.StatusCode)

		// Check response
		var body dto.AccountOutput
		err = e2e.ReadResponseBody(resp, &body)
		assert.NoError(t, err)
		assert.Equal(t, reqBody.Email, body.Email)

		// Check cache
		code, _, err := testEnvironment.RedisC.Exec(context.Background(), []string{"redis-cli", "get", "email-verify.new_user_for_update_email_1@example.com"})
		assert.NoError(t, err)
		assert.Equal(t, 0, code)

		// Check email
		apiPort, err := testEnvironment.SmtpC.MappedPort(context.Background(), "8025")
		assert.NoError(t, err)

		mailhogResp, err := http.Get(fmt.Sprintf("http://localhost:%s/api/v2/messages?start=0&limit=10", apiPort.Port()))
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, mailhogResp.StatusCode)

		var mailhogMessages mailhogMessagesResponse
		err = e2e.ReadResponseBody(mailhogResp, &mailhogMessages)
		assert.NoError(t, err)
		assert.True(t, mailhogMessages.Total > 0)
	})
}

func Test_AccountHandler_VerifyEmail(t *testing.T) {
	e2e.MustLoadFixtures(fixtures)

	url := server.URL + "/v0/account/email/verify"
	userEntity := &model2.UserEntity{
		ID:              uuid.MustParse("8bffbff2-1653-4aa7-8402-1d29c4a5cae1"),
		Username:        "user_for_verify_email",
		Email:           "user_for_verify_email@example.com",
		Password:        "$2a$12$ALU3HAOtZpp22.WQVFZvnO.17WcwFrxCabnVuuz3FvPzh7TsXU8Ve",
		Role:            model2.RoleEntity{Name: model2.RoleUser},
		IsEmailVerified: false,
		IsEnabled:       true,
		DeletedAt:       nil,
	}
	accessToken, _, _ := security.GenerateTokens(testEnvironment.Container.Config.Security.JWT, userEntity)

	t.Run("Unauthorized", func(t *testing.T) {
		reqBody := dto.VerifyEmailInput{
			OTP:   "123456",
			Email: "new_user_for_verify_email_1@example.com",
		}
		reqBodyReader, _ := e2e.NewJSONReader(reqBody)

		req, _ := http.NewRequest("POST", url, reqBodyReader)
		resp, err := server.Client().Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

		var body dto2.ErrorResponse
		err = e2e.ReadResponseBody(resp, &body)
		assert.NoError(t, err)
	})

	t.Run("Banned user", func(t *testing.T) {
		anotherUserEntity := model2.UserEntity{
			ID:        uuid.New(),
			Username:  "user",
			Email:     "user@example.com",
			Role:      model2.RoleEntity{Name: model2.RoleUser},
			IsEnabled: false,
			DeletedAt: nil,
		}
		anotherAccessToken, _, _ := security.GenerateTokens(testEnvironment.Container.Config.Security.JWT, &anotherUserEntity)
		req, _ := http.NewRequest("POST", url, nil)
		req.Header.Set("Authorization", "Bearer "+anotherAccessToken)

		resp, err := server.Client().Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	})

	t.Run("Deleted user", func(t *testing.T) {
		deletedTime := time.Now().UTC()
		anotherUserEntity := model2.UserEntity{
			ID:        uuid.New(),
			Username:  "user",
			Email:     "user@example.com",
			Role:      model2.RoleEntity{Name: model2.RoleUser},
			IsEnabled: true,
			DeletedAt: &deletedTime,
		}
		anotherAccessToken, _, _ := security.GenerateTokens(testEnvironment.Container.Config.Security.JWT, &anotherUserEntity)
		req, _ := http.NewRequest("POST", url, nil)
		req.Header.Set("Authorization", "Bearer "+anotherAccessToken)

		resp, err := server.Client().Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	})

	t.Run("Bad body", func(t *testing.T) {
		bodies := []*dto.VerifyEmailInput{
			nil,

			// Bad OTP
			{
				OTP:   "",
				Email: "new_user_for_verify_email_1@example.com",
			},

			// Bad email
			{
				OTP:   "123456",
				Email: "",
			},
			{
				OTP:   "123456",
				Email: "new_user_for_verify_email_1@example.com",
			},
		}

		for i, body := range bodies {
			t.Run(fmt.Sprintf("Bad body %d", i), func(t *testing.T) {
				var reqBodyReader io.Reader = nil
				if body != nil {
					reqBodyReader, _ = e2e.NewJSONReader(body)
				}

				req, _ := http.NewRequest("POST", url, reqBodyReader)
				req.Header.Set("Authorization", "Bearer "+accessToken)

				resp, err := server.Client().Do(req)
				assert.NoError(t, err)
				assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

				var body dto2.ErrorResponse
				err = e2e.ReadResponseBody(resp, &body)
				assert.NoError(t, err)
			})
		}

	})

	t.Run("Not found cache entry", func(t *testing.T) {
		reqBody := dto.VerifyEmailInput{
			OTP:   "123456",
			Email: "new_user_for_verify_email_1@example.com",
		}
		reqBodyReader, _ := e2e.NewJSONReader(reqBody)

		req, _ := http.NewRequest("POST", url, reqBodyReader)
		req.Header.Set("Authorization", "Bearer "+accessToken)

		resp, err := server.Client().Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		var body dto2.ErrorResponse
		err = e2e.ReadResponseBody(resp, &body)
		assert.NoError(t, err)
	})

	t.Run("Incorrect OTP", func(t *testing.T) {
		// Set cache entry
		cacheEntry := dto.EmailVerifyCache{
			Email: "new_user_for_verify_email_1@example.com",
			OTP:   "123456",
		}
		err := testEnvironment.Container.Cache.Manager.Set(
			context.Background(), "email-verify.new_user_for_verify_email_1@example.com", cacheEntry,
		)
		assert.NoError(t, err)

		// Send request
		reqBody := dto.VerifyEmailInput{
			OTP:   "654321",
			Email: "new_user_for_verify_email_1@example.com",
		}
		reqBodyReader, _ := e2e.NewJSONReader(reqBody)

		req, _ := http.NewRequest("POST", url, reqBodyReader)
		req.Header.Set("Authorization", "Bearer "+accessToken)

		resp, err := server.Client().Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		var body dto2.ErrorResponse
		err = e2e.ReadResponseBody(resp, &body)
		assert.NoError(t, err)
	})

	t.Run("User not found", func(t *testing.T) {
		// Set cache entry
		cacheEntry := dto.EmailVerifyCache{
			Email: "new_user_for_verify_email_1@example.com",
			OTP:   "123456",
		}
		err := testEnvironment.Container.Cache.Manager.Set(
			context.Background(), "email-verify.new_user_for_verify_email_1@example.com", cacheEntry,
		)
		assert.NoError(t, err)

		// Another access token
		anotherAccessToken, _, _ := security.GenerateTokens(
			testEnvironment.Container.Config.Security.JWT,
			&model2.UserEntity{
				ID:        uuid.New(),
				IsEnabled: true,
				DeletedAt: nil,
			},
		)
		// Send request
		reqBody := dto.VerifyEmailInput{
			OTP:   "123456",
			Email: "new_user_for_verify_email_1@example.com",
		}
		reqBodyReader, _ := e2e.NewJSONReader(reqBody)

		req, _ := http.NewRequest("POST", url, reqBodyReader)
		req.Header.Set("Authorization", "Bearer "+anotherAccessToken)

		resp, err := server.Client().Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)

		var body dto2.ErrorResponse
		err = e2e.ReadResponseBody(resp, &body)
		assert.NoError(t, err)
	})

	t.Run("Incorrect email", func(t *testing.T) {
		// Set cache entry
		cacheEntry := dto.EmailVerifyCache{
			Email: "new_user_for_verify_email_2@example.com",
			OTP:   "123456",
		}
		err := testEnvironment.Container.Cache.Manager.Set(
			context.Background(), "email-verify.new_user_for_verify_email_2@example.com", cacheEntry,
		)
		assert.NoError(t, err)

		// Send request
		reqBody := dto.VerifyEmailInput{
			OTP:   "123456",
			Email: "new_user_for_verify_email_2@example.com",
		}
		reqBodyReader, _ := e2e.NewJSONReader(reqBody)

		req, _ := http.NewRequest("POST", url, reqBodyReader)
		req.Header.Set("Authorization", "Bearer "+accessToken)

		resp, err := server.Client().Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		var body dto2.ErrorResponse
		err = e2e.ReadResponseBody(resp, &body)
		assert.NoError(t, err)
	})

	t.Run("Success", func(t *testing.T) {
		// Set cache entry
		cacheEntry := dto.EmailVerifyCache{
			Email: "user_for_verify_email@example.com",
			OTP:   "123456",
		}
		err := testEnvironment.Container.Cache.Manager.Set(
			context.Background(), "email-verify.user_for_verify_email@example.com", cacheEntry,
		)
		assert.NoError(t, err)

		// Send request
		reqBody := dto.VerifyEmailInput{
			OTP:   "123456",
			Email: "user_for_verify_email@example.com",
		}
		reqBodyReader, _ := e2e.NewJSONReader(reqBody)

		req, _ := http.NewRequest("POST", url, reqBodyReader)
		req.Header.Set("Authorization", "Bearer "+accessToken)

		resp, err := server.Client().Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}

func Test_AccountHandler_SetPassword(t *testing.T) {
	e2e.MustLoadFixtures(fixtures)

	url := server.URL + "/v0/account/password"
	userEntity := model2.UserEntity{
		ID:             uuid.MustParse("e5750b61-7eda-41dd-b4ab-0097d4dbc92e"),
		Username:       "user_for_set_password_1",
		Email:          "user_for_set_password_1@example.com",
		Role:           model2.RoleEntity{Name: model2.RoleUser},
		IsPasswordTemp: true,
		IsEnabled:      true,
		DeletedAt:      nil,
	}

	accessToken, _, _ := security.GenerateTokens(testEnvironment.Container.Config.Security.JWT, &userEntity)

	t.Run("Unauthorized", func(t *testing.T) {
		// Send request
		reqBody := dto.SetPasswordInput{
			Password: "7676393c-0d28-4f68-807d-12aa6b88c039",
		}
		reqBodyReader, _ := e2e.NewJSONReader(reqBody)

		req, _ := http.NewRequest("POST", url, reqBodyReader)

		resp, err := server.Client().Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

		var body dto2.ErrorResponse
		err = e2e.ReadResponseBody(resp, &body)
		assert.NoError(t, err)
	})

	t.Run("Banned user", func(t *testing.T) {
		anotherUserEntity := model2.UserEntity{
			ID:        uuid.New(),
			Username:  "user",
			Email:     "user@example.com",
			Role:      model2.RoleEntity{Name: model2.RoleUser},
			IsEnabled: false,
			DeletedAt: nil,
		}
		anotherAccessToken, _, _ := security.GenerateTokens(testEnvironment.Container.Config.Security.JWT, &anotherUserEntity)
		req, _ := http.NewRequest("POST", url, nil)
		req.Header.Set("Authorization", "Bearer "+anotherAccessToken)

		resp, err := server.Client().Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	})

	t.Run("Deleted user", func(t *testing.T) {
		deletedTime := time.Now().UTC()
		anotherUserEntity := model2.UserEntity{
			ID:        uuid.New(),
			Username:  "user",
			Email:     "user@example.com",
			Role:      model2.RoleEntity{Name: model2.RoleUser},
			IsEnabled: true,
			DeletedAt: &deletedTime,
		}
		anotherAccessToken, _, _ := security.GenerateTokens(testEnvironment.Container.Config.Security.JWT, &anotherUserEntity)
		req, _ := http.NewRequest("POST", url, nil)
		req.Header.Set("Authorization", "Bearer "+anotherAccessToken)

		resp, err := server.Client().Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	})

	t.Run("Bad body", func(t *testing.T) {
		bodies := []*dto.SetPasswordInput{
			nil,

			// Bad password
			{
				Password: "",
			},
			{
				Password: "weak",
			},
		}

		for i, body := range bodies {
			t.Run(fmt.Sprintf("Bad body %d", i), func(t *testing.T) {
				var reqBodyReader io.Reader = nil
				if body != nil {
					reqBodyReader, _ = e2e.NewJSONReader(body)
				}

				// Send request
				req, _ := http.NewRequest("POST", url, reqBodyReader)
				req.Header.Set("Authorization", "Bearer "+accessToken)

				resp, err := server.Client().Do(req)
				assert.NoError(t, err)
				assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

				var body dto2.ErrorResponse
				err = e2e.ReadResponseBody(resp, &body)
				assert.NoError(t, err)
			})
		}
	})

	t.Run("User not found", func(t *testing.T) {
		// New access token for another user
		anotherAccessToken, _, _ := security.GenerateTokens(
			testEnvironment.Container.Config.Security.JWT,
			&model2.UserEntity{
				ID:        uuid.New(),
				IsEnabled: true,
				DeletedAt: nil,
			},
		)

		// Send request
		reqBody := dto.SetPasswordInput{
			Password: "7676393c-0d28-4f68-807d-12aa6b88c039",
		}
		reqBodyReader, _ := e2e.NewJSONReader(reqBody)

		req, _ := http.NewRequest("POST", url, reqBodyReader)
		req.Header.Set("Authorization", "Bearer "+anotherAccessToken)

		resp, err := server.Client().Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)

		var body dto2.ErrorResponse
		err = e2e.ReadResponseBody(resp, &body)
		assert.NoError(t, err)
	})

	t.Run("Password is not temporary", func(t *testing.T) {
		userEntity2 := model2.UserEntity{
			ID:             uuid.MustParse("33dcd273-f31f-4dde-a37e-9cd2b3e16fcc"),
			Username:       "user_for_set_password_2",
			Email:          "user_for_set_password_2@example.com",
			Role:           model2.RoleEntity{Name: model2.RoleUser},
			IsPasswordTemp: true,
			IsEnabled:      true,
			DeletedAt:      nil,
		}
		anotherAccessToken, _, _ := security.GenerateTokens(testEnvironment.Container.Config.Security.JWT, &userEntity2)

		// Send request
		reqBody := dto.SetPasswordInput{
			Password: "7676393c-0d28-4f68-807d-12aa6b88c039",
		}
		reqBodyReader, _ := e2e.NewJSONReader(reqBody)

		req, _ := http.NewRequest("POST", url, reqBodyReader)
		req.Header.Set("Authorization", "Bearer "+anotherAccessToken)

		resp, err := server.Client().Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusConflict, resp.StatusCode)

		var body dto2.ErrorResponse
		err = e2e.ReadResponseBody(resp, &body)
		assert.NoError(t, err)
	})

	t.Run("Success", func(t *testing.T) {
		// Send request
		reqBody := dto.SetPasswordInput{
			Password: "7676393c-0d28-4f68-807d-12aa6b88c039",
		}
		reqBodyReader, _ := e2e.NewJSONReader(reqBody)

		req, _ := http.NewRequest("POST", url, reqBodyReader)
		req.Header.Set("Authorization", "Bearer "+accessToken)

		resp, err := server.Client().Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}

func Test_AccountHandler_UpdatePassword(t *testing.T) {
	e2e.MustLoadFixtures(fixtures)

	url := server.URL + "/v0/account/password"
	userEntity := model2.UserEntity{
		ID:        uuid.MustParse("b03643e2-263d-406d-ac8f-02e51fe8927e"),
		Username:  "user_for_update_password_1",
		Email:     "user_for_update_password_1@example.com",
		Password:  "$2a$12$.s1AcMGZNlbGfWCUebmcGeHpV0bMUWorZ4Zmx/YwY5RVfp.OZbpDG",
		Role:      model2.RoleEntity{Name: model2.RoleUser},
		IsEnabled: true,
		DeletedAt: nil,
	}

	accessToken, _, _ := security.GenerateTokens(testEnvironment.Container.Config.Security.JWT, &userEntity)

	t.Run("Unauthorized", func(t *testing.T) {
		// Send request
		reqBody := dto.UpdatePasswordInput{
			OldPassword: "user_for_update_password_1",
			NewPassword: "7676393c-0d28-4f68-807d-12aa6b88c039",
		}
		reqBodyReader, _ := e2e.NewJSONReader(reqBody)

		req, _ := http.NewRequest("PATCH", url, reqBodyReader)

		resp, err := server.Client().Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

		var body dto2.ErrorResponse
		err = e2e.ReadResponseBody(resp, &body)
		assert.NoError(t, err)
	})

	t.Run("Banned user", func(t *testing.T) {
		anotherUserEntity := model2.UserEntity{
			ID:        uuid.New(),
			Username:  "user",
			Email:     "user@example.com",
			Role:      model2.RoleEntity{Name: model2.RoleUser},
			IsEnabled: false,
			DeletedAt: nil,
		}
		anotherAccessToken, _, _ := security.GenerateTokens(testEnvironment.Container.Config.Security.JWT, &anotherUserEntity)
		req, _ := http.NewRequest("PATCH", url, nil)
		req.Header.Set("Authorization", "Bearer "+anotherAccessToken)

		resp, err := server.Client().Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	})

	t.Run("Deleted user", func(t *testing.T) {
		deletedTime := time.Now().UTC()
		anotherUserEntity := model2.UserEntity{
			ID:        uuid.New(),
			Username:  "user",
			Email:     "user@example.com",
			Role:      model2.RoleEntity{Name: model2.RoleUser},
			IsEnabled: true,
			DeletedAt: &deletedTime,
		}
		anotherAccessToken, _, _ := security.GenerateTokens(testEnvironment.Container.Config.Security.JWT, &anotherUserEntity)
		req, _ := http.NewRequest("PATCH", url, nil)
		req.Header.Set("Authorization", "Bearer "+anotherAccessToken)

		resp, err := server.Client().Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	})

	t.Run("Bad body", func(t *testing.T) {
		bodies := []*dto.UpdatePasswordInput{
			nil,

			// Bad old password
			{
				OldPassword: "",
				NewPassword: "7676393c-0d28-4f68-807d-12aa6b88c039",
			},

			// Bad new password
			{
				OldPassword: "user_for_update_password_1",
				NewPassword: "",
			},
			{
				OldPassword: "user_for_update_password_1",
				NewPassword: "weak",
			},
		}

		for i, body := range bodies {
			t.Run(fmt.Sprintf("Bad body %d", i), func(t *testing.T) {
				var reqBodyReader io.Reader = nil
				if body != nil {
					reqBodyReader, _ = e2e.NewJSONReader(body)
				}

				// Send request
				req, _ := http.NewRequest("PATCH", url, reqBodyReader)
				req.Header.Set("Authorization", "Bearer "+accessToken)

				resp, err := server.Client().Do(req)
				assert.NoError(t, err)
				assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

				var body dto2.ErrorResponse
				err = e2e.ReadResponseBody(resp, &body)
				assert.NoError(t, err)
			})
		}
	})

	t.Run("User not found", func(t *testing.T) {
		// New access token for another user
		anotherAccessToken, _, _ := security.GenerateTokens(
			testEnvironment.Container.Config.Security.JWT,
			&model2.UserEntity{
				ID:        uuid.New(),
				IsEnabled: true,
				DeletedAt: nil,
			},
		)

		// Send request
		reqBody := dto.UpdatePasswordInput{
			OldPassword: "user_for_update_password_1",
			NewPassword: "7676393c-0d28-4f68-807d-12aa6b88c039",
		}
		reqBodyReader, _ := e2e.NewJSONReader(reqBody)

		req, _ := http.NewRequest("PATCH", url, reqBodyReader)
		req.Header.Set("Authorization", "Bearer "+anotherAccessToken)

		resp, err := server.Client().Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)

		var body dto2.ErrorResponse
		err = e2e.ReadResponseBody(resp, &body)
		assert.NoError(t, err)
	})

	t.Run("Password is incorrect", func(t *testing.T) {
		// Send request
		reqBody := dto.UpdatePasswordInput{
			OldPassword: "incorrect_user_for_update_password_1",
			NewPassword: "7676393c-0d28-4f68-807d-12aa6b88c039",
		}
		reqBodyReader, _ := e2e.NewJSONReader(reqBody)

		req, _ := http.NewRequest("PATCH", url, reqBodyReader)
		req.Header.Set("Authorization", "Bearer "+accessToken)

		resp, err := server.Client().Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		var body dto2.ErrorResponse
		err = e2e.ReadResponseBody(resp, &body)
		assert.NoError(t, err)
	})

	t.Run("Success", func(t *testing.T) {
		// Send request
		reqBody := dto.UpdatePasswordInput{
			OldPassword: "user_for_update_password_1",
			NewPassword: "7676393c-0d28-4f68-807d-12aa6b88c039",
		}
		reqBodyReader, _ := e2e.NewJSONReader(reqBody)

		req, _ := http.NewRequest("PATCH", url, reqBodyReader)
		req.Header.Set("Authorization", "Bearer "+accessToken)

		resp, err := server.Client().Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}

func Test_AccountHandler_RestoreAccount(t *testing.T) {
	e2e.MustLoadFixtures(fixtures)

	url := server.URL + "/v0/account/restore"

	userEntity := model2.UserEntity{
		ID:        uuid.MustParse("fb22c374-e2c9-4f47-aa17-a9853bd29a58"),
		Username:  "user_for_restore_1",
		Email:     "user_for_restore_1@example.com",
		Role:      model2.RoleEntity{Name: model2.RoleUser},
		IsEnabled: true,
		DeletedAt: nil,
	}
	accessToken, _, _ := security.GenerateTokens(testEnvironment.Container.Config.Security.JWT, &userEntity)

	t.Run("Unauthorized", func(t *testing.T) {
		// Send request
		req, _ := http.NewRequest("GET", url, nil)

		resp, err := server.Client().Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("Banned user", func(t *testing.T) {
		anotherUserEntity := model2.UserEntity{
			ID:        uuid.New(),
			Username:  "user",
			Email:     "user@example.com",
			Role:      model2.RoleEntity{Name: model2.RoleUser},
			IsEnabled: false,
			DeletedAt: nil,
		}
		anotherAccessToken, _, _ := security.GenerateTokens(testEnvironment.Container.Config.Security.JWT, &anotherUserEntity)
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("Authorization", "Bearer "+anotherAccessToken)

		resp, err := server.Client().Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	})

	t.Run("User not found", func(t *testing.T) {
		// New access token for another user
		anotherAccessToken, _, _ := security.GenerateTokens(
			testEnvironment.Container.Config.Security.JWT,
			&model2.UserEntity{
				ID:        uuid.New(),
				IsEnabled: true,
				DeletedAt: nil,
			},
		)

		// Send request
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("Authorization", "Bearer "+anotherAccessToken)

		resp, err := server.Client().Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("User not deleted", func(t *testing.T) {
		// New access token for another user
		anotherUserEntity := model2.UserEntity{
			ID:        uuid.MustParse("8e2e4465-87d9-4156-84c8-49fa0afe2809"),
			Username:  "user_for_restore_2",
			Email:     "user_for_restore_2@example.com",
			Role:      model2.RoleEntity{Name: model2.RoleUser},
			IsEnabled: true,
			DeletedAt: nil,
		}
		anotherAccessToken, _, _ := security.GenerateTokens(testEnvironment.Container.Config.Security.JWT, &anotherUserEntity)

		// Send request
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("Authorization", "Bearer "+anotherAccessToken)

		resp, err := server.Client().Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusConflict, resp.StatusCode)
	})

	t.Run("Success", func(t *testing.T) {
		// Send request
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("Authorization", "Bearer "+accessToken)

		resp, err := server.Client().Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}

func Test_AccountHandler_DeleteAccount(t *testing.T) {
	e2e.MustLoadFixtures(fixtures)

	url := server.URL + "/v0/account"

	userEntity := model2.UserEntity{
		ID:        uuid.MustParse("aa693ac7-5f36-47d4-b7e0-582c5eab3d0f"),
		Username:  "user_for_delete_1",
		Email:     "user_for_delete_1@example.com",
		Role:      model2.RoleEntity{Name: model2.RoleUser},
		IsEnabled: true,
		DeletedAt: nil,
	}
	accessToken, _, _ := security.GenerateTokens(testEnvironment.Container.Config.Security.JWT, &userEntity)

	t.Run("Unauthorized", func(t *testing.T) {
		// Send request
		req, _ := http.NewRequest("DELETE", url, nil)

		resp, err := server.Client().Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("Banned user", func(t *testing.T) {
		anotherUserEntity := model2.UserEntity{
			ID:        uuid.New(),
			Username:  "user",
			Email:     "user@example.com",
			Role:      model2.RoleEntity{Name: model2.RoleUser},
			IsEnabled: false,
			DeletedAt: nil,
		}
		anotherAccessToken, _, _ := security.GenerateTokens(testEnvironment.Container.Config.Security.JWT, &anotherUserEntity)
		req, _ := http.NewRequest("DELETE", url, nil)
		req.Header.Set("Authorization", "Bearer "+anotherAccessToken)

		resp, err := server.Client().Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	})

	t.Run("Deleted user", func(t *testing.T) {
		deletedTime := time.Now().UTC()
		anotherUserEntity := model2.UserEntity{
			ID:        uuid.New(),
			Username:  "user",
			Email:     "user@example.com",
			Role:      model2.RoleEntity{Name: model2.RoleUser},
			IsEnabled: true,
			DeletedAt: &deletedTime,
		}
		anotherAccessToken, _, _ := security.GenerateTokens(testEnvironment.Container.Config.Security.JWT, &anotherUserEntity)
		req, _ := http.NewRequest("DELETE", url, nil)
		req.Header.Set("Authorization", "Bearer "+anotherAccessToken)

		resp, err := server.Client().Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	})

	t.Run("User not found", func(t *testing.T) {
		// New access token for another user
		anotherAccessToken, _, _ := security.GenerateTokens(
			testEnvironment.Container.Config.Security.JWT,
			&model2.UserEntity{
				ID:        uuid.New(),
				IsEnabled: true,
				DeletedAt: nil,
			},
		)

		// Send request
		req, _ := http.NewRequest("DELETE", url, nil)
		req.Header.Set("Authorization", "Bearer "+anotherAccessToken)

		resp, err := server.Client().Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("User already deleted", func(t *testing.T) {
		// New access token for another user
		anotherUserEntity := model2.UserEntity{
			ID:        uuid.MustParse("585d8642-e4bf-40aa-9084-89890fc1639f"),
			Username:  "user_for_delete_2",
			Email:     "user_for_delete_2@example.com",
			Role:      model2.RoleEntity{Name: model2.RoleUser},
			IsEnabled: true,
			DeletedAt: nil,
		}
		anotherAccessToken, _, _ := security.GenerateTokens(testEnvironment.Container.Config.Security.JWT, &anotherUserEntity)

		// Send request
		req, _ := http.NewRequest("DELETE", url, nil)
		req.Header.Set("Authorization", "Bearer "+anotherAccessToken)

		resp, err := server.Client().Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusConflict, resp.StatusCode)
	})

	t.Run("Success", func(t *testing.T) {
		// Send request
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("Authorization", "Bearer "+accessToken)

		resp, err := server.Client().Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}
