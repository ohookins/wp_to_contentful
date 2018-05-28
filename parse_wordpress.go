package main

import (
	"encoding/xml"
)

type rss struct {
	XMLName xml.Name `xml:"rss"`
	Channel channel  `xml:"channel"`
}

type channel struct {
	Title      string       `xml:"title"`
	Author     []author     `xml:"author"`
	WPCategory []wpcategory `xml:"category"`
	WPTag      []wptag      `xml:"tag"`
	// Ignore terms, they don't add any value.
	Item []item `xml:"item"`
}

type author struct {
	Id        int    `xml:"author_id"`
	Login     string `xml:"author_login",cdata`
	Email     string `xml:"author_email",cdata`
	FirstName string `xml:"author_first_name",cdata`
	LastName  string `xml:"author_last_name",cdata`
}

type wpcategory struct {
	NiceName string `xml:"category_nicename",cdata`
	CatName  string `xml:"cat_name",cdata`
}

type wptag struct {
	Slug string `xml:"tag_slug",cdata`
	Name string `xml:"tag_name",cdata`
}

// item is any item of content - it can be an image/attachment/media or a
// blog post.
type item struct {
	Title   string `xml:"title"`
	Link    string `xml:"link"`
	PubDate string `xml:"pubDate"`       // non-default format, need to parse later
	Creator string `xml:"creator",cdata` // Encoded author string
	Status  string `xml:"status",cdata`  // Encoded status string - published/draft

	// Permalink of blog post, or full URL to binary if an image/attachement
	Guid string `xml:"guid"`

	// Post this image/attachment is linked from
	PostID int `xml:"post_id"`

	// Body of post
	Content []string `xml:"encoded",cdata`

	// Effectively the slug
	PostName string `xml:"post_name",cdata`

	// Categories and tags
	Categories []category `xml:"category"`
}

type category struct {
	Domain string `xml:"domain,attr"`
	Name   string `xml:"nicename,attr"`
}

func parseDoc(xmldoc []byte) (*rss, error) {
	body := &rss{}

	err := xml.Unmarshal(xmldoc, body)

	return body, err
}
