package main

import (
	"fmt"

	ctf "github.com/ohookins/contentful-go"
)

func createTags(cma *ctf.Contentful, tags []wptag, space string) error {
	ct := &ctf.ContentType{
		Sys:          &ctf.Sys{ID: "tag"},
		Name:         "Tag",
		Description:  "Common tags for types of content",
		DisplayField: "realname",
		Fields: []*ctf.Field{
			&ctf.Field{
				ID:   "realname",
				Name: "Realname",
				Type: ctf.FieldTypeText,
			},
		},
	}

	fmt.Println("creating new 'tag' content type")
	if err := cma.ContentTypes.Upsert(space, ct); err != nil {
		fmt.Println(err.Error())
		return err
	}

	fmt.Println("activating new 'tag' content type")
	if err := cma.ContentTypes.Activate(space, ct); err != nil {
		fmt.Println(err.Error())
		return err
	}

	fmt.Printf("creating %d new tags\n", len(tags))
	for _, tag := range tags {
		entry := &ctf.Entry{
			Sys: &ctf.Sys{
				ID:          "tag_" + tag.Slug,
				ContentType: ct,
			},
			Fields: map[string]ctf.LocalizedField{
				"realname": ctf.LocalizedField{
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
