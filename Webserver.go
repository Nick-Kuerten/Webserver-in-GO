package main

import (
	"embed"
	"fmt"
	"html/template"
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
	tmpl, err := template.ParseFS(content, "index.html")
	if err != nil {
		fallbackTmpl, tmplErr := template.New("error").Parse(`There was an unexpected error loading index.html: {{.}}`)
		if tmplErr != nil {
			ctx.SetStatusCode(500)
			ctx.SetContentType("text/plain; charset=utf-8")
			ctx.WriteString("error creating error-template: " + tmplErr.Error())
			return
		}

		ctx.SetStatusCode(500)
		ctx.SetContentType("text/html; charset=utf-8")

		if execErr := fallbackTmpl.Execute(ctx.Response.BodyWriter(), err.Error()); execErr != nil {
			ctx.SetStatusCode(500)
			ctx.SetContentType("text/plain; charset=utf-8")
			ctx.WriteString("error executing error-template: " + execErr.Error())
		}
		return
	}

	ctx.SetContentType("text/html; charset=utf-8")

	data := map[string]any{
		"Title":  "Webserver",
		"Now":    time.Now().Format(time.RFC3339),
		"Path":   ctx.Path(),
		"Method": ctx.Method(),
		"Host":   ctx.Host(),
	}

	if execErr := tmpl.ExecuteTemplate(ctx.Response.BodyWriter(), "index.html", data); execErr != nil {
		ctx.SetStatusCode(500)
		ctx.SetContentType("text/plain; charset=utf-8")
		ctx.WriteString("error executing template: " + execErr.Error())
		return
	}
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
