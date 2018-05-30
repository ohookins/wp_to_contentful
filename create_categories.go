package main

import (
	"fmt"

	ctf "github.com/ohookins/contentful-go"
)

var categoryContentType = &ctf.ContentType{
	Sys:          &ctf.Sys{ID: "category"},
	Name:         "Category",
	Description:  "Main categories/themes of content",
	DisplayField: "realname",
	Fields: []*ctf.Field{
		&ctf.Field{
			ID:   "realname",
			Name: "Realname",
			Type: ctf.FieldTypeSymbol,
		},
	},
}

func createCategories(cma *ctf.Contentful, categories []wpcategory, space string) error {
	fmt.Println("creating new 'category' content type")
	if err := cma.ContentTypes.Upsert(space, categoryContentType); err != nil {
		fmt.Println("creating category error: ", err.Error())
		return err
	}

	fmt.Println("activating new 'category' content type")
	if err := cma.ContentTypes.Activate(space, categoryContentType); err != nil {
		fmt.Println("activating category error: ", err.Error())
		return err
	}

	fmt.Printf("creating %d new categories\n", len(categories))
	for _, category := range categories {
		entry := &ctf.Entry{
			Sys: &ctf.Sys{
				ID:          "cat_" + category.NiceName,
				ContentType: categoryContentType,
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
