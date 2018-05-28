package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	ctf "github.com/ohookins/contentful-go"
)

func getContentTypeFor(url string) string {
	resp, err := http.Head(url)
	if err != nil {
		fmt.Println(err.Error())
		return ""
	}

	return resp.Header.Get("Content-Type")
}

func getUpdatedAsset(cma *ctf.Contentful, space, assetID string) (*ctf.Asset, error) {
	for {
		asset, err := cma.Assets.Get(space, assetID)
		if err != nil {
			return nil, err
		}

		// URL property of the localized asset will be replaced with the
		// ctfassets.net URL of the processed asset when it is completed.
		if asset.Fields.File["en-US"].URL != "" {
			return asset, nil
		}
		time.Sleep(time.Second)
	}
}

func createAttachments(cma *ctf.Contentful, items []item, space string) error {
	// filter out all actual attachments
	attachments := []item{}
	for _, item := range items {
		if strings.Contains(item.Link, "attachment_id") {
			attachments = append(attachments, item)
		}
	}

	fmt.Printf("creating %d new assets\n", len(attachments))
	for _, item := range attachments {
		// Get the actual filename of the attachment
		parts := strings.Split(item.Guid, "/")
		filename := parts[len(parts)-1]

		// Determine the content type of the currently live attachment in WP
		contentType := getContentTypeFor(item.Guid)

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
						UploadURL:   item.Guid,
						ContentType: contentType,
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

		// Block until the asset is processed
		asset, err := getUpdatedAsset(cma, space, asset.Sys.ID)
		if err != nil {
			return err
		}

		// Add the new URL to the map for later replacement in post links.
		replacementURLs[strings.Replace(item.Guid, "http:", "https:", -1)] = asset.Fields.File["en-US"].URL
		replacementURLs[item.Guid] = asset.Fields.File["en-US"].URL

		if err := cma.Assets.Publish(space, asset); err != nil {
			return err
		}
	}

	return nil
}
