package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/valyala/fasthttp"
)

func notFoundHandler(ctx *fasthttp.RequestCtx) {
	ctx.SetStatusCode(404)
	ctx.SetContentType("text/html")
	html := `
	<!DOCTYPE html>
	<html>
	<head>
		<title>404 - Page Not Found</title>
	</head>
	<body>
		<h1>404 - Page Not Found</h1>
		<p>The page you are looking for does not exist.</p>
		<p>Path: ` + string(ctx.Path()) + `</p>
		<a href="/">Back to Homepage</a>
	</body>
	</html>
	`
	ctx.WriteString(html)
}

func router(ctx *fasthttp.RequestCtx) {
	path := string(ctx.Path())

	if path == "/test" {
		ctx.SetContentType("text/html")
		ctx.WriteString("<h1>Test Seite</h1>")
	} else if path == "/" {
		fs := &fasthttp.FS{
			Root:       ".",
			IndexNames: []string{"index.html"},
			Compress:   true,
		}
		fs.NewRequestHandler()(ctx)
	} else {
		filePath := "." + path
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			notFoundHandler(ctx)
		} else {
			fs := &fasthttp.FS{
				Root:       ".",
				IndexNames: []string{"index.html"},
				Compress:   true,
			}
			fs.NewRequestHandler()(ctx)
		}
	}
}

func main() {
	var question string
	fmt.Print("   Wana run the Server? (yes/no)")
	fmt.Print("\n-> ")
	fmt.Scan(&question)
	if question == "yes" {
		fmt.Println("Running on http://localhost:8080")

		server := &fasthttp.Server{
			Handler:            router,
			ReadTimeout:        5 * time.Minute,
			WriteTimeout:       10 * time.Second,
			MaxRequestBodySize: 2 * 1024 * 1024,
		}

		err := server.ListenAndServe(":8080")
		if err != nil {
			log.Fatal(err)
		}
	} else {
		fmt.Println("Exiting...")
	}
}
