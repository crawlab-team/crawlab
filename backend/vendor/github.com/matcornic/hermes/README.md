# Hermes

[![Build Status](https://travis-ci.org/matcornic/hermes.svg?branch=master)](https://travis-ci.org/matcornic/hermes)
[![Go Report Card](https://goreportcard.com/badge/github.com/matcornic/hermes)](https://goreportcard.com/report/github.com/matcornic/hermes)
[![Go Coverage](https://codecov.io/github/matcornic/hermes/coverage.svg)](https://codecov.io/github/matcornic/hermes/)
[![Godoc](https://godoc.org/github.com/matcornic/hermes?status.svg)](https://godoc.org/github.com/matcornic/hermes)

Hermes is the Go port of the great [mailgen](https://github.com/eladnava/mailgen) engine for Node.js. Check their work, it's awesome!
It's a package that generates clean, responsive HTML e-mails for sending transactional e-mails (welcome e-mails, reset password e-mails, receipt e-mails and so on), and associated plain text fallback.

# Demo

<img src="screens/default/welcome.png" height="400" /> <img src="screens/default/reset.png" height="400" /> <img src="screens/default/receipt.png" height="400" />

# Usage

First install the package:

```
go get -u github.com/matcornic/hermes
```

Then, start using the package by importing and configuring it:

```go
// Configure hermes by setting a theme and your product info
h := hermes.Hermes{
    // Optional Theme
    // Theme: new(Default) 
    Product: hermes.Product{
        // Appears in header & footer of e-mails
        Name: "Hermes",
        Link: "https://example-hermes.com/",
        // Optional product logo
        Logo: "http://www.duchess-france.org/wp-content/uploads/2016/01/gopher.png",
    },
}
```

Next, generate an e-mail using the following code:

```go
email := hermes.Email{
    Body: hermes.Body{
        Name: "Jon Snow",
        Intros: []string{
            "Welcome to Hermes! We're very excited to have you on board.",
        },
        Actions: []hermes.Action{
            {
                Instructions: "To get started with Hermes, please click here:",
                Button: hermes.Button{
                    Color: "#22BC66", // Optional action button color
                    Text:  "Confirm your account",
                    Link:  "https://hermes-example.com/confirm?token=d9729feb74992cc3482b350163a1a010",
                },
            },
        },
        Outros: []string{
            "Need help, or have questions? Just reply to this email, we'd love to help.",
        },
    },
}

// Generate an HTML email with the provided contents (for modern clients)
emailBody, err := h.GenerateHTML(email)
if err != nil {
    panic(err) // Tip: Handle error with something else than a panic ;)
}

// Generate the plaintext version of the e-mail (for clients that do not support xHTML)
emailText, err := h.GeneratePlainText(email)
if err != nil {
    panic(err) // Tip: Handle error with something else than a panic ;)
}

// Optionally, preview the generated HTML e-mail by writing it to a local file
err = ioutil.WriteFile("preview.html", []byte(emailBody), 0644)
if err != nil {
    panic(err) // Tip: Handle error with something else than a panic ;)
}
```

This code would output the following HTML template:

<img src="screens/demo.png" height="400" />

And the following plain text:

```
------------
Hi Jon Snow,
------------

Welcome to Hermes! We're very excited to have you on board.

To get started with Hermes, please click here: https://hermes-example.com/confirm?token=d9729feb74992cc3482b350163a1a010

Need help, or have questions? Just reply to this email, we'd love to help.

Yours truly,
Hermes - https://example-hermes.com/

Copyright © 2017 Hermes. All rights reserved.
```

> Theme templates will be embedded in your application binary. If you want to use external templates (for configuration), use your own theme by implementing `hermes.Theme` interface with code searching for your files.

## More Examples

* [Welcome](examples/welcome.go)
* [Receipt](examples/receipt.go)
* [Password Reset](examples/reset.go)
* [Maintenance](examples/maintenance.go)

To run the examples, go to `examples` folder, then run `go run *.go`. HTML and Plaintext example should be created in given theme folders.

## Plaintext E-mails

To generate a [plaintext version of the e-mail](https://litmus.com/blog/best-practices-for-plain-text-emails-a-look-at-why-theyre-important), simply call `GeneratePlainText` function:

```go
// Generate plaintext email using hermes
emailText, err := h.GeneratePlainText(email)
if err != nil {
    panic(err) // Tip: Handle error with something else than a panic ;)
}
```

## Supported Themes

The following open-source themes are bundled with this package:

* `default` by [Postmark Transactional Email Templates](https://github.com/wildbit/postmark-templates)

<img src="screens/default/welcome.png" height="200" /> <img src="screens/default/reset.png" height="200" /> <img src="screens/default/receipt.png" height="200" />

* `flat`, slightly modified from [Postmark Transactional Email Templates](https://github.com/wildbit/postmark-templates)

<img src="screens/flat/welcome.png" height="200" /> <img src="screens/flat/reset.png" height="200" /> <img src="screens/flat/receipt.png" height="200" />

## RTL Support

To change the default text direction (left-to-right), simply override it as follows:

```go
// Configure hermes by setting a theme and your product info
h := hermes.Hermes {
    // Custom text direction
    TextDirection: hermes.TDRightToLeft,
}
```

## Language Customizations

To customize the e-mail's greeting ("Hi") or signature ("Yours truly"), supply custom strings within the e-mail's `Body`:

```go
email := hermes.Email{
    Body: hermes.Body{
        Greeting: "Dear",
        Signature: "Sincerely",
    },
}
```

To use a custom title string rather than a greeting/name introduction, provide it instead of `Name`:

```go
email := hermes.Email{
    Body: hermes.Body{
        // Title will override `Name`
        Title: "Welcome to Hermes",
    },
}
```

To customize the `Copyright`, override it when initializing `Hermes` within your `Product` as follows:

```go
// Configure hermes by setting a theme and your product info
h := hermes.Hermes{
    // Optional Theme
    // Theme: new(Default)
    Product: hermes.Product{
        // Appears in header & footer of e-mails
        Name: "Hermes",
        Link: "https://example-hermes.com/",
        // Custom copyright notice
        Copyright: "Copyright © 2017 Dharma Initiative. All rights reserved."
    },
}
```

To use a custom fallback text at the end of the email, change the `TroubleText` field of the `hermes.Product` struct. The default value is `If you’re having trouble with the button '{ACTION}', copy and paste the URL below into your web browser.`. The `{ACTION}` placeholder will be replaced with the corresponding text of the supplied action button:

```go
// Configure hermes by setting a theme and your product info
h := hermes.Hermes{
    // Optional Theme
    // Theme: new(Default)
    Product: hermes.Product{
        // Custom trouble text
        TroubleText: "If the {ACTION}-button is not working for you, just copy and paste the URL below into your web browser."
    },
}
```

## Elements

Hermes supports injecting custom elements such as dictionaries, tables and action buttons into e-mails.

### Action

To inject an action button in to the e-mail, supply the `Actions` object as follows:

```go
email := hermes.Email{
    Body: hermes.Body{
        Actions: []hermes.Action{
            {
                Instructions: "To get started with Hermes, please click here:",
                Button: hermes.Button{
                    Color: "#22BC66", // Optional action button color
                    Text:  "Confirm your account",
                    Link:  "https://hermes-example.com/confirm?token=d9729feb74992cc3482b350163a1a010",
                },
            },
        },
    },
}
```

To inject multiple action buttons in to the e-mail, supply another struct in Actions slice `Action`.

### Table

To inject a table into the e-mail, supply the `Table` object as follows:

```go
email := hermes.Email{
    Body: hermes.Body{
        Table: hermes.Table{
            Data: [][]hermes.Entry{
                // List of rows
                {   
                    // Key is the column name, Value is the cell value
                    // First object defines what columns will be displayed
                    {Key: "Item", Value: "Golang"},
                    {Key: "Description", Value: "Open source programming language that makes it easy to build simple, reliable, and efficient software"},
                    {Key: "Price", Value: "$10.99"},
                },
                {
                    {Key: "Item", Value: "Hermes"},
                    {Key: "Description", Value: "Programmatically create beautiful e-mails using Golang."},
                    {Key: "Price", Value: "$1.99"},
                },
            },
            Columns: hermes.Columns{
                // Custom style for each rows
                CustomWidth: map[string]string{
                    "Item":  "20%",
                    "Price": "15%",
                },
                CustomAlignment: map[string]string{
                    "Price": "right",
                },
            },
        },
    },
}
```

### Dictionary

To inject key-value pairs of data into the e-mail, supply the `Dictionary` object as follows:

```go
email := hermes.Email{
    Body: hermes.Body{
        Dictionary: []hermes.Entry{
            {Key: "Date", Value: "20 November 1887"},
            {Key: "Address", Value: "221B Baker Street, London"},
        },
    },
}
```

### Free Markdown

If you need more flexibility in the content of your generated e-mail, while keeping the same format than any other e-mail, use Markdown content. Supply the `FreeMarkdown` object as follows:

```go
email := hermes.Email{
		Body: hermes.Body{
			FreeMarkdown: `
> _Hermes_ service will shutdown the **1st August 2017** for maintenance operations. 

Services will be unavailable based on the following schedule:

| Services | Downtime |
| :------:| :-----------: |
| Service A | 2AM to 3AM |
| Service B | 4AM to 5AM |
| Service C | 5AM to 6AM |

---

Feel free to contact us for any question regarding this matter at [support@hermes-example.com](mailto:support@hermes-example.com) or in our [Gitter](https://gitter.im/)

`,
		},
	}
}
```

This code would output the following HTML template:

<img src="screens/free-markdown.png" height="400" />

And the following plaintext:

```
------------
Hi Jon Snow,
------------

> 
> 
> 
> Hermes service will shutdown the *1st August 2017* for maintenance
> operations.
> 
> 

Services will be unavailable based on the following schedule:

+-----------+------------+
| SERVICES  |  DOWNTIME  |
+-----------+------------+
| Service A | 2AM to 3AM |
| Service B | 4AM to 5AM |
| Service C | 5AM to 6AM |
+-----------+------------+

Feel free to contact us for any question regarding this matter at support@hermes-example.com ( support@hermes-example.com ) or in our Gitter ( https://gitter.im/ )

Yours truly,
Hermes - https://example-hermes.com/

Copyright © 2017 Hermes. All rights reserved.
```

Be aware that this content will replace existing tables, dictionary and actions. Only intros, outros, header and footer will be kept.

This is helpful when your application needs sending e-mails, wrote on-the-fly by adminstrators.

> Markdown is rendered with [Blackfriday](https://github.com/russross/blackfriday), so every thing Blackfriday can do, Hermes can do it as well.

## Troubleshooting

1. After sending multiple e-mails to the same Gmail / Inbox address, they become grouped and truncated since they contain similar text, breaking the responsive e-mail layout.

> Simply sending the `X-Entity-Ref-ID` header with your e-mails will prevent grouping / truncation.

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md)

## License

Apache 2.0

