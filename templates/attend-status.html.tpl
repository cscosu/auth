<div id="attend" {{ if .oob }}hx-swap-oob="true" {{ end }}>
  {{ if .canAttend }}
  <h2 class="mb-2 text-xl font-bold">At a meeting right now?</h2>
  <div class="flex gap-2">
    <button
      hx-post="/attend/in-person"
      class="grow justify-center secondary-button"
    >
      I'm here in person
    </button>
    <button
      hx-post="/attend/online"
      class="grow justify-center secondary-button"
    >
      I'm here online
    </button>
  </div>
  {{ else }}
  <p class="text-xl font-bold">Thanks for marking your attendance today!</p>
  {{ end }}
</div>
