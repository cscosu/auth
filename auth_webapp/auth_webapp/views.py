from django.http import HttpResponse

def home(request):
    if request.session is not None and 'headers' in request.session:
        return HttpResponse(str(request.session['headers']))
    else:
        return HttpResponse("not logged in")

def login(request):
    request.session['headers'] = str(request.headers)
    return HttpResponse("you're now logged in!\n"+str(request.headers))