{{ template "base.html.tpl" . }}

{{ define "title" }}Cyber Security Club @ Ohio State{{ end }}

{{ define "content" }}
<div class="grid lg:grid-cols-2 gap-6">
  <div class="card">
    <div class="card-title">Contact</div>
    <div class="card-content flex flex-col gap-2">
      <a href="/discord/signin" class="secondary-button flex flex-row gap-2">
        {{ if .hasLinkedDiscord }}
        {{ template "checkmark.html.tpl" }}
        {{ else }}
        {{ template "discord.html.tpl" }}
        {{ end }}
        Join/Link Discord
      </a>
      {{ template "mailchimp.html.tpl" . }}
      {{ if .isOnMailingList }}
      <div>
        If you removed yourself from the mailing list and need to resubscribe,
        click
        <a
          class="external-link"
          href="https://mailinglist.osucyber.club"
          target="_blank"
          >here</a
        >.
      </div>
      {{ end }}
      <div>
        Issues? Ask in the Discord or email
        <a class="external-link" href="mailto:info@osucyber.club"
          >info@osucyber.club</a
        >.
      </div>
    </div>
  </div>
  <div class="card">
    <div class="card-title">Meetings</div>
    <div class="card-content">
      <p class="mb-2">
        We meet Tuesdays at 7pm in Enarson 230. Meeting reminders are sent as
        part of a weekly newsletter.
      </p>

      <a
        hx-boost="true"
        href="/attendance"
        class="external-link text-sm italic text-gray-600"
        >View past attendance</a
      >
      {{ template "attend-status.html.tpl" . }}
    </div>
  </div>
  <div class="card">
    <div class="card-title">Resources</div>
    <div class="card-content">
      <ul class="list-disc list-inside">
        <li>
          Check out our
          <a
            class="external-link"
            href="https://wiki.osucyber.club"
            target="_blank"
            >Wiki</a
          >! We have a lot of content from past meetings available.
        </li>
        <li>
          Play in our
          <a
            class="external-link"
            href="https://wiki.osucyber.club/Bootcamp-CTF/Welcome"
            target="_blank"
            >24/7 Bootcamp CTF</a
          >! This is our series of hacking challenges you can try, right now!
        </li>
      </ul>
    </div>
  </div>
</div>
{{ end }}
