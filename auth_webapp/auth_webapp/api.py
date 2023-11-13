from django.http import JsonResponse
from django.conf import settings
from django.views.decorators.http import require_http_methods
from .models import OSUUser, AttendanceRecord
from datetime import datetime
from django.views.decorators.csrf import csrf_exempt


def auth(view):
    def new_view(request, *args, **kwargs):
        header = request.headers.get("Authorization")
        if header:
            if header != f"Bearer {settings.API_TOKEN}":
                print(header)
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


@csrf_exempt
@require_http_methods(["POST"])
@auth
def attend(request, buckid):
    try:
        user = OSUUser.objects.get(shib_id=buckid)
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
            return JsonResponse({"error": "user cannot submit attendance"}, status=403)
        return JsonResponse({"success": "user attendance recorded"})
    except OSUUser.DoesNotExist:
        return JsonResponse({"error": "not found"}, status=404)
