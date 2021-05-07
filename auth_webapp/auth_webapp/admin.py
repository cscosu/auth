from django.contrib.auth.admin import UserAdmin as BaseUserAdmin
from django.contrib import admin
from .models import OSUUser, AttendanceRecord
from django import forms


class OSUUserChangeForm(forms.ModelForm):
    """Admin form for updating users

    intentionally few attributes here
    """
    class Meta:
        model = OSUUser
        fields = ('affiliation', 'is_admin')

class OSUUserAdmin(BaseUserAdmin):
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

admin.site.register(OSUUser, OSUUserAdmin)
admin.site.register(AttendanceRecord)