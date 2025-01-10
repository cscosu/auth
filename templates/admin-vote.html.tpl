{{ template "admin.html.tpl" . }}

{{ define "title" }}Admin Vote | Cyber Security Club @ Ohio State{{ end }}

{{ define "admin" }}
<div class="card">
  <div class="card-title">Admin Votes</div>
  <div class="card-content">
    <a hx-boost="true" href="/admin/vote/new" class="secondary-button"
      >New Vote</a
    >

    <div>
      <h1>Past Votes</h1>
      <ul class="list-disc list-inside">
        {{
          range.pastElections
        }}
        <li>
          <a
            class="font-bold"
            hx-boost="true"
            href="/admin/vote/{{ .ElectionId }}"
            >{{ .Name }}</a
          >
          {{ if .Published }}<em>published</em>{{ end }}

          {{ if .Done }}
          <em>completed {{ .DoneTime }}</em>
          {{ end }}
        </li>
        {{
          end
        }}
      </ul>
    </div>
  </div>
</div>
{{ end }}
