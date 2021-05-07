FROM tiangolo/meinheld-gunicorn:python3.7

ENV MODULE_NAME=auth_webapp.wsgi
ENV VARIABLE_NAME=application
ENV PORT=8000

COPY . /app/
WORKDIR /app/
RUN pip install -r requirements.txt
RUN python3 manage.py check --deploy