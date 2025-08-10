module github.com/aman/internal/services/user

go 1.24.6

require (
	github.com/aman/internal/services/user/pb v0.0.0-00010101000000-000000000000
	google.golang.org/grpc v1.74.2
)

replace github.com/aman/internal/services/user/pb => ./pb

require (
	golang.org/x/net v0.43.0 // indirect
	golang.org/x/sys v0.35.0 // indirect
	golang.org/x/text v0.28.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250804133106-a7a43d27e69b // indirect
	google.golang.org/protobuf v1.36.7 // indirect
)
