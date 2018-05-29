# Summary

This is just a set of tooling to help me migrate my WordPress blog to
Contentful. Probably not very useful for the general case, but I'm
documenting how it is used in case I forget.

## Prerequisites

You need to have the following installed:
 * `dep` for go dependency management
 * `pandoc` (e.g. from homebrew) for HTML to Markdown conversion

## Building/Running

```
dep ensure
go build
./wp_to_contentful -filename <WordpressDump.xml> -space <SpaceID> -token <CMAToken>
```

The following are created:
 * tags
 * authors
 * categories
 * assets (both images and other attachments)
 * blog posts

Each of the above with exception of assets get their own content type. Posts
make use of links to link to the entries for tags/categories/authors rather
than duplicating the data directly in the post. The data in these other content
types is minimal, but it's better to do it this way for content searching.

Additionally, a Netlify-compatible `_redirects` file is generated from the
renamed paths of posts and assets.
