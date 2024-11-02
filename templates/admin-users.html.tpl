{{ template "admin.html.tpl" . }}

{{ define "title" }}Admin Users | Cyber Security Club @ Ohio State{{ end }}

{{ define "admin" }}
<div class="card">
  <div class="card-title">Users</div>
  <div class="card-content">
    <table class="table-auto">
      <thead>
        <tr class="border-b-2">
          <th class="px-4 py-2 text-left">
            <a
              target="_blank"
              href="https://webauth.service.ohio-state.edu/~shibboleth/user-attribute-reference.html?article=mail"
              >Name.#</a
            >
          </th>
          <th class="px-4 py-2 text-left">Discord ID</th>
          <th class="px-4 py-2 text-left">
            <a
              target="_blank"
              href="https://webauth.service.ohio-state.edu/~shibboleth/user-attribute-reference.html?article=employeenumber"
              >Buck ID</a
            >
          </th>
          <th class="px-4 py-2 text-left">
            <a
              target="_blank"
              href="https://webauth.service.ohio-state.edu/~shibboleth/user-attribute-reference.html?article=preferred-names"
              >Name</a
            >
          </th>
          <th class="px-4 py-2 text-left">
            <a
              class="inline-flex"
              target="_blank"
              href="https://webauth.service.ohio-state.edu/~shibboleth/user-attribute-reference.html?article=idm-id"
              >IDM ID {{ template "key.html.tpl" }}</a
            >
          </th>
        </tr>
      </thead>
      <tbody class="[&_td]:px-4 [&_td]:py-2 [&_tr:not(:last-child)]:border-b">
        {{
          range.users
        }}
        <tr>
          <td>{{ .NameNum }}</td>
          <td>{{ if not (eq .DiscordID 0) }}{{ .DiscordID }}{{ end }}</td>
          <td>{{ .BuckID }}</td>
          <td>{{ .DisplayName }}</td>
          <td>{{ .IDMID }}</td>
        </tr>
        {{
          end
        }}
      </tbody>
    </table>
  </div>
</div>
{{ end }}
