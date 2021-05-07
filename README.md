# auth.osucyber.club

First of all, credit goes to SIGPwny (UIUC's Cybersecurity Club) for their [dockerized Shibboleth SP setup](https://github.com/sigpwny/sigpwny-shibboleth-auth). The shibboleth part of this is only lightly modified.

Features:

- Login with OSU
- Assign discord role based on affiliation (Student, Faculty/Staff, Alumni) w/ our [discord bot](https://github.com/cscosu/discord_bot)
- One-click add to mailing list

Planned:

- Ability to download private files (stored on S3?) for students


## How to run just web app

```
cd auth_webapp
docker-compose up
```

Navigate to http://localhost:8000 to view the app
Use the debug login feature: http://localhost:8000/debug_login/?id=my_cool_id&super=1

Then you can check out django's automated admin pages at http://localhost:8000/admin/ and
you can change attributes like the user's affiliation, to test stuff.

The debug login feature logs in the same way a real user would log in. request.user will
have all the same attributes.

## How to run the whole thing (including local shib) -- NOT RECOMMENDED

This won't work unless you have the shib key (which we can't really give out).

First you need to generate some self-signed keys for debug purposes:

```
mkdir -p shib_docker/keys/web/
openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout shib_docker/keys/web/auth-osucyber-club-selfsigned.key -out shib_docker/keys/web/auth-osucyber-club-selfsigned.crt
```

Then you should generate the static files folder:
```
mkdir shib_docker/static/
cd auth_webapp
python3 manage.py collectstatic
```

todo: write more information about how this works