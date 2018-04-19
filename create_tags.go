package main

import (
	"fmt"

	ctf "github.com/ohookins/contentful-go"
)

func createTags(cma *ctf.Contentful, tags []wptag, space string) error {
	ct := &ctf.ContentType{
		Sys:         &ctf.Sys{ID: "tag"},
		Name:        "Tag",
		Description: "Common tags for types of content",
		Fields: []*ctf.Field{
			&ctf.Field{
                ID: "realname",
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
                ID: tag.Slug,
                ContentType: ct,
            },
		}

		fmt.Printf("creating new tag with ID %s\n", tag.Slug)
		if err := cma.Entries.Publish(space, entry); err != nil {
            fmt.Println(err)
            return err
        }
	}

	return nil
}

func deleteTags(cma *ctf.Contentful, space string) error {
	ct := &ctf.ContentType{Sys: &ctf.Sys{ID: "tag"}}

	collection := cma.Entries.List(space)
	collection.Query.ContentType("tag")

	for {
		collection.Next()
		if len(collection.Items) == 0 {
			break
		}

		for _, entry := range collection.ToEntry() {
			fmt.Printf("deleting tag with ID %s\n", entry.Sys.ID)
			_ = cma.Entries.Unpublish(space, entry)
			_ = cma.Entries.Delete(space, entry.Sys.ID)
		}
	}

	fmt.Println("deactivating 'tag' content type")
	if err := cma.ContentTypes.Deactivate(space, ct); err != nil {
		if _, ok := err.(ctf.NotFoundError); !ok {
			return err
		}
	}

	fmt.Println("deleting 'tag' content type")
	if err := cma.ContentTypes.Delete(space, ct); err != nil {
		if _, ok := err.(ctf.NotFoundError); !ok {
			return err
		}
	}
	return nil
}
