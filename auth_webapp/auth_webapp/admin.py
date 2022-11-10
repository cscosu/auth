import csv
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

    actions = ["export_as_csv"]


admin.site.register(OSUUser, OSUUserAdmin)
admin.site.register(AttendanceRecord)
