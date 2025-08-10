module github.com/aman/internal/libraries/auth

go 1.24.6

require (
	github.com/aman/internal/logging v0.0.0-00010101000000-000000000000
	github.com/golang-jwt/jwt/v5 v5.3.0
)

require (
	go.uber.org/multierr v1.10.0 // indirect
	go.uber.org/zap v1.27.0 // indirect
)

replace github.com/aman/internal/logging => ../../logging
