package utils

import (
	"context"
	"fmt"

	"github.com/cloudinary/cloudinary-go"
	"github.com/cloudinary/cloudinary-go/api/uploader"
)

func UploadToCloudinary(cld *cloudinary.Cloudinary, filePath string) (*uploader.UploadResult, error) {
	var ctx = context.Background()
	resp, err := cld.Upload.Upload(ctx, filePath, uploader.UploadParams{})
	if err != nil {
		return nil, fmt.Errorf("cannot upload image: %w", err)
	}

	return resp, nil
}
