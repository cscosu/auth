{{ template "base.html.tpl" . }}

{{ define "title" }}OSU Cyber Security Club{{ end }}

{{ define "content" }}
<div class="card">
  <div class="card-title">Attendance History</div>
  <div class="card-content">
    <table class="table-auto">
      <thead>
        <tr class="border-b-2">
          <th class="px-4 py-2 text-left">Date</th>
          <th class="px-4 py-2 text-left">Type</th>
        </tr>
      </thead>
      <tbody class="[&_td]:px-4 [&_td]:py-2 [&_tr:not(:last-child)]:border-b">
        {{
          range.records
        }}
        <tr>
          <td>{{ .timestamp }}</td>
          <td>{{ .type }}</td>
        </tr>
        {{
          end
        }}
      </tbody>
    </table>
  </div>
</div>
{{ end }}
