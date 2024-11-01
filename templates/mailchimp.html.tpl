{{ if .isOnMailingList }}
<div
  class="secondary-button hover:border-gray-300 hover:cursor-default flex flex-row gap-2"
>
  {{ template "checkmark.html.tpl" }}
  <span class="cursor-text">Subscribed to our mailing list</span>
</div>
{{ else }}
<button
  hx-post="/mailchimp"
  hx-swap="outerHTML"
  class="secondary-button flex flex-row gap-2"
>
  {{ template "email.html.tpl" }}
  Subscribe to our mailing list
</button>
{{ end }}
