{{ template "admin.html.tpl" . }}

{{ define "title" }}Admin Vote | Cyber Security Club @ Ohio State{{ end }}

{{ define "admin" }}
<div class="card">
  <div class="card-title">Admin Vote Edit</div>
  <div class="card-content flex justify-between">
    <div>
      {{ template "admin-vote-edit-partial.html.tpl" . }}
      <a
        hx-boost="true"
        hx-confirm="Are you sure you want to permanently delete this vote and all associated data?"
        href="/admin/vote/{{ .electionId }}/delete"
        class="secondary-button"
      >
        Delete
      </a>
    </div>
    {{ if .done }}
    <div>
      <canvas id="candidatesChart"></canvas>
      <div id="labels" class="hidden">
        {{ range.candidates }}
        <p>{{ .Name }}</p>
        {{ end }}
      </div>
      <div id="datas" class="hidden">
        {{ range.candidates }}
        <p>{{ .Votes }}</p>
        {{ end }}
      </div>
      <script>
        (function () {
          const candidatesChart = document.getElementById("candidatesChart");

          const labels = [...document.getElementById("labels").children].map(
            (child) => child.innerHTML
          );
          const datas = [...document.getElementById("datas").children].map(
            (child) => parseInt(child.innerHTML)
          );

          const data = {
            labels: labels,
            datasets: [
              {
                label: "Votes",
                data: datas,
              },
            ],
          };

          new Chart(candidatesChart, {
            type: "pie",
            data: data,
          });
        })();
      </script>
    </div>
    {{ end }}
  </div>
</div>
{{ end }}
