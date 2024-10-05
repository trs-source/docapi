package docapi

type Contact struct {
	Name  string `json:"name,omitempty"`
	URL   string `json:"url,omitempty"`
	Email string `json:"email,omitempty"`
}

type OptsContact func(*Contact)

func WithContactWebSite(url string) OptsContact {
	return func(c *Contact) {
		c.URL = url
	}
}

func WithContactEmail(email string) OptsContact {
	return func(c *Contact) {
		c.Email = email
	}
}
