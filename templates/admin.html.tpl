{{ template "base.html.tpl" . }}

{{ define "title" }}Cyber Security Club @ Ohio State{{ end }}

{{ define "content" }}
<div class="grid lg:grid-cols-2 gap-6">
  <div class="card">
    <div class="card-title">Admin</div>
    <div class="card-content">Content</div>
  </div>
</div>
{{ end }}
