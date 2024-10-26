{{ template "base.html.tpl" . }}

{{ define "title" }}Cyber Security Club @ Ohio State{{ end }}

{{ define "content" }}
<div class="card">
  <div class="card-title">Sign in to use this site</div>
  <div class="card-content">
    <p>This site will allow you to:</p>

    <ul class="list-disc list-inside">
      <li>Link your Discord account (access non-public channels)</li>
      <li>Download members-only files</li>
      <li>Submit attendance at meetings</li>
      <li>Join the mailing list (if you aren't already on it)</li>
    </ul>

    <p>
      OSU Students, faculty, and alumni can click the button below to sign in
    </p>
    <a href="/signin" class="secondary-button">Sign in with OSU</a>
  </div>
</div>
{{ end }}
