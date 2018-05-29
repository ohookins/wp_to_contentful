package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"

	ctf "github.com/ohookins/contentful-go"
)

var (
	sanitisePostRE  = regexp.MustCompile(`\\(["'.$^#_-~\n<>])`)
	imageResizeRE   = regexp.MustCompile(`\{\.alignleft[^}]+\}`)
	postContentType = &ctf.ContentType{
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
				ID:   "title",
				Name: "Title",
				Type: ctf.FieldTypeText,
			},
			&ctf.Field{
				ID:       "author",
				Name:     "Author",
				Type:     ctf.FieldTypeLink,
				LinkType: "Entry",
			},
			&ctf.Field{
				ID:   "published",
				Name: "Publishing Date",
				Type: ctf.FieldTypeDate,
			},
			&ctf.Field{
				ID:       "category",
				Name:     "Category",
				Type:     ctf.FieldTypeLink,
				LinkType: "Entry",
			},
			&ctf.Field{
				ID:   "tags",
				Name: "Tags",
				Type: ctf.FieldTypeArray,
				Items: &ctf.FieldTypeArrayItem{
					Type:     ctf.FieldTypeLink,
					LinkType: "Entry",
				},
			},
			&ctf.Field{
				ID:   "body",
				Name: "Body",
				Type: ctf.FieldTypeText,
			},
		},
	}
)

// post IDs can be no longer than 64 characters long
func createValidSlug(slug string) string {
	if len(slug) > 64 {
		return slug[:64]
	}

	return slug
}

// For some reason, the converted markdown contains a lot of escaped characters
func sanitisePost(body string) string {
	s := sanitisePostRE.ReplaceAllString(body, "$1")

	// images have some extraneous resizing styling which is unnecessary, e.g.:
	// {.alignleft .size-medium .wp-image-115 width="300" height="300"}
	return imageResizeRE.ReplaceAllString(s, "")
}

// Replace original wordpress URLs in the content with their Contentful
// counterparts.
func replaceURLs(body string) string {
	for origURL, newURL := range replacementURLs {
		body = strings.Replace(body, origURL, newURL, -1)
	}

	return body
}

func convertSourceDocumentLineEnds(content string) []byte {
	return []byte(strings.Replace(content, "\n", "<br/>", -1))
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
	fsrc.Write(convertSourceDocumentLineEnds(content[0]))

	// Call pandoc on the document to convert it
	cmd := exec.Command("pandoc", "--from", "html", "--to", "markdown", "-o", fdst.Name(), fsrc.Name())
	err = cmd.Run()
	if err != nil {
		fmt.Println(err.Error())
		return ""
	}

	md, _ := ioutil.ReadFile(fdst.Name())
	postBody := sanitisePost(string(md))
	postBody = replaceURLs(postBody)
	return postBody
}

func extractCategoryAndTags(categories []category) (cat string, tags []string) {
	for _, c := range categories {
		if c.Domain == "category" {
			cat = c.Name
		}

		if c.Domain == "post_tag" {
			tags = append(tags, c.Name)
		}
	}
	return
}

func reformatPubDate(wpdate string) string {
	// e.g. Sat, 11 Sep 2010 22:00:54 +0000
	t, err := time.Parse("Mon, 2 Jan 2006 15:04:05 -0700", wpdate)
	if err != nil {
		return ""
	}

	// ISO8601
	return t.Format("2006-01-02T15:04:05")
}

func generateTagArray(tags []string) []map[string]*ctf.Sys {
	result := []map[string]*ctf.Sys{}

	for _, tag := range tags {
		result = append(result, map[string]*ctf.Sys{
			"sys": &ctf.Sys{
				Type:     "Link",
				LinkType: "Entry",
				ID:       "tag_" + tag,
			},
		})
	}

	return result
}

func createPosts(cma *ctf.Contentful, items []item, space string) error {
	fmt.Println("creating new 'post' content type")
	if err := cma.ContentTypes.Upsert(space, postContentType); err != nil {
		fmt.Println(err.Error())
		return err
	}

	fmt.Println("activating new 'post' content type")
	if err := cma.ContentTypes.Activate(space, postContentType); err != nil {
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

		category, tags := extractCategoryAndTags(post.Categories)

		entry := &ctf.Entry{
			Sys: &ctf.Sys{
				ID:          finalSlug,
				ContentType: postContentType,
			},
			Fields: map[string]interface{}{
				"slug": map[string]string{
					"en-US": post.PostName,
				},
				"title": map[string]string{
					"en-US": post.Title,
				},
				"body": map[string]string{
					"en-US": content,
				},
				"author": map[string]interface{}{
					"en-US": map[string]*ctf.Sys{
						"sys": &ctf.Sys{
							Type:     "Link",
							LinkType: "Entry",
							ID:       "author_" + post.Creator,
						},
					},
				},
				"category": map[string]interface{}{
					"en-US": map[string]*ctf.Sys{
						"sys": &ctf.Sys{
							Type:     "Link",
							LinkType: "Entry",
							ID:       "cat_" + category,
						},
					},
				},
				"tags": map[string]interface{}{
					"en-US": generateTagArray(tags),
				},
				"published": map[string]string{
					"en-US": reformatPubDate(post.PubDate),
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

		// Add this entry to the list of URLs to replace. Entries are processed
		// in chronological order from the XML dump so if we have
		// back-references we should be able to replace all references.
		replacementURLs[strings.Replace(post.Guid, "http:", "https:", -1)] = "https://paperairoplane.net/" + post.PostName
		replacementURLs[post.Guid] = "https://paperairoplane.net/" + post.PostName
	}

	return nil
}
