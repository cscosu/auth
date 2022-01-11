from django.db import models
from django.contrib.auth.models import AbstractBaseUser, PermissionsMixin
import mailchimp_marketing as MailchimpMarketing
from mailchimp_marketing.api_client import ApiClientError
import os
import hashlib

try:
    client = MailchimpMarketing.Client()
    client.set_config({"api_key": os.environ["MAILCHIMP_API_KEY"], "server": "us16"})
except ApiClientError as error:
    print("MAILCHIMP ERROR: " + error)

# class OSUUserManager(BaseUserManager):
#     # This makes stuff work like manage.py createsuperuser
#     def create_user(shib_id, display_name, affiliation, name_num):
#         print("Creating user with shib %s name %s aff %s name.# %s" % (shib_id, display_name, affiliation, name_num))
#         user, created = OSUUser.objects.get_or_create(shib_id=shib_id, defaults={"display_name": display_name, "name_num": name_num, "affiliation": affiliation})
#         return user

#     def create_superuser(shib_id, display_name=None, affiliation=None, name_num=None, password=None):
#         print("Creating superuser with shib %s name %s aff %s name.# %s" % (shib_id, display_name, affiliation, name_num))
#         user, created = OSUUser.objects.get_or_create(shib_id=shib_id, defaults={"display_name": display_name, "name_num": name_num, "affiliation": affiliation})
#         user.is_superuser = True
#         return user

# If you change these, then if you fully restart the docker containers it will automatically run
# the migrations for you (until you CTRL-C to stop your containers and then re-start them,
# # you'll get db errors)


class OSUUser(AbstractBaseUser, PermissionsMixin):
    USERNAME_FIELD = "shib_id"
    REQUIRED_FIELDS = ["display_name", "name_num"]
    shib_id = models.CharField(max_length=50, unique=True, primary_key=True)
    name_num = models.CharField(max_length=50)
    display_name = models.CharField(max_length=100)
    # You can change this type if needed
    discord_id = models.CharField(max_length=100, null=True, blank=True)
    last_login = models.DateTimeField(null=True)
    voted = models.BooleanField(default=False)

    # Affiliation is going to be interesting. In theory one user can have more than one affiliation.
    # But here I'm assuming we pick one.
    # More info: https://webauth.service.ohio-state.edu/~shibboleth/user-attribute-reference.html#edupersonscopedaffiliation
    aff_choices = models.TextChoices("Affiliation", "STUDENT FACULTY_STAFF ALUMNI NONE")
    affiliation = models.CharField(
        max_length=20, choices=aff_choices.choices, default=aff_choices.NONE
    )
    is_admin = models.BooleanField(default=False)

    # This is effectively a cache for the function below
    added_to_mailing_list = models.BooleanField(default=False)

    def is_currently_on_mailing_list(self):
        if self.added_to_mailing_list:
            return True
        email = (self.name_num + "@osu.edu").lower()
        try:
            list_info = client.lists.get_list_member(
                os.environ["MAILCHIMP_LIST_ID"], hashlib.md5(email.encode()).hexdigest()
            )
            if list_info["status"] == "subscribed":
                self.added_to_mailing_list = True
                self.save()
                return True
            else:
                print("Subscriber exists but not subscribed: " + email)
                return False
        except ApiClientError:
            return False

    def add_to_mailing_list(self):
        email = (self.name_num + "@osu.edu").lower()
        try:
            list_info = client.lists.add_list_member(
                os.environ["MAILCHIMP_LIST_ID"],
                {"email_address": email, "status": "subscribed"},
            )
            if list_info["status"] == "subscribed":
                self.added_to_mailing_list = True
                self.save()
                return True
            else:
                print("New sub status not subscribed: " + str(list_info))
        except ApiClientError as e:
            print("api error")
            print(e.text)
            pass
        return False

    # Some helpers the authentication module uses
    def get_full_name(self):
        return self.display_name

    def get_short_name(self):
        return self.name_num

    def is_staff(self):  # gives access to admin site
        return self.is_admin  # comes from PermissionMixin

    def is_active(self):
        return True


class AttendanceRecord(models.Model):
    id = models.BigAutoField(primary_key=True)
    user = models.ForeignKey(OSUUser, on_delete=models.CASCADE)
    date_recorded = models.DateTimeField()

    ATTENDANCE_TYPE = models.TextChoices("AttendType", "IN_PERSON ONLINE DEFAULT")
    attend_type = models.CharField(
        max_length=20, choices=ATTENDANCE_TYPE.choices, default=ATTENDANCE_TYPE.DEFAULT
    )


# TODO: Add Meeting model and possible per-meeting survey questions or something
