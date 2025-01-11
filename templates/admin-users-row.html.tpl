{{ if .Editable }}
<tr id="user-{{ .BuckID }}">
  <td>{{ .BuckID }}</td>
  <td><input name="nameNum" value="{{ .NameNum }}" /></td>
  <td>
    <input
      name="discordId"
      value="{{ if not (eq .DiscordID 0) }}{{ .DiscordID }}{{ end }}"
    />
  </td>
  <td><input name="displayName" value="{{ .DisplayName }}" /></td>
  <td>{{ .LastSeenTime }}</td>
  <td>
    {{ if .LastAttendedTime }}
    {{ .LastAttendedTime }}
    {{ end }}
  </td>
  <td>
    <input
      type="checkbox"
      name="addedToMailingList"
      {{
      if
      .AddedToMailingList
      }}checked{{
      end
      }}
    />
  </td>
  <td>
    <input type="checkbox" name="isStudent" {{ if .Student }}checked{{ end }} />
  </td>
  <td>
    <input type="checkbox" name="isAlum" {{ if .Alum }}checked{{ end }} />
  </td>
  <td>
    <input
      type="checkbox"
      name="isEmployee"
      {{
      if
      .Employee
      }}checked{{
      end
      }}
    />
  </td>
  <td>
    <input type="checkbox" name="isFaculty" {{ if .Faculty }}checked{{ end }} />
  </td>
  <td>
    <input type="checkbox" name="isAdmin" {{ if .IsAdmin }}checked{{ end }} />
  </td>
  <td>
    <button
      hx-patch="/admin/users/{{ .BuckID }}"
      hx-swap="outerHTML"
      hx-target="#user-{{ .BuckID }}"
      hx-include="closest tr"
    >
      Save
    </button>
    <button
      hx-get="/admin/users/{{ .BuckID }}"
      hx-vals='{"cancel": true}'
      hx-swap="outerHTML"
      hx-target="#user-{{ .BuckID }}"
    >
      Cancel
    </button>
  </td>
</tr>
{{ else }}
<tr id="user-{{ .BuckID }}">
  <td>{{ .BuckID }}</td>
  <td>{{ .NameNum }}</td>
  <td>{{ if not (eq .DiscordID 0) }}{{ .DiscordID }}{{ end }}</td>
  <td>{{ .DisplayName }}</td>
  <td>{{ .LastSeenTime }}</td>
  <td>
    {{ if .LastAttendedTime }}
    {{ .LastAttendedTime }}
    {{ end }}
  </td>
  <td>
    {{ if .AddedToMailingList }}
    {{ template "checkmark.html.tpl" }}
    {{ else }}
    {{ template "x.html.tpl" }}
    {{ end }}
  </td>
  <td>
    {{ if .Student }}
    {{ template "checkmark.html.tpl" }}
    {{ else }}
    {{ template "x.html.tpl" }}
    {{ end }}
  </td>
  <td>
    {{ if .Alum }}
    {{ template "checkmark.html.tpl" }}
    {{ else }}
    {{ template "x.html.tpl" }}
    {{ end }}
  </td>
  <td>
    {{ if .Employee }}
    {{ template "checkmark.html.tpl" }}
    {{ else }}
    {{ template "x.html.tpl" }}
    {{ end }}
  </td>
  <td>
    {{ if .Faculty }}
    {{ template "checkmark.html.tpl" }}
    {{ else }}
    {{ template "x.html.tpl" }}
    {{ end }}
  </td>
  <td>
    {{ if .IsAdmin }}
    {{ template "checkmark.html.tpl" }}
    {{ else }}
    {{ template "x.html.tpl" }}
    {{ end }}
  </td>
  <td>
    <button
      hx-get="/admin/users/{{ .BuckID }}"
      hx-swap="outerHTML"
      hx-target="#user-{{ .BuckID }}"
    >
      {{ template "pencil.html.tpl" }}
    </button>
  </td>
</tr>
{{ end }}
