{{ template "base.html.tpl" . }}

{{ define "content" }}
<div class="flex gap-4 flex-col md:flex-row">
  <div class="flex flex-col gap-2 w-40" hx-boost="true">
    <a
      href="/admin/users"
      class="rounded-sm px-4 py-2 hover:bg-gray-100
        {{ if eq .path `/admin/users` }}bg-gray-100{{ end }}"
      >Users</a
    >
    <a
      href="/admin/vote"
      class="rounded-sm px-4 py-2 hover:bg-gray-100
        {{ if eq .path `/admin/vote` }}bg-gray-100{{ end }}"
      >Vote</a
    >
  </div>
  <div class="flex-1 min-w-0">{{ block "admin" . }}{{ end }}</div>
</div>
{{ end }}
