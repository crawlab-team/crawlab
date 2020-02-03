# Contributing

Thanks so much for wanting to help! We really appreciate it.

* Have an idea for a new feature?
* Want to add a new built-in theme?

Excellent! You've come to the right place.

1. If you find a bug or wish to suggest a new feature, please create an issue first
2. Make sure your code & comment conventions are in-line with the project's style (execute gometalinter as in [.travis.yml](.travis.yml) file)
3. Make your commits and PRs as tiny as possible - one feature or bugfix at a time
4. Write detailed commit messages, in-line with the project's commit naming conventions

## Theming Instructions

This file contains instructions on adding themes to Hermes:

* [Using a Custom Theme](#using-a-custom-theme)
* [Creating a Built-In Theme](#creating-a-built-in-theme)

> We use Golang templates under the hood to inject the e-mail body into themes.
> - [Official guide](https://golang.org/pkg/text/template/)
> - [Tutorial](https://astaxie.gitbooks.io/build-web-application-with-golang/en/07.4.html)
> - [Hugo guide](https://gohugo.io/templates/go-templates/)

### Using a Custom Theme

If you want to supply your own **custom theme** for Hermes to use (but don't want it included with Hermes):

1. Create a new struct implementing `Theme` interface ([hermes.go](hermes.go)). A real-life example is in [default.go](default.go)
2. Supply your new theme at hermes creation

```go

type MyCustomTheme struct{}

func (dt *MyCustomTheme) Name() string {
	return "mycustomthem"
}

func (dt *MyCustomTheme) HTMLTemplate() string {
    // Get the template from a file (if you want to be able to change the template live without retstarting your application)
    // Or write the template by returning pure string here (if you want embbeded template and do not bother with external dependencies)
    return "<A go html template with wanted information>" 
}

func (dt *MyCustomTheme) PlainTextTemplate() string {
    // Get the template from a file (if you want to be able to change the template live without retstarting your application)
    // Or write the template by returning pure string here (if you want embbeded template and do not bother with external dependencies)
    return "<A go plaintext template with wanter information>"
}

h := hermes.Hermes{
    Theme: new(MyCustomTheme) // Set your fresh new theme here
    Product: hermes.Product{
        Name: "Hermes",
        Link: "https://example-hermes.com/",
    },
}

// ...
// Continue with the rest as usual, create your email and generate the content.
// ...
```

3. That's it.

### Creating a Built-In Theme

If you want to create a new **built-in** Hermes theme:

1. Fork the repository to your GitHub account and clone it to your computer
2. Create a new Go file named after your new theme
3. Copy content of [default.go](default.go) file in new file and make any necessary changes
4. Scroll down to the [injection snippets](#injection-snippets) and copy and paste each code snippet into the relevant area of your template markup
5. Test the theme by adding the theme to slice of tested themes (see [hermes_test.go](hermes_test.go)). A set of tests will be run to check that your theme follows features of Hermes.
6. Create examples in new folder for your theme in `examples` folder and run `go run *.go`. It will generate the different `html` and `plaintext` emails for your different examples. Follow the same examples as default theme (3 examples: Welcome, Reset and Receipt)
7. Add the theme name, credit, and screenshots to the `README.md` file's [Supported Themes](README.md#supported-themes) section (copy one of the existing themes' markup and modify it accordingly)
8. Submit a pull request with your changes and we'll let you know if anything's missing!

Thanks again for your contribution!

# Injection Snippets

## Product Branding Injection

The following will inject either the product logo or name into the template.

```html
<a href="{{.Hermes.Product.Link}}" target="_blank">
        {{ if .Hermes.Product.Logo }}
            <img src="{{.Hermes.Product.Logo}}" class="email-logo" />
        {{ else }}
            {{ .Hermes.Product.Name }}
        {{ end }}
</a>
```

It's a good idea to add the following CSS declaration to set `max-height: 50px` for the logo:

```css
.email-logo {
    max-height: 50px;
}
```

## Title Injection

The following will inject the e-mail title (Hi John Appleseed,) or a custom title provided by the user:

```html
<h1>{{if .Email.Body.Title }}{{ .Email.Body.Title }}{{ else }}{{ .Email.Body.Greeting }} {{ .Email.Body.Name }},{{ end }}</h1>
```

## Intro Injection

The following will inject the intro text (string or array) into the e-mail:

```html
{{ with .Email.Body.Intros }}
    {{ if gt (len .) 0 }}
        {{ range $line := . }}
            <p>{{ $line }}</p>
        {{ end }}
    {{ end }}
{{ end }}
```

## Dictionary Injection

The following will inject a `<dl>` of key-value pairs into the e-mail:

```html
{{ with .Email.Body.Dictionary }} 
    {{ if gt (len .) 0 }}
        <dl class="body-dictionary">
        {{ range $entry := . }}
            <dt>{{ $entry.Key }}:</dt>
            <dd>{{ $entry.Value }}</dd>
        {{ end }}
        </dl>
    {{ end }}
{{ end }}
```

It's a good idea to add this to the top of the template to improve the styling of the dictionary:

```css
/* Dictionary */
.dictionary {
    width: 100%;
    overflow: hidden;
    margin: 0 auto;
    padding: 0;
}
.dictionary dt {
    clear: both;
    color: #000;
    font-weight: bold;
    margin-right: 4px;
}
.dictionary dd {
    margin: 0 0 10px 0;
}
```

## Table Injection

The following will inject the table into the e-mail:

```html
<!-- Table -->
{{ with .Email.Body.Table }}
{{ $data := .Data }}
{{ $columns := .Columns }}
    {{ if gt (len $data) 0 }}
    <table class="data-wrapper" width="100%" cellpadding="0" cellspacing="0">
        <tr>
        <td colspan="2">
            <table class="data-table" width="100%" cellpadding="0" cellspacing="0">
            <tr>
                {{ $col := index $data 0 }}
                {{ range $entry := $col }}
                <th
                    {{ with $columns }}
                        {{ $width := index .CustomWidth $entry.Key }}
                        {{ with $width }}
                        width="{{ . }}"
                        {{ end }}
                        {{ $align := index .CustomAlignment $entry.Key }}
                        {{ with $align }}
                        style="text-align:{{ . }}"
                        {{ end }}
                    {{ end }}
                >
                    <p>{{ $entry.Key }}</p>
                </th>
                {{ end }}
            </tr>
            {{ range $row := $data }}
            <tr>
                {{ range $cell := $row }}
                <td
                    {{ with $columns }}
                        {{ $align := index .CustomAlignment $cell.Key }}
                        {{ with $align }}
                        style="text-align:{{ . }}"
                        {{ end }}
                    {{ end }}
                >
                {{ $cell.Value }}
                </td>
                {{ end }}
            </tr>
            {{ end }}
            </table>
        </td>
        </tr>
    </table>
    {{ end }}
{{ end }}
```

It's a good idea to add this to the top of the template to improve the styling of the table:

```css
/* Table */
.data-wrapper {
    width: 100%;
    margin: 0;
    padding: 35px 0;
}
.data-table {
    width: 100%;
    margin: 0;
}
.data-table th {
    text-align: left;
    padding: 0px 5px;
    padding-bottom: 8px;
    border-bottom: 1px solid #DEDEDE;
}
.data-table th p {
    margin: 0;
    font-size: 12px;
}
.data-table td {
    text-align: left;
    padding: 10px 5px;
    font-size: 15px;
    line-height: 18px;
}
```

## Action Injection

The following will inject the action link (or button) into the e-mail:

```html
{{ with .Email.Body.Actions }}
    {{ if gt (len .) 0 }}
        {{ range $action := . }}
        <p>{{ $action.Instructions }}</p>
        <table class="body-action" align="center" width="100%" cellpadding="0" cellspacing="0">
            <tr>
                <td align="center">
                <div>
                    <a href="{{ $action.Button.Link }}" class="button" style="background-color: {{ $action.Button.Color }}" target="_blank">{{ $action.Button.Text }}</a>
                </div>
                </td>
            </tr>
        </table>
        {{ end }}
    {{ end }}
{{ end }}
```

A good practice is to describe action in footer in case of problem when displaying button and CSS. The text for the description is provided through the `TroubleText` field of the `Product` struct. The text may contain a placeholder `{ACTION}` which is expected to be replaced with the text of the button. The default value of `TroubleText` is `If you’re having trouble with the button '{ACTION}', copy and paste the URL below into your web browser.`

```html
{{ with .Email.Body.Actions }}
<table class="body-sub">
    <tbody><tr>
    {{ range $action := . }}
    <td>
        <p class="sub">{{$.Hermes.Product.TroubleText | replace "{ACTION}" $action.Button.Text}}</p>
        <p class="sub"><a href="{{ $action.Button.Link }}">{{ $action.Button.Link }}</a></p>
    </td>
    {{ end }}
    </tr>
    </tbody>
</table>
{{ end }}
```

## Outro Injection

The following will inject the outro text (string or array) into the e-mail:

```html
{{ with .Email.Body.Outros }} 
    {{ if gt (len .) 0 }}
        {{ range $line := . }}
        <p>{{ $line }}</p>
        {{ end }}
    {{ end }}
{{ end }}
```

## Signature Injection

The following will inject the signature phrase (e.g. Yours truly) along with the product name into the e-mail:

```html
{{.Email.Body.Signature}},
<br>
{{.Hermes.Product.Name}}
```

## Copyright Injection

The following will inject the copyright notice into the e-mail:

```html
{{.Hermes.Product.Copyright}}
```

## Text Direction Injection

In order to support generating RTL e-mails, inject the `textDirection` variable into the `<body>` tag:

```html
<body dir="{{.Hermes.TextDirection}}">
```

## FreeMarkdown Injection

In order to support Markdown free content, inject the following code:

````html
{{ if (ne .Email.Body.FreeMarkdown "") }}
    {{ .Email.Body.FreeMarkdown.ToHTML }}
{{ else }}
    [... Here is the templating for dictionary, table and actions]
{{ end }}
```
