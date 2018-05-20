package main

import (
	"fmt"

	ctf "github.com/ohookins/contentful-go"
)

func deleteContentAndType(cma *ctf.Contentful, space, ctName string) error {
	ct := &ctf.ContentType{Sys: &ctf.Sys{ID: ctName}}

	collection := cma.Entries.List(space)
	collection.Query.ContentType(ctName)

	for {
		collection.Next()
		if len(collection.Items) == 0 {
			break
		}

		for _, entry := range collection.ToEntry() {
			fmt.Printf("deleting %s with ID %s\n", ctName, entry.Sys.ID)
			_ = cma.Entries.Unpublish(space, entry)
			_ = cma.Entries.Delete(space, entry.Sys.ID)
		}
	}

	fmt.Printf("deactivating '%s' content type\n", ctName)
	if err := cma.ContentTypes.Deactivate(space, ct); err != nil {
		if _, ok := err.(ctf.NotFoundError); !ok {
			return err
		}
	}

	fmt.Printf("deleting '%s' content type\n", ctName)
	if err := cma.ContentTypes.Delete(space, ct); err != nil {
		if _, ok := err.(ctf.NotFoundError); !ok {
			return err
		}
	}
	return nil
}
