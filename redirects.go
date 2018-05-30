package main

import (
	"fmt"
	"os"
	"regexp"
)

var sitePrefix = regexp.MustCompile(`^https?://paperairoplane.net`)

func generateRedirects(m map[string]string) {
	f, err := os.OpenFile("_redirects", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to open _redirects file for writing")
		return
	}
	defer f.Close()

	// Make another map for eliminating duplicates due to protocol changes
	uniquePaths := make(map[string]bool)

	// TODO: Probably need to generate redirects as well for WP redirect links
	// to attachments? Check access logs!
	for src, dst := range m {
		srcPath := sitePrefix.ReplaceAllString(src, "")
		dstPath := sitePrefix.ReplaceAllString(dst, "")

		if !uniquePaths[srcPath] {
			fmt.Fprintf(f, "%s %s\n", srcPath, dstPath)
			uniquePaths[srcPath] = true
		}
	}
}
