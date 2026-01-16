package server

import (
	"log"
	"net/http"
)

// Swagger UI is served from a CDN and points to /swagger/openapi.yaml.
const swaggerUIHTML = `<!doctype html>
<html lang="en">
  <head>
    <meta charset="utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    <title>URL Shortener API - Swagger UI</title>
    <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css" />
    <style>
      html, body { margin: 0; padding: 0; }
    </style>
  </head>
  <body>
    <div id="swagger-ui"></div>
    <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
    <script>
      window.onload = () => {
        window.ui = SwaggerUIBundle({
          url: '/swagger/openapi.yaml',
          dom_id: '#swagger-ui',
          deepLinking: true,
          presets: [SwaggerUIBundle.presets.apis],
          layout: 'BaseLayout'
        });
      };
    </script>
  </body>
</html>`

// SwaggerHandler serves:
//   - GET /swagger/              : Swagger UI HTML
//   - GET /swagger/openapi.yaml  : OpenAPI YAML contract
func SwaggerHandler() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /swagger", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/swagger/", http.StatusPermanentRedirect)
	})

	mux.HandleFunc("GET /swagger/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if _, err := w.Write([]byte(swaggerUIHTML)); err != nil {
			// Too late to change HTTP response – just log it.
			log.Printf("swagger ui write failed: %v", err)
		}
	})

	mux.HandleFunc("GET /swagger/openapi.yaml", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "openapi/openapi.yaml")
	})

	return mux
}
