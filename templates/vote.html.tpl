{{ template "base.html.tpl" . }}

{{ define "title" }}Vote | Cyber Security Club @ Ohio State{{ end }}

{{ define "content" }}

<div id="voting-form">
  {{ if .hasVoted }}
  <div class="card">
    <div class="card-title">Thanks for voting!</div>
    <div class="card-content">
      Your vote in the <strong>{{ .electionName }}</strong> election has been
      recorded.
    </div>
  </div>
  {{ else if not .electionName }}
  <div class="card">
    <div class="card-title">No active vote</div>
    <div class="card-content">
      There is no active election right now. Try again later.
    </div>
  </div>
  {{ else }}
  {{ template "voting-form.html.tpl" . }}
  {{ end }}
</div>

{{ end }}
