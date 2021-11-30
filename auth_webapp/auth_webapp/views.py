from django.contrib.auth.models import AnonymousUser
from django.http import HttpResponse
from django.utils.http import url_has_allowed_host_and_scheme
from django.utils.encoding import iri_to_uri
from django.http.response import HttpResponseBadRequest, HttpResponseForbidden, HttpResponseServerError
from django.shortcuts import render, redirect
from django.conf import settings
from .models import OSUUser, AttendanceRecord
from django.contrib.auth import authenticate, login, logout
from django.contrib import messages
from django.views.decorators.http import require_http_methods
import datetime
import time
import pytz
import jwt

ohio_tz = pytz.timezone('America/New_York')


def home(request):
    template_data = {'session': request.session}
    return render(request, 'home.html', template_data)


def new_member(request):
    template_data = {'session': request.session}
    if request.user.is_authenticated:
        discord_token = jwt.encode({"buck_id": request.user.shib_id, "date": str(
            datetime.date.today())}, settings.JWT_SECRET2, algorithm='HS256')
        template_data['discord_token_msg'] = f"!connect {discord_token}"
        template_data['can_subscribe'] = not request.user.is_currently_on_mailing_list()
    return render(request, 'new_members.html', template_data)


@require_http_methods(["GET"])
def login_view(request):
    # In production, request.headers for /login ONLY are TRUSTED -- nginx is the only
    # one that can set them, and will set them only on successful shib

    # https://webauth.service.ohio-state.edu/~shibboleth/user-attribute-reference.html#edupersonscopedaffiliation
    if 'Employeenumber' not in request.headers or 'Displayname' not in request.headers or 'Eppn' not in request.headers:
        return HttpResponseServerError("Sorry, something went wrong. Contact info@osucyber.club. :(")
    name_n = request.headers['Eppn'].split("@")[0]
    shib_id = request.headers['Employeenumber']

    # Pick the best affiliation for this user
    # I'm enforcing a weird heirarchy here
    best_affiliation = OSUUser.aff_choices.NONE
    print(request.headers['Affiliation'])
    for a in request.headers['Affiliation'].split(';'):
        if a == "student@osu.edu" and best_affiliation != OSUUser.aff_choices.FACULTY_STAFF:
            best_affiliation = OSUUser.aff_choices.STUDENT
        elif a == "faculty@osu.edu":
            best_affiliation = OSUUser.aff_choices.FACULTY_STAFF
        elif a == "alum@osu.edu" and best_affiliation == OSUUser.aff_choices.NONE:
            best_affiliation = OSUUser.aff_choices.ALUMNI

    user, created = OSUUser.objects.get_or_create(shib_id=shib_id, defaults={
                                                  "display_name": request.headers['Displayname'], "name_num": name_n})

    # Users can be updated in shib, we should record those updates
    if best_affiliation != user.affiliation or request.headers['Displayname'] != user.display_name:
        user.affiliation = best_affiliation
        user.display_name = request.headers['Displayname']
        user.save()
    login(request, user)

    if 'next' in request.GET and url_has_allowed_host_and_scheme(request.GET['next'], None):
        url = iri_to_uri(request.GET['next'])
        return redirect(url)
    else:
        return redirect("new_member" if created else "home")


@require_http_methods(["GET"])
def debug_login(request):
    if not settings.DEBUG:
        return HttpResponseForbidden()
    if request.user.is_authenticated:
        return HttpResponseBadRequest("You are already logged in")

    if 'id' not in request.GET:
        return HttpResponseBadRequest("Please include id in the query string :)")
    # Lets create a user
    user, created = OSUUser.objects.get_or_create(shib_id=request.GET['id'], defaults={
                                                  "display_name": "test user %s" % request.GET['id'], "name_num": "test.%s" % request.GET['id']})
    if 'super' in request.GET and request.GET['super'] == '1':
        user.is_superuser = True
        user.save()
    login(request, user)
    return redirect("home")


def logout_view(request):
    logout(request)
    return redirect("home")


def _user_can_submit_attendance(user):
    try:
        last_attend = AttendanceRecord.objects.filter(
            user=user).order_by('-date_recorded')[0]
        return datetime.datetime.now(tz=ohio_tz) > last_attend.date_recorded + datetime.timedelta(hours=2)
    except IndexError:
        return True


@require_http_methods(["GET"])
def attendance(request):
    if not request.user.is_authenticated:
        return redirect("/login?return=/attend")
    return render(request, 'attendance.html', {
        'can_attend': _user_can_submit_attendance(request.user),
        'can_subscribe': not request.user.is_currently_on_mailing_list(),
        'attendance': AttendanceRecord.objects.filter(user=request.user).order_by('-date_recorded')
    })


@require_http_methods(["POST"])
def attend(request):
    if _user_can_submit_attendance(request.user):
        ar = AttendanceRecord(
            user=request.user, date_recorded=datetime.datetime.now().astimezone(ohio_tz))
        ar.save()
    return redirect("attendance")


@require_http_methods(["POST"])
def subscribe(request):
    if not request.user.is_authenticated:
        return redirect("/login?return=/new")

    if request.user.is_currently_on_mailing_list():
        return redirect('new_member')

    if request.user.add_to_mailing_list():
        messages.add_message(request, messages.INFO,
                             'You have been added to our mailing list!')
    else:
        messages.add_message(
            request, messages.ERROR, 'Sorry, something went wrong adding you to the mailing list. Try http://mailinglist.osucyber.club instead, and let us know.')

    return redirect('new_member')
