<div class="card">
  {{ if not .hasVoted }}
  <div class="card-title">Cast your vote for {{ .electionName }}</div>
  <div class="card-content">
    <form>
      <ul>
        {{
          range.candidates
        }}
        <li>
          <input type="radio" id="{{ .Id }}" name="vote" value="{{ .Id }}" />
          <label for="{{ .Id }}"> {{ .Name }}</label>
        </li>
        {{
          end
        }}
      </ul>
      <button
        hx-post="/vote"
        hx-target="#voting-form"
        class="grow justify-center secondary-button"
      >
        Submit
      </button>
    </form>
  </div>
  {{ else }}
  <p class="card-title">Thanks for voting!</p>
  {{ end }}
</div>
