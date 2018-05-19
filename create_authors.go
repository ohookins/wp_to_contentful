package main

import (
	"fmt"

	ctf "github.com/ohookins/contentful-go"
)

func createAuthors(cma *ctf.Contentful, authors []author, space string) error {
	ct := &ctf.ContentType{
		Sys:          &ctf.Sys{ID: "author"},
		Name:         "Author",
		Description:  "Authors of content entries",
		DisplayField: "login",
		Fields: []*ctf.Field{
			&ctf.Field{
				ID:   "login",
				Name: "Login",
				Type: ctf.FieldTypeText,
			},
			&ctf.Field{
				ID:   "email",
				Name: "Email",
				Type: ctf.FieldTypeText,
			},
			&ctf.Field{
				ID:   "firstName",
				Name: "First Name",
				Type: ctf.FieldTypeText,
			},
			&ctf.Field{
				ID:   "lastName",
				Name: "First Name",
				Type: ctf.FieldTypeText,
			},
		},
	}

	fmt.Println("creating new 'author' content type")
	if err := cma.ContentTypes.Upsert(space, ct); err != nil {
		fmt.Println(err.Error())
		return err
	}

	fmt.Println("activating new 'author' content type")
	if err := cma.ContentTypes.Activate(space, ct); err != nil {
		fmt.Println(err.Error())
		return err
	}

	fmt.Printf("creating %d new authors\n", len(authors))
	for _, a := range authors {
		entry := &ctf.Entry{
			Sys: &ctf.Sys{
				ID:          "author_" + a.Login,
				ContentType: ct,
			},
			Fields: map[string]interface{}{
				"login": map[string]string{
					"en-US": a.Login,
				},
				"email": map[string]string{
					"en-US": a.Email,
				},
				"firstName": map[string]string{
					"en-US": a.FirstName,
				},
				"lastName": map[string]string{
					"en-US": a.LastName,
				},
			},
		}

		fmt.Printf("creating new author with ID author_%s\n", a.Login)
		if err := cma.Entries.Upsert(space, entry); err != nil {
			return err
		}
		if err := cma.Entries.Publish(space, entry); err != nil {
			return err
		}
	}

	return nil
}

func deleteAuthors(cma *ctf.Contentful, space string) error {
	ct := &ctf.ContentType{Sys: &ctf.Sys{ID: "author"}}

	collection := cma.Entries.List(space)
	collection.Query.ContentType("author")

	for {
		collection.Next()
		if len(collection.Items) == 0 {
			break
		}

		for _, entry := range collection.ToEntry() {
			fmt.Printf("deleting author with ID %s\n", entry.Sys.ID)
			_ = cma.Entries.Unpublish(space, entry)
			_ = cma.Entries.Delete(space, entry.Sys.ID)
		}
	}

	fmt.Println("deactivating 'author' content type")
	if err := cma.ContentTypes.Deactivate(space, ct); err != nil {
		if _, ok := err.(ctf.NotFoundError); !ok {
			return err
		}
	}

	fmt.Println("deleting 'author' content type")
	if err := cma.ContentTypes.Delete(space, ct); err != nil {
		if _, ok := err.(ctf.NotFoundError); !ok {
			return err
		}
	}
	return nil
}
