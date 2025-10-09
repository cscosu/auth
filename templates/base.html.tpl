<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    <meta name="title" content="{{ block `title` . }}{{ end }}" />
    <script src="https://unpkg.com/htmx.org@2.0.3"></script>
    <script src="https://cdn.jsdelivr.net/npm/chart.js"></script>
    <script src="https://cdn.tailwindcss.com"></script>
    <style type="text/tailwindcss">
      @tailwind base;
      @tailwind components;
      @tailwind utilities;
      .card { 
				@apply shadow-md rounded-sm 
				border; 
				padding-right: 25px;
			}
			.container {
				display: flex; 
				margin-right: 20px;
			}
			.container div { margin: 5px 10px; }
      .subcard { @apply w-full max-w-sm p-4 bg-white border border-gray-200 rounded-lg shadow-sm sm:p-6 dark:bg-gray-100 dark:border-gray-100; }
      .card-title { @apply shadow-sm py-2 px-3 font-bold; }
      .card-content { @apply py-2 px-3; }
      .external-link { @apply text-blue-500; }
      .secondary-button { @apply inline-flex border border-gray-300 hover:border-gray-400 active:border-gray-500 text-gray-800 rounded-md cursor-pointer px-4 py-2 text-center; margin-left: 20px; }
      .primary-button { @apply inline-flex bg-teal-400 hover:bg-teal-500 active:bg-teal-600 font-bold text-gray-100 rounded-md cursor-pointer px-4 py-2 text-center; }
    </style>
    <title>{{ block "title" . }}{{ end }}</title>
  </head>
  <body class="min-h-screen flex flex-col">
    <header class="border-b flex justify-center py-2">
      <div class="container mx-2">
        <nav class="flex items-center justify-between">
          <div class="flex items-center">
            <a
              hx-boost="true"
              href="/"
              class="mr-6 flex items-center space-x-2"
            >
              <img src="/static/logo.png" width="28" />
            </a>
            <a hx-boost="true" href="/">OSU Cyber Security Club Auth</a>
          </div>
          {{ if .nameNum }}
          <div>
            <span class="mr-4">Signed in as {{ .nameNum }}</span>
            <a href="/signout" class="primary-button">Sign out</a>
          </div>
          {{ else }}
          <a href="/signin" class="primary-button">Sign in with OSU</a>
          {{ end }}
        </nav>
      </div>
    </header>
    <div class="flex flex-1 justify-center">
      <div class="container mx-2 mt-6">{{ block "content" . }}{{ end }}</div>
    </div>
    <footer class="flex justify-center">
      <div class="container py-6 mx-2">
        <a href="https://osucyber.club"
          >Cyber Security Club @ The Ohio State University</a
        >
      </div>
    </footer>
  </body>
</html>
