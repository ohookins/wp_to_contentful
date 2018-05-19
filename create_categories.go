package main

import (
	"fmt"

	ctf "github.com/ohookins/contentful-go"
)

func createCategories(cma *ctf.Contentful, categories []wpcategory, space string) error {
	ct := &ctf.ContentType{
		Sys:          &ctf.Sys{ID: "category"},
		Name:         "Category",
		Description:  "Main categories/themes of content",
		DisplayField: "realname",
		Fields: []*ctf.Field{
			&ctf.Field{
				ID:   "realname",
				Name: "Realname",
				Type: ctf.FieldTypeText,
			},
		},
	}

	fmt.Println("creating new 'category' content type")
	if err := cma.ContentTypes.Upsert(space, ct); err != nil {
		fmt.Println(err.Error())
		return err
	}

	fmt.Println("activating new 'category' content type")
	if err := cma.ContentTypes.Activate(space, ct); err != nil {
		fmt.Println(err.Error())
		return err
	}

	fmt.Printf("creating %d new categories\n", len(categories))
	for _, category := range categories {
		entry := &ctf.Entry{
			Sys: &ctf.Sys{
				ID:          "cat_" + category.NiceName,
				ContentType: ct,
			},
			Fields: map[string]interface{}{
				"realname": map[string]string{
					"en-US": category.CatName,
				},
			},
		}

		fmt.Printf("creating new category with ID cat_%s\n", category.NiceName)
		if err := cma.Entries.Upsert(space, entry); err != nil {
			return err
		}
		if err := cma.Entries.Publish(space, entry); err != nil {
			return err
		}
	}

	return nil
}

func deleteCategories(cma *ctf.Contentful, space string) error {
	ct := &ctf.ContentType{Sys: &ctf.Sys{ID: "category"}}

	collection := cma.Entries.List(space)
	collection.Query.ContentType("category")

	for {
		collection.Next()
		if len(collection.Items) == 0 {
			break
		}

		for _, entry := range collection.ToEntry() {
			fmt.Printf("deleting category with ID %s\n", entry.Sys.ID)
			_ = cma.Entries.Unpublish(space, entry)
			_ = cma.Entries.Delete(space, entry.Sys.ID)
		}
	}

	fmt.Println("deactivating 'category' content type")
	if err := cma.ContentTypes.Deactivate(space, ct); err != nil {
		if _, ok := err.(ctf.NotFoundError); !ok {
			return err
		}
	}

	fmt.Println("deleting 'category' content type")
	if err := cma.ContentTypes.Delete(space, ct); err != nil {
		if _, ok := err.(ctf.NotFoundError); !ok {
			return err
		}
	}
	return nil
}
