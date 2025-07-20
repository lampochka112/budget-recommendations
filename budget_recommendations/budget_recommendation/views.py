from django.views.generic import ListView
from .models import Transaction

class SimpleTransactionView(ListView):
    model = Transaction
    template_name = 'finance/simple_list.html'
    context_object_name = 'transactions'
    
    def get_queryset(self):
        return Transaction.objects.filter(user=self.request.user).order_by('-date')[:20]  # Только 20 последних