package main

import (
	"flag"
	"fmt"
	"io/ioutil"

	ctf "github.com/ohookins/contentful-go"
)

// Map of original wordpress URLs to new CTF URLs for attachments and posts.
var replacementURLs map[string]string

func main() {
	replacementURLs = make(map[string]string)

	filename := flag.String("filename", "", "Filename of XML export of Wordpress blog")
	space := flag.String("space", "", "Space ID on Contentful")
	token := flag.String("token", "", "Personal CMA Token")
	skipTags := flag.Bool("skiptags", false, "Skip tag deletion/creation")
	skipPosts := flag.Bool("skipposts", false, "Skip post deletion/creation")
	skipAssets := flag.Bool("skipassets", false, "Skip asset/attachment deletion/creation")
	flag.Parse()

	if *filename == "" {
		fmt.Println("Please supply a filename input")
		return
	}

	if *space == "" {
		fmt.Println("Please supply a Space ID")
		return
	}

	if *token == "" {
		fmt.Println("Please supply a CMA Token")
		return
	}
	cma := ctf.NewCMA(*token)

	file, err := ioutil.ReadFile(*filename)
	if err != nil {
		fmt.Println(err)
		return
	}

	body, err := parseDoc(file)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = deleteContentAndType(cma, *space, "author")
	if err != nil {
		fmt.Printf("%v\n", err)
	}

	err = createAuthors(cma, body.Channel.Author, *space)
	if err != nil {
		fmt.Println(err)
	}

	err = deleteContentAndType(cma, *space, "category")
	if err != nil {
		fmt.Printf("%v\n", err)
	}

	err = createCategories(cma, body.Channel.WPCategory, *space)
	if err != nil {
		fmt.Println(err)
	}

	if !*skipTags {
		err = deleteContentAndType(cma, *space, "tag")
		if err != nil {
			fmt.Printf("%v\n", err)
		}

		err = createTags(cma, body.Channel.WPTag, *space)
		if err != nil {
			fmt.Println(err)
		}
	}

	if !*skipAssets {
		err = deleteAssets(cma, *space)
		if err != nil {
			fmt.Printf("%v\n", err)
		}

		err = createAttachments(cma, body.Channel.Item, *space)
		if err != nil {
			fmt.Println(err)
		}
	}

	if !*skipPosts {
		err = deleteContentAndType(cma, *space, "post")
		if err != nil {
			fmt.Printf("%v\n", err)
		}

		err = createPosts(cma, body.Channel.Item, *space)
		if err != nil {
			fmt.Println(err)
		}
	}

	// Generate a redirects file for Netlify from the replacement URLs map
	generateRedirects(replacementURLs)
}
