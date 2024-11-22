{{ template "base.html.tpl" . }}

{{ define "title" }}Not Found | Cyber Security Club @ Ohio State{{ end }}

{{ define "content" }}
<div class="h-full flex flex-col justify-center items-center">
  <div class="text-4xl font-bold">404</div>
  <div class="text-lg">Page not found</div>
</div>
{{ end }}
