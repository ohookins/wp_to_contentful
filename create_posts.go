package main

import (
	"fmt"
	"strings"

	ctf "github.com/ohookins/contentful-go"
)

func createPosts(cma *ctf.Contentful, items []item, space string) error {
	ct := &ctf.ContentType{
		Sys:          &ctf.Sys{ID: "post"},
		Name:         "Post",
		Description:  "A blog post/entry",
		DisplayField: "slug",
		Fields: []*ctf.Field{
			&ctf.Field{
				ID:   "slug",
				Name: "Slug",
				Type: ctf.FieldTypeText,
			},
		},
	}

	fmt.Println("creating new 'post' content type")
	if err := cma.ContentTypes.Upsert(space, ct); err != nil {
		fmt.Println(err.Error())
		return err
	}

	fmt.Println("activating new 'post' content type")
	if err := cma.ContentTypes.Activate(space, ct); err != nil {
		fmt.Println(err.Error())
		return err
	}

	fmt.Printf("creating %d new posts\n", len(items))
	for _, item := range items {
		// Skip attachments which are created separately
		if strings.Contains(item.Link, "attachment_id") {
			continue
		}

		entry := &ctf.Entry{
			Sys: &ctf.Sys{
				ID:          item.PostName,
				ContentType: ct,
			},
			Fields: map[string]ctf.LocalizedField{
				"slug": ctf.LocalizedField{
					"en-US": item.PostName,
				},
			},
		}

		fmt.Printf("creating new post with ID %s\n", item.PostName)
		if err := cma.Entries.Upsert(space, entry); err != nil {
			return err
		}
		if err := cma.Entries.Publish(space, entry); err != nil {
			return err
		}
	}

	return nil
}
