package main

import (
	"fmt"
	"strings"

	ctf "github.com/ohookins/contentful-go"
)

func createAttachments(cma *ctf.Contentful, items []item, space string) error {
	ct := &ctf.ContentType{
		Sys:          &ctf.Sys{ID: "attachment"},
		Name:         "Attachment",
		Description:  "A file attachment or image",
		DisplayField: "slug",
		Fields: []*ctf.Field{
			&ctf.Field{
				ID:   "slug",
				Name: "Slug",
				Type: ctf.FieldTypeText,
			},
		},
	}

	fmt.Println("creating new 'attachment' content type")
	if err := cma.ContentTypes.Upsert(space, ct); err != nil {
		fmt.Println(err.Error())
		return err
	}

	fmt.Println("activating new 'attachment' content type")
	if err := cma.ContentTypes.Activate(space, ct); err != nil {
		fmt.Println(err.Error())
		return err
	}

	fmt.Printf("creating %d new attachments\n", len(items))
	for _, item := range items {
		// Skip anything which is not an attachment
		if !strings.Contains(item.Link, "attachment_id") {
			continue
		}

		entry := &ctf.Entry{
			Sys: &ctf.Sys{
				ID:          item.PostName,
				ContentType: ct,
			},
			Fields: map[string]interface{}{
				"slug": map[string]string{
					"en-US": item.PostName,
				},
			},
		}

		fmt.Printf("creating new attachment with ID %s\n", item.PostName)
		if err := cma.Entries.Upsert(space, entry); err != nil {
			return err
		}
		if err := cma.Entries.Publish(space, entry); err != nil {
			return err
		}
	}

	return nil
}
