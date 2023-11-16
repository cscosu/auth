from django.http import JsonResponse
from django.conf import settings
from django.views.decorators.http import require_http_methods
from .models import OSUUser, AttendanceRecord
from datetime import datetime
from django.views.decorators.csrf import csrf_exempt
import json


def auth(view):
    def new_view(request, *args, **kwargs):
        header = request.headers.get("Authorization")
        if header:
            if header != f"Bearer {settings.API_TOKEN}":
                return JsonResponse({"error": "invalid api token"}, status=403)
            return view(request, *args, **kwargs)
        else:
            return JsonResponse({"error": "missing api token"}, status=403)

    return new_view


@auth
@require_http_methods(["GET"])
def user_by_buckid(request, buckid):
    try:
        user = OSUUser.objects.get(shib_id=buckid)
        return JsonResponse(
            {
                "buckid": user.shib_id,
                "nameDotNumber": user.name_num,
                "displayName": user.display_name,
                "discordId": user.discord_id,
            }
        )
    except OSUUser.DoesNotExist:
        return JsonResponse({"error": "not found"}, status=404)


@auth
@require_http_methods(["GET"])
def user_by_discordid(request, discord_id):
    try:
        user = OSUUser.objects.get(discord_id=discord_id)
        return JsonResponse(
            {
                "buckid": user.shib_id,
                "nameDotNumber": user.name_num,
                "displayName": user.display_name,
                "discordId": user.discord_id,
            }
        )
    except OSUUser.DoesNotExist:
        return JsonResponse({"error": "not found"}, status=404)


@auth
@require_http_methods(["GET"])
def user_by_buckid(request, buckid):
    try:
        user = OSUUser.objects.get(shib_id=buckid)
        return JsonResponse(
            {
                "buckid": user.shib_id,
                "nameDotNumber": user.name_num,
                "displayName": user.display_name,
                "discordId": user.discord_id,
            }
        )
    except OSUUser.DoesNotExist:
        return JsonResponse({"error": "not found"}, status=404)


@csrf_exempt
@require_http_methods(["POST"])
@auth
def attend(request):
    try:
        print("got attend request")
        body = json.loads(request.body)
        print(body)
        buckid = int(body["buckid"])
        name_num = str(body["nameDotNumber"])
        display_name = str(body["displayName"])
    except:
        return JsonResponse({"error": "bad request"}, status=400)

    try:
        user = OSUUser.objects.get(shib_id=buckid)
        is_new = False
    except OSUUser.DoesNotExist:
        user = OSUUser(
            shib_id=buckid,
            display_name=display_name,
            name_num=name_num,
            affiliation="STUDENT",
        )
        user.save()
        is_new = True

    in_person = True  # todo: make configurable
    attend_type = (
        AttendanceRecord.ATTENDANCE_TYPE.IN_PERSON
        if in_person
        else AttendanceRecord.ATTENDANCE_TYPE.ONLINE
    )

    if user.can_submit_attendance():
        ar = AttendanceRecord(
            user=user,
            date_recorded=datetime.now().astimezone(settings.TIMEZONE),
            attend_type=attend_type,
        )
        ar.save()
    else:
        return JsonResponse(
            {"error": "user cannot submit attendance", "isNew": is_new}, status=403
        )
    return JsonResponse({"success": "user attendance recorded", "isNew": is_new})


@csrf_exempt
@require_http_methods(["POST"])
@auth
def set_discordid_by_buckid(request, buckid):
    try:
        discord_id = int(request.body)

        try:
            discord_linked = OSUUser.objects.get(discord_id=discord_id)

            if discord_linked.shib_id != buckid:
                print(f"ID is already linked to user {discord_linked.shib_id}")
                return JsonResponse(
                    {
                        "error": "Your discord is already linked to another OSU user.",
                    }
                )
        except OSUUser.DoesNotExist:
            pass

        user = OSUUser.objects.get(shib_id=buckid)

        old_discord = user.discord_id  # possibly None
        user.discord_id = discord_id
        user.save()
        print(f"Successfully linked {user.discord_id} to shib {user.shib_id}")
        return JsonResponse(
            {
                "success": True,
                "affiliation": user.affiliation,
                "oldDiscord": old_discord,
            }
        )

    except OSUUser.DoesNotExist:
        return JsonResponse({"error": "not found"}, status=404)
    except ValueError:
        return JsonResponse({"error": "invalid discordid"}, status=400)
