{{ template "base.html.tpl" . }}

{{ define "title" }}OSU Cyber Security Club{{ end }}

{{ define "content" }}
<div class="grid lg:grid-cols-2 gap-6">
  <div class="card">
    <div class="card-title">Contact</div>
    <div class="card-content flex flex-col gap-2">
      <a href="/discord/signin" class="secondary-button flex flex-row gap-2">
        {{ if .hasLinkedDiscord }}
        {{ template "checkmark.html.tpl" }}
        {{ end }}
        Link Discord
      </a>
      <a class="secondary-button flex flex-row gap-2">
        Subscribe to our mailing list
      </a>
    </div>
  </div>
  <div class="card">
    <div class="card-title">Meetings</div>
    <div class="card-content">
      <p>We meet Tuesdays at 7pm in Enarson 230.</p>
      <p>Meeting reminders are sent as part of a weekly newsletter.</p>
      <h2 class="mt-4 mb-2 text-xl font-bold">At a meeting right now?</h2>
      <div class="flex gap-2">
        <a class="grow justify-center secondary-button">I'm here in person</a>
        <a class="grow justify-center secondary-button">I'm here online</a>
      </div>
      <a href="/attendance" class="external-link text-sm italic text-gray-600"
        >View past attendance</a
      >
    </div>
  </div>
  <div class="card">
    <div class="card-title">Resources</div>
    <div class="card-content">
      <ul class="list-disc list-inside">
        <li>
          Check out our
          <a class="external-link" href="https://wiki.osucyber.club">Wiki</a>!
          We have a lot of content from past meetings available.
        </li>
        <li>
          Play in our
          <a
            class="external-link"
            href="https://wiki.osucyber.club/Bootcamp-CTF/Welcome"
            >24/7 Bootcamp CTF</a
          >! This is our series of hacking challenges you can try, right now!
        </li>
      </ul>
    </div>
  </div>
</div>
{{ end }}
