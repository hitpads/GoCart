package main

import (
	"bytes"
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	inventoryServiceURL = "http://localhost:8001"
	orderServiceURL     = "http://localhost:8002"
)

// redir incoming request to target port
func proxyRequest(c *gin.Context, targetURL string) {

	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read request body"})
		return
	}
	log.Printf("Incoming request body: %s", string(bodyBytes))

	// RECREATE request body for redirecting
	c.Request.Body = io.NopCloser(io.MultiReader(bytes.NewReader(bodyBytes)))

	req, err := http.NewRequest(c.Request.Method, targetURL+c.Request.RequestURI, c.Request.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create request"})
		return
	}
	req.Header = c.Request.Header

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "service unavailable"})
		return
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read response"})
		return
	}
	c.Data(resp.StatusCode, resp.Header.Get("Content-Type"), respBody)
}

func main() {
	r := gin.Default()

	r.Static("/static", "./static")

	r.GET("/", func(ctx *gin.Context) {
		ctx.File("./static/index.html")
	})

	r.Any("/products/*any", func(c *gin.Context) {
		proxyRequest(c, inventoryServiceURL)
	})

	r.Any("/orders/*any", func(c *gin.Context) {
		proxyRequest(c, orderServiceURL)
	})

	if err := r.Run(":8000"); err != nil {
		log.Fatal(err)
	}
}
