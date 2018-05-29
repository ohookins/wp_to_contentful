package main

import (
	"fmt"

	ctf "github.com/ohookins/contentful-go"
)

var authorContentType = &ctf.ContentType{
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

func createAuthors(cma *ctf.Contentful, authors []author, space string) error {
	fmt.Println("creating new 'author' content type")
	if err := cma.ContentTypes.Upsert(space, authorContentType); err != nil {
		fmt.Println(err.Error())
		return err
	}

	fmt.Println("activating new 'author' content type")
	if err := cma.ContentTypes.Activate(space, authorContentType); err != nil {
		fmt.Println(err.Error())
		return err
	}

	fmt.Printf("creating %d new authors\n", len(authors))
	for _, a := range authors {
		entry := &ctf.Entry{
			Sys: &ctf.Sys{
				ID:          "author_" + a.Login,
				ContentType: authorContentType,
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
