# Auth

This is an auth platform and discord bot that interacts with SAML as a single binary, without Shibboleth or Nginx.

## Development

By default, the SAML provider is mocked out, so it is possible to develop without Service Provider secrets. Just run:

```
docker compose up --watch --build --remove-orphans
```

This will start a local server at http://localhost:3000

You will also want to install `gopls` for the language server in your IDE.

## Development with SAML

The certificate for our current credentials will expire `May 4 20:08:05 2031 GMT`, at which point it will need to be renewed.

To develop with SAML credentials, put the Service Provider keys into `keys/sp-cert.pem` and `keys/sp-key.pem`. Then run:

```
docker compose -f docker-compose-saml.yaml up --watch --build --remove-orphans
```

This will start a local server at https://test-auth.osucyber.club. It will generate a self signed certificate for local https.

## How does it work?

Read this for specific details on OSU's SSO: https://webauth.service.ohio-state.edu/~shibboleth/index.html

In summary:

The OSU authentication system uses SAML2, like many other universities. This means that OSU runs an Identity Provider (IdP), which stores user data (like `Name.#`, email, and `BuckID`). Third parties can interact with the IdP by becoming a Service Provider (SP), think Schedule Planner. The OSU Cyber Security Club also has been granted private keys to be an SP at 2 hosts: `https://auth-test.osucyber.club` and `https://auth.osucyber.club`. We've configured DNS to point `auth-test.osucyber.club` to `127.0.0.1` for local testing, and it still requires self signed certificates. `https://auth.osucyber.club` is the public, production web server.

Here are some other important facts about our setup in particular:
- Our registered entity ID is `https://auth.osucyber.club/shibboleth`
- Shibboleth assertions will be sent to our `/Shibboleth.sso/` and below (for example `/Shibboleth.sso/SAML2/POST`)

The recommended way to actually use those keys as an SP is using [Shibboleth](https://shibboleth.net). Shibboleth was primarily developed in 2004 by an OSU employee and integrates nicely into Apache and IIS. But, it is kind of a nightmare. OSUCyber's former auth system, modelled after [sigpwny's auth system](https://github.com/sigpwny/sigpwny-shibboleth-auth) dockerized the process, which works this way:

There are 2 containers side-by-side:
- Shibboleth
  - Shibboleth's 3 processes, `shibd`, `shibauthorizer`, and `shibresponder`
    - Configured with `OSU-attribute-policy.xml`, `OSU-metadata.cer`, `shibboleth2.xml`, `attribute-map.xml`, `sessionError.html`
  - Nginx
    - Using [a plugin to interact with Shibboleth](https://github.com/nginx-shib/nginx-http-shibboleth)
    - Configured `/login` endpoint to be restricted by Shibboleth authentication, using redirects to sign in. Then once it is successful, it sets trusted HTTP headers like `Employeenumber` and `Displayname`. This is how the webapp will receive information about the OSU user.
- Webapp
  - Whatever webapp you want to write
  - Has a `/login` which reads from the trusted HTTP headers set by Nginx. Then the app can put it in its own database, or however it wants to handle it.

However, updating any of the Shibboleth/Nginx stuff is really scary. Instead, it is possible to just use a SAML2 library to be an SP, all in one single place. That is what this repository is. We use [`crewjam/saml`](https://github.com/crewjam/saml) to handle being an SP in golang, and then build the rest of the auth app around it. It requires picking out only the important parts of the `shibboleth2.xml`, and makes significantly easier to read.
