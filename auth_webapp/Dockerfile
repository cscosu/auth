FROM python:3
ENV PYTHONUNBUFFERED=1
ENV DJANGO_DEBUG=False
WORKDIR /code
COPY requirements.txt /code/
RUN pip install -r requirements.txt
COPY . /code/
RUN python3 manage.py check --deploy
EXPOSE 8000
CMD python manage.py runserver 0.0.0.0:8000

