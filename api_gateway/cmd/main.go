package main

import (
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	inventoryServiceURL = "http://localhost:8001"
	orderServiceURL     = "http://localhost:8002"
)

// forward the incoming HTTP request to the target URL
func proxyRequest(c *gin.Context, targetURL string) {
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

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read response"})
		return
	}
	c.Data(resp.StatusCode, resp.Header.Get("Content-Type"), body)
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
