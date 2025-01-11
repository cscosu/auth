{{ if .published }}
<div id="form">
  <p class="font-bold">Election Name</p>
  <p>{{ .electionName }}</p>
  <p class="font-bold">Candidates</p>
  <ul class="flex flex-col gap-1 list-disc list-inside">
    {{
      range.candidates
    }}
    <li>{{ .Name }}{{ if $.done }}: {{ .Votes }}{{ end }}</li>
    {{
      end
    }}
  </ul>
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
  <input
    type="text"
    name="electionName"
    value="{{ .electionName }}"
    hx-trigger="input changed delay:500ms"
    hx-patch="/admin/vote/{{ .electionId }}"
  />
  <p class="font-bold">Candidates</p>
  <ul>
    {{
      range.candidates
    }}
    <li class="flex items-center before:content-['•'] before:mr-2 group">
      <input
        type="text"
        name="candidateName"
        value="{{ .Name }}"
        hx-trigger="input changed delay:500ms"
        hx-patch="/admin/vote/{{ $.electionId }}/{{ .Id }}"
      />
      <button
        class="ml-1 text-red-700 hidden group-hover:block"
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
    {{
      end
    }}
    <li>
      <button
        class="italic"
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
