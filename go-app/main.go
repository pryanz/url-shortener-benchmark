package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type URLShortenRequest struct{
	URL string `json:"url" binding:"required,url"`
}

func shortenURL(c *gin.Context){
	var req URLShortenRequest
	
	if err := c.ShouldBindJSON(&req); err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid URL format"})
		return
	}

	id , err := RedisClient.Incr(Ctx,"global:next_id").Result()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	shortened_code := Encode(int(id)) 
	err = RedisClient.Set(Ctx, fmt.Sprintf("url:%d",id), req.URL, 0).Err()
	if err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save URL"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"short_code": shortened_code})
}

func redirectURL(c *gin.Context){
	shortened_code := c.Param("short_code")
	id := Decode(shortened_code)
	url , err := RedisClient.Get(Ctx, fmt.Sprintf("url:%d",id)).Result()
	if err == redis.Nil {
		c.JSON(http.StatusNotFound, gin.H{"error":"url does not exist"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	c.Redirect(http.StatusFound, url) 
}

func main(){
	InitRedis()
	
	router := gin.Default()

	router.GET("/", func(c * gin.Context){
		c.String(200,"Go baseline active")
	})

	router.POST("/shorten", shortenURL)
	router.GET("/:short_code", redirectURL)
	router.HEAD("/:short_code", redirectURL)

	router.Run(":8081")
}