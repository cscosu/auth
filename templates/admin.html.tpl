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
    <a
      href="/admin/auth.db"
      target="_top"
      class="rounded-sm px-4 py-2 hover:bg-gray-100 inline-flex justify-between
        {{ if eq .path `/admin/download` }}bg-gray-100{{ end }}"
      >
        <span>Database</span>
        <svg
          xmlns="http://www.w3.org/2000/svg"
          width="24"
          height="24"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="2"
          stroke-linecap="round"
          stroke-linejoin="round"
          class="lucide lucide-download-icon lucide-download"
        >
          <path d="M12 15V3" />
          <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4" />
          <path d="m7 10 5 5 5-5" />
        </svg>
      </a
    >
  </div>
  <div class="flex-1 min-w-0">{{ block "admin" . }}{{ end }}</div>
</div>
{{ end }}
