package notification

// MailThemeFlat is a theme
type MailThemeFlat struct{}

// Name returns the name of the flat theme
func (dt *MailThemeFlat) Name() string {
	return "flat"
}

// HTMLTemplate returns a Golang template that will generate an HTML email.
func (dt *MailThemeFlat) HTMLTemplate() string {
	return `
<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">
<html xmlns="http://www.w3.org/1999/xhtml">
<head>
  <meta name="viewport" content="width=device-width, initial-scale=1.0" />
  <meta http-equiv="Content-Type" content="text/html; charset=UTF-8" />
</head>
<body dir="{{.Hermes.TextDirection}}">

  <table class="email-wrapper" width="100%" cellpadding="0" cellspacing="0">
    <tr>
      <td class="content">
        <table class="email-content" width="100%" cellpadding="0" cellspacing="0">
          <!-- Logo -->
          <tr>
            <td class="email-masthead">
              <a class="email-masthead_name" href="{{.Hermes.Product.Link}}" target="_blank">
                {{ if .Hermes.Product.Logo }}
                  <img src="{{.Hermes.Product.Logo | url }}" class="email-logo" style="height: 48px"/>
                  <span style="font-size:36px;font-weight:600;margin-left:12px;color:#409eff">{{ .Hermes.Product.Name}} </span>
                {{ else }}
                  {{ .Hermes.Product.Name }}
                {{ end }}
                </a>
            </td>
          </tr>

          <!-- Email Body -->
          <tr>
            <td class="email-body" width="100%">
              <table class="email-body_inner" align="center" width="570" cellpadding="0" cellspacing="0">
                <!-- Body content -->
                <tr>
                  <td class="content-cell">
                    {{ if (ne .Email.Body.FreeMarkdown "") }}
                      {{ .Email.Body.FreeMarkdown.ToHTML }}
                    {{ else }}

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

                      <!-- Action -->
                      {{ with .Email.Body.Actions }}
                        {{ if gt (len .) 0 }}
                          {{ range $action := . }}
                            <p>{{ $action.Instructions }}</p>
                            <table class="body-action" align="center" width="100%" cellpadding="0" cellspacing="0">
                              <tr>
                                <td align="center">
                                  <div>
                                    <a href="{{ $action.Button.Link }}" class="button" style="background-color: {{ $action.Button.Color }}; color: {{ $action.Button.TextColor }}" target="_blank">
                                      {{ $action.Button.Text }}
                                    </a>
                                  </div>
                                </td>
                              </tr>
                            </table>
                          {{ end }}
                        {{ end }}
                      {{ end }}

                    {{ end }}
                    {{ with .Email.Body.Outros }}
                        {{ if gt (len .) 0 }}
                          {{ range $line := . }}
                            <p>{{ $line }}</p>
                          {{ end }}
                        {{ end }}
                      {{ end }}

                    <p>
                      {{.Email.Body.Signature}}
                    </p>

                    {{ if (eq .Email.Body.FreeMarkdown "") }}
                      {{ with .Email.Body.Actions }}
                        <table class="body-sub">
                          <tbody>
                              {{ range $action := . }}
                                <tr>
                                  <td>
                                    <p class="sub">{{$.Hermes.Product.TroubleText | replace "{ACTION}" $action.Button.Text}}</p>
                                    <p class="sub"><a href="{{ $action.Button.Link }}">{{ $action.Button.Link }}</a></p>
                                  </td>
                                </tr>
                              {{ end }}
                          </tbody>
                        </table>
                      {{ end }}
                    {{ end }}
                  </td>
                </tr>
              </table>
            </td>
          </tr>
          <tr>
            <td>
              <table class="email-footer" align="center" width="570" cellpadding="0" cellspacing="0">
                <tr>
                  <td class="content-cell">
                    <p class="sub center">
                      {{.Hermes.Product.Copyright}}
                    </p>
                  </td>
                </tr>
              </table>
            </td>
          </tr>
        </table>
      </td>
    </tr>
  </table>
</body>
</html>
`
}

// PlainTextTemplate returns a Golang template that will generate an plain text email.
func (dt *MailThemeFlat) PlainTextTemplate() string {
	return `{{ with .Email.Body.Intros }}
  {{ range $line := . }}
    <p>{{ $line }}</p>
  {{ end }}
{{ end }}
{{ if (ne .Email.Body.FreeMarkdown "") }}
  {{ .Email.Body.FreeMarkdown.ToHTML }}
{{ else }}
  {{ with .Email.Body.Dictionary }}
    <ul>
    {{ range $entry := . }}
      <li>{{ $entry.Key }}: {{ $entry.Value }}</li>
    {{ end }}
    </ul>
  {{ end }}
  {{ with .Email.Body.Table }}
    {{ $data := .Data }}
    {{ $columns := .Columns }}
    {{ if gt (len $data) 0 }}
      <table class="data-table" width="100%" cellpadding="0" cellspacing="0">
        <tr>
          {{ $col := index $data 0 }}
          {{ range $entry := $col }}
            <th>{{ $entry.Key }} </th>
          {{ end }}
        </tr>
        {{ range $row := $data }}
          <tr>
            {{ range $cell := $row }}
              <td>
                {{ $cell.Value }}
              </td>
            {{ end }}
          </tr>
        {{ end }}
      </table>
    {{ end }}
  {{ end }}
  {{ with .Email.Body.Actions }}
    {{ range $action := . }}
      <p>{{ $action.Instructions }} {{ $action.Button.Link }}</p>
    {{ end }}
  {{ end }}
{{ end }}
{{ with .Email.Body.Outros }}
  {{ range $line := . }}
    <p>{{ $line }}<p>
  {{ end }}
{{ end }}
<p>{{.Email.Body.Signature}},<br>{{.Hermes.Product.Name}} - {{.Hermes.Product.Link}}</p>

<p>{{.Hermes.Product.Copyright}}</p>
`
}

func (dt *MailThemeFlat) GetStyle() string {
	return `
<style>
.content-cell table {
  width: 100%;
  border-collapse: collapse;
}
.content-cell table,
.content-cell th,
.content-cell td {
  border: 1px solid #EDEFF2;
}
.content-cell th,
.content-cell td {
  padding: 10px;
  font-size: 14px;
  line-height: 18px;
}
.content-cell th {
  background: #409eff;
  color: white;
}
.content-cell td {
  color: #606266;
}
.content-cell p {
  color: #606266;
}
.content-cell a {
  color: #409eff;
}
.email-masthead .email-masthead_name {
  display: flex;
  justify-content: center;
  align-items: center;
  text-decoration: none;
  color: #409eff;
  margin-bottom: 20px;
}
</style>
`
}
