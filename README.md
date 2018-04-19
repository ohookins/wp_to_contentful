# Summary

This is just a set of tooling to help me migrate my WordPress blog to
Contentful. Probably not very useful for the general case, but I'm
documenting how it is used in case I forget.

dep ensure
go build
./wp_to_contentful -filename <WordpressDump.xml> -space <SpaceID> -token <CMAToken>
