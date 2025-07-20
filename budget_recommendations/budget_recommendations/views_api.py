from django.http import JsonResponse

def get_transactions(request):
    if not request.user.is_authenticated:
        return JsonResponse([], safe=False)
    
    transactions = Transaction.objects.filter(
        user=request.user
    ).values('amount', 'date', 'category__name')[:50]  # Только нужные поля
    
    return JsonResponse(list(transactions), safe=False)