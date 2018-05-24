package main

import (
	"fmt"
	"strings"

	ctf "github.com/ohookins/contentful-go"
)

func createAttachments(cma *ctf.Contentful, items []item, space string) error {
	fmt.Printf("creating %d new assets\n", len(items))
	for _, item := range items {
		// Skip anything which is not an attachment
		if !strings.Contains(item.Link, "attachment_id") {
			continue
		}

		// Get the actual filename of the attachment
		parts := strings.Split(item.Link, "/")
		filename := parts[len(parts)-1]

		asset := &ctf.Asset{
			Sys: &ctf.Sys{
				ID: item.PostName,
			},
			Fields: &ctf.FileFields{
				Title:       map[string]string{"en-US": item.Title},
				Description: map[string]string{"en-US": fmt.Sprintf("Linked from post ID %d", item.PostID)},
				File: map[string]*ctf.File{
					"en-US": &ctf.File{
						Name:        filename,
						URL:         item.Guid,
						ContentType: "image/jpeg",
					},
				},
			},
		}

		fmt.Printf("creating new asset with ID %s\n", item.PostName)
		if err := cma.Assets.Upsert(space, asset); err != nil {
			return err
		}
		if err := cma.Assets.Process(space, asset); err != nil {
			return err
		}
		if err := cma.Assets.Publish(space, asset); err != nil {
			return err
		}
	}

	return nil
}
