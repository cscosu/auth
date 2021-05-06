# This is adapted from SIGPwny's implementation https://github.com/sigpwny/sigpwny-shibboleth-auth (thanks @arxenix)

FROM debian:buster

# -- VARIOUS DEPENDENCIES --
RUN apt-get update

# saml2/shib stuff
RUN apt-get install -y opensaml2-schemas xmltooling-schemas libshibsp8 libshibsp-plugins shibboleth-sp2-common shibboleth-sp2-utils

# nginx stuff
RUN apt-get install -y build-essential libpcre3 libpcre3-dev libssl-dev zlib1g-dev php7.3-fpm

# etc stuff
RUN apt-get install -y wget curl git


# -- BUILD AND INSTALL NGINX W/ nginx-shib MODULE --
RUN echo "deb http://nginx.org/packages/debian/ buster nginx" >> /etc/apt/sources.list
RUN echo "deb-src http://nginx.org/packages/debian/ buster nginx" >> /etc/apt/sources.list

WORKDIR /tmp
RUN wget http://nginx.org/keys/nginx_signing.key
RUN apt-key add nginx_signing.key
RUN apt-get update

WORKDIR /opt
RUN git clone https://github.com/openresty/headers-more-nginx-module.git
RUN git clone https://github.com/nginx-shib/nginx-http-shibboleth.git
RUN mkdir buildnginx && cd buildnginx
RUN cat /etc/apt/sources.list
RUN apt-get source nginx
RUN apt-get build-dep -y nginx

# add the shib module to the build file
WORKDIR /opt/nginx-1.20.0
RUN sed -i 's;--with-ld-opt="$(LDFLAGS)";--with-ld-opt="$(LDFLAGS)" --add-dynamic-module=/opt/headers-more-nginx-module --add-dynamic-module=/opt/nginx-http-shibboleth;g' debian/rules
# build the package
RUN dpkg-buildpackage -b
# install it
WORKDIR /opt
RUN dpkg -i nginx_1.20.0-1~buster_amd64.deb
COPY shib_clear_headers /etc/nginx/shib_clear_headers



# -- INSTALL & CONFIGURE SUPERVISOR --
RUN apt-get install -y supervisor
COPY shibboleth.conf /etc/supervisor/conf.d/shibboleth.conf


# -- KEYS --
COPY ./keys/sp-cert.pem /etc/shibboleth/sp-cert.pem
COPY ./keys/sp-key.pem /etc/shibboleth/sp-key.pem
COPY ./keys/web/auth-osucyber-club-selfsigned.crt /etc/nginx/auth-osucyber-club-selfsigned.crt
COPY ./keys/web/auth-osucyber-club-selfsigned.key /etc/nginx/auth-osucyber-club-selfsigned.key

# -- CONFIGURATION --
COPY ./etc-shib/ /etc/shibboleth/

CMD ["/usr/bin/supervisord"]
# CMD ["service", "supervisor", "start"]
