package main

import (
	"fmt"

	ctf "github.com/ohookins/contentful-go"
)

var tagContentType = &ctf.ContentType{
	Sys:          &ctf.Sys{ID: "tag"},
	Name:         "Tag",
	Description:  "Common tags for types of content",
	DisplayField: "realname",
	Fields: []*ctf.Field{
		&ctf.Field{
			ID:   "realname",
			Name: "Realname",
			Type: ctf.FieldTypeSymbol,
		},
	},
}

func createTags(cma *ctf.Contentful, tags []wptag, space string) error {
	fmt.Println("creating new 'tag' content type")
	if err := cma.ContentTypes.Upsert(space, tagContentType); err != nil {
		fmt.Println(err.Error())
		return err
	}

	fmt.Println("activating new 'tag' content type")
	if err := cma.ContentTypes.Activate(space, tagContentType); err != nil {
		fmt.Println(err.Error())
		return err
	}

	fmt.Printf("creating %d new tags\n", len(tags))
	for _, tag := range tags {
		entry := &ctf.Entry{
			Sys: &ctf.Sys{
				ID:          "tag_" + tag.Slug,
				ContentType: tagContentType,
			},
			Fields: map[string]interface{}{
				"realname": map[string]string{
					"en-US": tag.Name,
				},
			},
		}

		fmt.Printf("creating new tag with ID tag_%s\n", tag.Slug)
		if err := cma.Entries.Upsert(space, entry); err != nil {
			return err
		}
		if err := cma.Entries.Publish(space, entry); err != nil {
			return err
		}
	}

	return nil
}
