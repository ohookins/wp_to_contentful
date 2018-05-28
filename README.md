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
