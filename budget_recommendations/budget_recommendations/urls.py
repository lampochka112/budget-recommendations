from django.urls import path
from . import views_api

urlpatterns = [
    path('api/transactions/', views_api.get_transactions),
]