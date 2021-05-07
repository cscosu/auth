# Integration with discord bot

from django.contrib.auth.models import AnonymousUser
from django.http import HttpResponse, JsonResponse
from django.utils.http import url_has_allowed_host_and_scheme
from django.utils.encoding import iri_to_uri
from django.http.response import HttpResponseBadRequest, HttpResponseForbidden, HttpResponseServerError
from django.shortcuts import render, redirect
from django.conf import settings
from .models import OSUUser, AttendanceRecord
from django.views.decorators.http import require_http_methods
from django.views.decorators.csrf import csrf_exempt
import jwt

@require_http_methods('POST')
@csrf_exempt
def link_discord(request):
    if 'token' not in request.POST:
        print("missing token")
        return HttpResponseBadRequest("bad token")
    
    try:
        token = jwt.decode(request.POST['token'], settings.JWT_SECRET1, algorithms='HS256', options={"require": ["discord_id", "auth_token"]})
        user_token = jwt.decode(token['auth_token'], settings.JWT_SECRET2, algorithms='HS256', options={"require": ["buck_id"]})
        
        print(f"Processing token for {token['discord_id']}")
        try:
            discord_linked = OSUUser.objects.get(discord_id=token['discord_id'])
            print(f"ID is already linked to user {discord_linked.shib_id}")
            return JsonResponse({"success": False, "msg": "Your discord is already linked to another OSU user. Please contact a club officer and we can fix it up for you."})
        except OSUUser.DoesNotExist:
            pass

        user = OSUUser.objects.get(shib_id=user_token["buck_id"])
        if user.discord_id is None:
            user.discord_id = token["discord_id"]
            user.save()
            print(f"Successfully linked {user.discord_id} to shib {user.shib_id}")
            return JsonResponse({"success": True, "affiliation": user.affiliation})
        else:
            print(f"User {user.shib_id} is already linked to another discord user")
            return JsonResponse({"success": False, "msg": "You already have another discord account linked. Please contact a club officer and we can fix it up for you."})

    except Exception as e: # This catches jwt decode errors, or missing user errors
        print(e)
        return HttpResponseBadRequest("bad token")
