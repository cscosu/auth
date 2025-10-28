{{ if .published }}
<div id="form">
  <p class="font-bold">Election Name</p>
  <p>{{ .electionName }}</p>
  <p class="font-bold">Candidates</p>
	<div class="relative overflow-x-auto">
			<table class="w-full text-sm text-left rtl:text-right">
					<thead class="dark:bg-gray-400">
							<tr>
									<th class="px-6 py-3">
										Name
									</th>
									<th class="px-6 py-3">
										Vote Count
									</th>
									<th class="px-6 py-3">
										Percentage
									</th>
							</tr>
					</thead>
					<tbody>
							{{
								range .candidates 
							}}
							<tr class="bg-white border-b dark:bg-gray-200 dark:border-gray-300 border-gray-200">
									<th scope="row" class="px-6 py-4 font-medium text-gray-900 whitespace-nowrap">
										{{ .Name }}
									</th>
									<td class="px-6 py-4">
										{{ .Votes }}
									</td>
									<td class="px-6 py-4">
										{{ printf "%.2f" .Percentage }}%
									</td>
							</tr>
							{{
								end
							}}
					</tbody>
			</table>
	</div>
  <p>{{ .totalVotes }} total vote{{ if ne .totalVotes 1 }}s{{ end }}</p>
  {{ if not .done }}
  <a
    hx-boost="true"
    href="/admin/vote/{{ .electionId }}/close"
    class="secondary-button"
  >
    Close vote
  </a>
  {{ end }}
</div>
{{ else }}
<div id="form">
  <p class="font-bold">Election Name</p>
	<input type="text" name="electionName" value="{{ .electionName }}" hx-trigger="input changed delay:500ms" hx-patch="/admin/vote/{{ .electionId }}" class="bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500" placeholder="New Name" required />
  <p class="font-bold">Candidates</p>
  <ul>
		{{ range.candidates }}
    <li class="flex items-center before:mr-2 group">
			<input type="text" name="candidateName" value="{{ .Name }}" hx-trigger="input changed delay:500ms" hx-patch="/admin/vote/{{ $.electionId }}/{{ .Id }}" class="bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500" placeholder="New Name" required />
      <button
        class="ml-1 text-red-700 group-hover:block"
        hx-delete="/admin/vote/{{ $.electionId }}/{{ .Id }}"
        hx-target="#form"
      >
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
          class="lucide lucide-trash-2"
        >
          <path d="M3 6h18" />
          <path d="M19 6v14c0 1-1 2-2 2H7c-1 0-2-1-2-2V6" />
          <path d="M8 6V4c0-1 1-2 2-2h4c1 0 2 1 2 2v2" />
          <line x1="10" x2="10" y1="11" y2="17" />
          <line x1="14" x2="14" y1="11" y2="17" />
        </svg>
      </button>
    </li>
		{{ end }}
    <li>
      <button
				class="secondary-button"
        hx-put="/admin/vote/{{ .electionId }}"
        hx-target="#form"
      >
        Add New Candidate
      </button>
    </li>
  </ul>
  <button
    hx-post="/admin/vote/{{ .electionId }}/publish"
    hx-target="#form"
    class="secondary-button"
  >
    Publish vote
  </button>
</div>
{{ end }}
