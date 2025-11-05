# Start the server with .env loaded (PowerShell)
if (Test-Path .env) {
    Write-Output "Using .env"
} else {
    Write-Output "No .env file found. Create it from .env.example"
}

# Run the Go server
go run ./cmd/main.go
