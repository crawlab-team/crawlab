package hermes

import (
	"bytes"
	"github.com/Masterminds/sprig"
	"github.com/imdario/mergo"
	"github.com/jaytaylor/html2text"
	"gopkg.in/russross/blackfriday.v2"
	"html/template"
)

// Hermes is an instance of the hermes email generator
type Hermes struct {
	Theme         Theme
	TextDirection TextDirection
	Product       Product
}

// Theme is an interface to implement when creating a new theme
type Theme interface {
	Name() string              // The name of the theme
	HTMLTemplate() string      // The golang template for HTML emails
	PlainTextTemplate() string // The golang templte for plain text emails (can be basic HTML)
}

// TextDirection of the text in HTML email
type TextDirection string

var templateFuncs = template.FuncMap{
	"url": func(s string) template.URL {
		return template.URL(s)
	},
}

// TDLeftToRight is the text direction from left to right (default)
const TDLeftToRight TextDirection = "ltr"

// TDRightToLeft is the text direction from right to left
const TDRightToLeft TextDirection = "rtl"

// Product represents your company product (brand)
// Appears in header & footer of e-mails
type Product struct {
	Name        string
	Link        string // e.g. https://matcornic.github.io
	Logo        string // e.g. https://matcornic.github.io/img/logo.png
	Copyright   string // Copyright © 2017 Hermes. All rights reserved.
	TroubleText string // TroubleText is the sentence at the end of the email for users having trouble with the button (default to `If you’re having trouble with the button '{ACTION}', copy and paste the URL below into your web browser.`)
}

// Email is the email containing a body
type Email struct {
	Body Body
}

// Markdown is a HTML template (a string) representing Markdown content
// https://en.wikipedia.org/wiki/Markdown
type Markdown template.HTML

// Body is the body of the email, containing all interesting data
type Body struct {
	Name         string   // The name of the contacted person
	Intros       []string // Intro sentences, first displayed in the email
	Dictionary   []Entry  // A list of key+value (useful for displaying parameters/settings/personal info)
	Table        Table    // Table is an table where you can put data (pricing grid, a bill, and so on)
	Actions      []Action // Actions are a list of actions that the user will be able to execute via a button click
	Outros       []string // Outro sentences, last displayed in the email
	Greeting     string   // Greeting for the contacted person (default to 'Hi')
	Signature    string   // Signature for the contacted person (default to 'Yours truly')
	Title        string   // Title replaces the greeting+name when set
	FreeMarkdown Markdown // Free markdown content that replaces all content other than header and footer
}

// ToHTML converts Markdown to HTML
func (c Markdown) ToHTML() template.HTML {
	return template.HTML(blackfriday.Run([]byte(string(c))))
}

// Entry is a simple entry of a map
// Allows using a slice of entries instead of a map
// Because Golang maps are not ordered
type Entry struct {
	Key   string
	Value string
}

// Table is an table where you can put data (pricing grid, a bill, and so on)
type Table struct {
	Data    [][]Entry // Contains data
	Columns Columns   // Contains meta-data for display purpose (width, alignement)
}

// Columns contains meta-data for the different columns
type Columns struct {
	CustomWidth     map[string]string
	CustomAlignment map[string]string
}

// Action is an action the user can do on the email (click on a button)
type Action struct {
	Instructions string
	Button       Button
}

// Button defines an action to launch
type Button struct {
	Color     string
	TextColor string
	Text      string
	Link      string
}

// Template is the struct given to Golang templating
// Root object in a template is this struct
type Template struct {
	Hermes Hermes
	Email  Email
}

func setDefaultEmailValues(e *Email) error {
	// Default values of an email
	defaultEmail := Email{
		Body: Body{
			Intros:     []string{},
			Dictionary: []Entry{},
			Outros:     []string{},
			Signature:  "Yours truly",
			Greeting:   "Hi",
		},
	}
	// Merge the given email with default one
	// Default one overrides all zero values
	return mergo.Merge(e, defaultEmail)
}

// default values of the engine
func setDefaultHermesValues(h *Hermes) error {
	defaultTextDirection := TDLeftToRight
	defaultHermes := Hermes{
		Theme:         new(Default),
		TextDirection: defaultTextDirection,
		Product: Product{
			Name:        "Hermes",
			Copyright:   "Copyright © 2017 Hermes. All rights reserved.",
			TroubleText: "If you’re having trouble with the button '{ACTION}', copy and paste the URL below into your web browser.",
		},
	}
	// Merge the given hermes engine configuration with default one
	// Default one overrides all zero values
	err := mergo.Merge(h, defaultHermes)
	if err != nil {
		return err
	}
	if h.TextDirection != TDLeftToRight && h.TextDirection != TDRightToLeft {
		h.TextDirection = defaultTextDirection
	}
	return nil
}

// GenerateHTML generates the email body from data to an HTML Reader
// This is for modern email clients
func (h *Hermes) GenerateHTML(email Email) (string, error) {
	err := setDefaultHermesValues(h)
	if err != nil {
		return "", err
	}
	return h.generateTemplate(email, h.Theme.HTMLTemplate())
}

// GeneratePlainText generates the email body from data
// This is for old email clients
func (h *Hermes) GeneratePlainText(email Email) (string, error) {
	err := setDefaultHermesValues(h)
	if err != nil {
		return "", err
	}
	template, err := h.generateTemplate(email, h.Theme.PlainTextTemplate())
	if err != nil {
		return "", err
	}
	return html2text.FromString(template, html2text.Options{PrettyTables: true})
}

func (h *Hermes) generateTemplate(email Email, tplt string) (string, error) {

	err := setDefaultEmailValues(&email)
	if err != nil {
		return "", err
	}

	// Generate the email from Golang template
	// Allow usage of simple function from sprig : https://github.com/Masterminds/sprig
	t, err := template.New("hermes").Funcs(sprig.FuncMap()).Funcs(templateFuncs).Parse(tplt)
	if err != nil {
		return "", err
	}
	var b bytes.Buffer
	t.Execute(&b, Template{*h, email})
	return b.String(), nil
}
