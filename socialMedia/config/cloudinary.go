package config

import (
	"log"
	"os"

	"github.com/cloudinary/cloudinary-go"
)

var (
	Cld *cloudinary.Cloudinary
)

func ConnectToCloudinary() {
	var err error
	cloudName := os.Getenv("CLOUDINARY_CLOUD_NAME")
	apiKey := os.Getenv("CLOUDINARY_API_KEY")
	apiSecret := os.Getenv("CLOUDINARY_API_SECRET")

	Cld, err = cloudinary.NewFromParams(cloudName, apiKey, apiSecret)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Successfully connected to Cloudinary")
}
