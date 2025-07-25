from django.urls import path
from . import views_api
from . import views


urlpatterns = [
    path('api/transactions/', views_api.get_transactions),
    path('', views.index, name='index'),
    path('add/', views.add_budget, name='add_budget'),
    path('<int:id>/', views.budget_detail, name='detail'),
]