# Generated by Django 3.2.7 on 2021-11-30 05:18

from django.db import migrations, models


class Migration(migrations.Migration):

    dependencies = [
        ('auth_webapp', '0007_alter_osuuser_discord_id'),
    ]

    operations = [
        migrations.AddField(
            model_name='osuuser',
            name='voted',
            field=models.BooleanField(default=False),
        ),
    ]