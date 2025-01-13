templ generate

GOOS=js GOARCH=wasm go build -o ./public/lib.wasm ./front/web.go

tailwindcss -i ./global.css -o ./public/output.css

go run ./server/server.go