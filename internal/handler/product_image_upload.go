package handler

import (
	"fmt"
	"log"
	"mime"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

func UploadProductImageHandler(c *fiber.Ctx) error {
	// Handle file upload
	file, err := c.FormFile("image")
	if err != nil {
		return  c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to get uploaded file",
		})
	}

	//validate image
	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowedExts := map[string]bool{
		".jpg":  true,	
		".jpeg": true,
		".png":  true,
		".webp": true,
	}

	if !allowedExts[ext] {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid file type. Only JPG, JPEG, PNG, and WEBP are allowed",
		})
	}

	// validate content type
	contentType := file.Header.Get("Content-Type")
	allowedContentTypes := map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
		"image/webp": true,
		"image/jpg":  true,
	}
	if !allowedContentTypes[contentType] {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid content type. Only JPG, JPEG, PNG, and WEBP are allowed",
		})
	}




	// Format image name
	timestamp := time.Now().UnixNano()
	fileName := fmt.Sprintf("product_%d%s", timestamp,filepath.Ext(file.Filename))
	uploadPath := "./storage/product_images/" + fileName

	err = c.SaveFile(file,uploadPath)

	if err != nil {
		fmt.Println("Error saving file:", err)
		return  c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to save file",
		})
	}
	return c.JSON(fiber.Map{
		"success" : true,
		"message" : "Upload success",
		"file_name" : fileName,
	})
}

func GetProductImageHandler(c *fiber.Ctx) error {
	filenameparam := c.Params("filename")
	filePath := filepath.Join("storage", "product_images", filenameparam)
	if _, err := os.Stat(filePath); err !=  nil{
		if os.IsNotExist(err){
			return c.Status(http.StatusNotFound).SendString("File not found")
		}

		log.Println("Error accessing file:", err)
		return c.Status(http.StatusInternalServerError).SendString("Internal Server Error")
	}
	

	file, err := os.Open(filePath)
	if err != nil {
		log.Println("Error opening file:", err)
		return c.Status(http.StatusInternalServerError).SendString("Internal Server Error")
	}

	ext := path.Ext(filePath)
	mimeType := mime.TypeByExtension(ext)
	c.Set("Content-Type",mimeType)
	return c.SendStream(file)
}