package utils

import (
	"context"
	"fmt"

	"github.com/cloudinary/cloudinary-go"
	"github.com/cloudinary/cloudinary-go/api/admin"
)

func DeleteFromCloudinary(cld *cloudinary.Cloudinary, publicID string) error {
	var ctx = context.Background()
	_, err := cld.Admin.DeleteAssets(ctx, admin.DeleteAssetsParams{
		PublicIDs: []string{publicID},
	})
	if err != nil {
		return fmt.Errorf("cannot delete image: %w", err)
	}
	return nil
}
