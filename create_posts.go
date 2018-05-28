package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	ctf "github.com/ohookins/contentful-go"
)

// post IDs can be no longer than 64 characters long
func createValidSlug(slug string) string {
	if len(slug) > 64 {
		return slug[:64]
	}

	return slug
}

func convertToMarkdown(content []string) string {
	fsrc, err := ioutil.TempFile(".", "wpconvert")
	if err != nil {
		fmt.Println(err.Error())
		return ""
	}
	defer os.Remove(fsrc.Name())
	defer fsrc.Close()

	fdst, err := ioutil.TempFile(".", "md")
	if err != nil {
		fmt.Println(err.Error())
		return ""
	}
	defer os.Remove(fdst.Name())
	defer fdst.Close()

	// posts seem to have two content sections, the second is empty
	fsrc.Write([]byte(content[0]))

	// Call pandoc on the document to convert it
	cmd := exec.Command("pandoc", "--from", "html", "--to", "markdown", "-o", fdst.Name(), fsrc.Name())
	err = cmd.Run()
	if err != nil {
		fmt.Println(err.Error())
		return ""
	}

	md, _ := ioutil.ReadFile(fdst.Name())
	return string(md)
}

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
			&ctf.Field{
				ID:   "body",
				Name: "Body",
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

	// Extract just the posts and not attachments
	posts := []item{}
	for _, item := range items {
		// Skip attachments which are created separately
		// Also skip draft posts
		if !strings.Contains(item.Link, "attachment_id") && item.Status != "draft" {
			posts = append(posts, item)
		}
	}
	fmt.Printf("creating %d new posts\n", len(posts))

	for _, post := range posts {
		content := convertToMarkdown(post.Content)

		finalSlug := createValidSlug(post.PostName)

		entry := &ctf.Entry{
			Sys: &ctf.Sys{
				ID:          finalSlug,
				ContentType: ct,
			},
			Fields: map[string]ctf.LocalizedField{
				"slug": ctf.LocalizedField{
					"en-US": post.PostName,
				},
				"body": ctf.LocalizedField{
					"en-US": content,
				},
			},
		}

		fmt.Printf("creating new post with ID %s\n", finalSlug)
		if err := cma.Entries.Upsert(space, entry); err != nil {
			return err
		}
		if err := cma.Entries.Publish(space, entry); err != nil {
			return err
		}
	}

	return nil
}
