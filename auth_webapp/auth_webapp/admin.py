import csv
import base64
import requests as req
import json
from django.contrib.auth.admin import UserAdmin as BaseUserAdmin
from django.contrib import admin
from .models import OSUUser, AttendanceRecord
from django import forms
from django.http import HttpResponse


class ExportCsvMixin:
    def export_as_csv(self, request, queryset):

        meta = self.model._meta
        field_names = [field.name for field in meta.fields]

        response = HttpResponse(content_type='text/csv')
        response['Content-Disposition'] = 'attachment; filename={}.csv'.format(meta)
        writer = csv.writer(response)

        writer.writerow(field_names)
        for obj in queryset:
            row = writer.writerow([getattr(obj, field) for field in field_names])

        return response

    export_as_csv.short_description = "Export Selected"


def is_student(user):
    """Determins if name.num is current student

    Querys API to see if name.num is current student. Returns
    True if yes, returns False otherwise

    Note: errors if students not found in by API"""
    API_URL = base64.b64decode("aHR0cDovL2ttZGF0YS5vc3UuZWR1").decode()
    PEOPLE_URL = '/people'

    req_url = API_URL + PEOPLE_URL + "/" + user + ".json"
    response = req.get(req_url)
    if response.status_code == 200:
        affiliation = json.loads(response.content)["affiliation"]
        return "Student" in affiliation.split(", ")
    else:
        return False

def check_user_alum_queryset(queryset):
    checked = {}
    errors = {}
    for obj in queryset:
        if obj.affiliation == OSUUser.aff_choices.STUDENT:
            try:
                still_student = is_student(obj.name_num)
                checked[obj.name_num] = still_student
                if not still_student:
                    print(f"User {obj.shib_id} no longer student, marking as alumni.")
                    obj.affiliation = OSUUser.aff_choices.ALUMNI
                    obj.save()
            except Exception as e:
                errors[obj.name_num] = str(e)
                print(f"User {obj.shib_id} errored when checking student status: {str(e)}. Leaving status as f{str(obj.affiliation)}.")

    return checked, errors

@admin.action(description='Update student/alumi status')
def check_user_alum(self, request, queryset):
    """Admin action to update student/alumni status

    will remove student status if deemed not student.
    Querys an API for studen status
    Also returns json with list of checked students
    (called API on name.num) with T/F for student/alumni
    as well as a list of errors encountered (with name.num)
    """
    checked, errors = check_user_alum_queryset(queryset)
    rtnDict = {}
    rtnDict["checked"] = checked
    rtnDict["errors"] = errors
    rtnJson = json.dumps(rtnDict)
    response = HttpResponse(content_type="application/json")
    response.write(rtnJson)
    return response


class OSUUserChangeForm(forms.ModelForm):
    """Admin form for updating users

    intentionally few attributes here
    """
    class Meta:
        model = OSUUser
        fields = ('affiliation', 'is_admin')

class OSUUserAdmin(BaseUserAdmin, ExportCsvMixin):
    # The forms to add and change user instances
    form = OSUUserChangeForm
    add_form = None

    # The fields to be used in displaying the User model.
    # These override the definitions on the base UserAdmin
    # that reference specific fields on auth.User.
    list_display = ('shib_id', 'name_num', 'display_name', 'discord_id', 'last_login', 'affiliation', 'added_to_mailing_list')
    list_filter = ('is_superuser',)
    fieldsets = (
        (None, {'fields': ('shib_id', 'name_num')}),
        ('Personal info', {'fields': ('display_name','affiliation')}),
        ('Associations', {'fields': ('discord_id', 'added_to_mailing_list')}),
        ('Permissions', {'fields': ('is_superuser',)}),
    )
    search_fields = ('name_num',)
    ordering = ('last_login',)
    filter_horizontal = ()
    readonly_fields = ('shib_id', 'name_num')

    actions = ["export_as_csv", check_user_alum]


admin.site.register(OSUUser, OSUUserAdmin)
admin.site.register(AttendanceRecord)
