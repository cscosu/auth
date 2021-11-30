from django.contrib.auth.backends import ModelBackend
from .models import OSUUser


class ShibAuthBackend(ModelBackend):
    """This does not implement shibboleth authentication - we are using nginx
    for that. Trusted headers are sent from nginx to /login and then we explicitly
    call login. This is pretty much a dummy authentication backend.
    """

    def authenticate(self, request):
        # Deny the silly case where they try to login on /admin with a username/password
        return None

    def get_user(self, user_id):
        try:
            return OSUUser.objects.get(shib_id=user_id)
        except OSUUser.DoesNotExist:
            return None
