"""auth_webapp URL Configuration

The `urlpatterns` list routes URLs to views. For more information please see:
    https://docs.djangoproject.com/en/3.1/topics/http/urls/
Examples:
Function views
    1. Add an import:  from my_app import views
    2. Add a URL to urlpatterns:  path('', views.home, name='home')
Class-based views
    1. Add an import:  from other_app.views import Home
    2. Add a URL to urlpatterns:  path('', Home.as_view(), name='home')
Including another URLconf
    1. Import the include() function: from django.urls import include, path
    2. Add a URL to urlpatterns:  path('blog/', include('blog.urls'))
"""
from django.contrib import admin
from django.urls import path
import auth_webapp.views as views
import auth_webapp.discord_bot as discord_bot_views
import auth_webapp.api as api_views

urlpatterns = [
    path("", views.home, name="home"),
    path("login/", views.login_view, name="login"),
    path("debug_login/", views.debug_login, name="debug_login"),
    path("logout", views.logout_view, name="logout"),
    path("admin/", admin.site.urls),
    path("new", views.new_member, name="new_member"),
    path("attendance", views.attendance, name="attendance"),
    path("attendance/attend", views.attend, name="attend"),
    path("attendance/subscribe", views.subscribe, name="subscribe"),
    path("internal/link_discord", discord_bot_views.link_discord, name="link_discord"),
    path("elections", views.elections, name="elections"),
    path(
        "api/user/bybuckid/<int:buckid>",
        api_views.user_by_buckid,
        name="user_by_buckid",
    ),
    path("api/user/bybuckid/<int:buckid>/attend", api_views.attend, name="api_attend"),
    path(
        "api/user/bybuckid/<int:buckid>/discordid",
        api_views.attend,
        name="api_set_discordid_by_buckid",
    ),
]
