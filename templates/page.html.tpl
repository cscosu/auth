{{ template "base.html.tpl" . }}
{{ define "body" }}
<div>{{ block "content" . }}{{ end }}</div>
{{ end }}
