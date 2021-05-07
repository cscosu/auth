from django.http import HttpResponse
from django.shortcuts import render, redirect

def home(request):
    return render(request, 'home.html', {'session': request.session})

def login(request):
    request.session['headers'] = str(request.headers)
    return redirect("/")