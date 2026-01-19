package main

import (
	"embed"
	"fmt"
	"log"
	"time"

	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
)

//go:embed index.html
var content embed.FS

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

func Index(ctx *fasthttp.RequestCtx) {
	data, err := content.ReadFile("index.html")
	if err != nil {
		ctx.SetStatusCode(500)
		ctx.SetContentType("text/plain; charset=utf-8")
		ctx.WriteString("error loading index.html: " + err.Error())
		return
	}

	ctx.SetContentType("text/html; charset=utf-8")
	ctx.Write(data)
}

func main() {
	var question string
	fmt.Print("   Wana run the Server? (yes/no)")
	fmt.Print("\n-> ")
	fmt.Scan(&question)
	if question == "yes" {
		fmt.Println("Running on http://localhost:8080")

		r := router.New()
		r.GET("/", Index)
		log.Fatal(fasthttp.ListenAndServe(":8080", r.Handler))

		server := &fasthttp.Server{
			Handler:            r.Handler,
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
