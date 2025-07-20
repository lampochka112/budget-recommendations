# finance/services.py
import requests
from django.core.cache import cache
from django.conf import settings

def get_currency_rates():
    cached_rates = cache.get('currency_rates')
    if cached_rates:
        return cached_rates
    
    try:
        response = requests.get('https://www.cbr-xml-daily.ru/daily_json.js')
        data = response.json()
        rates = {
            'USD': data['Valute']['USD']['Value'],
            'EUR': data['Valute']['EUR']['Value'],
        }
        cache.set('currency_rates', rates, timeout=3600)  # Кешируем на 1 час
        return rates
    except Exception as e:
        # Логирование ошибки
        return None
    
    # finance/services.py
from .models import Transaction, Recommendation

def get_recommendations(user):
    # Анализ расходов пользователя
    categories = Transaction.objects.filter(
        user=user, 
        transaction_type=Transaction.EXPENSE
    ).values('category').annotate(total=models.Sum('amount')).order_by('-total')
    
    recommendations = []
    for cat in categories[:3]:  # Топ-3 категории по расходам
        recs = Recommendation.objects.filter(category_id=cat['category'], is_general=False)
        if not recs.exists():
            recs = Recommendation.objects.filter(is_general=True)
        recommendations.extend(recs)
    
    return recommendations