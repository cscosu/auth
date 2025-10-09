{{ template "admin.html.tpl" . }}

{{ define "title" }}Admin Vote | Cyber Security Club @ Ohio State{{ end }}

{{ define "admin" }}
<div class="card">
  <div class="card-title">Admin Votes</div>
	<a hx-boost="true" href="/admin/vote/new" class="secondary-button"
		>New Vote</a
	>
	<div class="container">
		<div class="subcard">
			<h5 class="mb-3 text-base font-semibold">
				Published
			</h5>
			<ul>
				{{ range.pastElections }}
				<li>
					{{ if .Published }} 
					<span class="flex-1 ms-3 whitespace-nowrap">
						<a
							hx-boost="true"
							href="/admin/vote/{{ .ElectionId }}"
							>{{ .Name }}</a
						>
					</span>
					{{ end }}
				</li>
				{{ end }}
			</ul>
		</div>
		<div class="subcard">
			<h5 class="mb-3 text-base font-semibold">
				Closed
			</h5>
			<ul>
				{{ range.pastElections }}
				<li>
					{{ if .Done }} 
					<span class="flex-1 ms-3 whitespace-nowrap">
						<a
							hx-boost="true"
							href="/admin/vote/{{ .ElectionId }}"
							>{{ .Name }}</a
						>
					</span>
					{{ end }}
				</li>
				{{ end }}
			</ul>
		</div>
		<div class="subcard">
			<h5 class="mb-3 text-base font-semibold">
				Unpublished
			</h5>
			<ul>
				{{ range.pastElections }}
				<li>
					{{ if and (not .Published) (not .Done) }} 
					<span class="flex-1 ms-3 whitespace-nowrap">
						<a
							hx-boost="true"
							href="/admin/vote/{{ .ElectionId }}"
							>{{ .Name }}</a
						>
					</span>
					{{ end }}
				</li>
				{{ end }}
			</ul>
		</div>
	</div>
</div>
{{ end }}
