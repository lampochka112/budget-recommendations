# finance/views.py
from django.views.generic import ListView, CreateView, UpdateView, DeleteView
from django.contrib.auth.mixins import LoginRequiredMixin
from django.urls import reverse_lazy
from .models import Transaction, BudgetCategory, FinancialGoal, Recommendation
from .forms import TransactionForm, BudgetCategoryForm, FinancialGoalForm
from .services import get_currency_rates

class TransactionListView(LoginRequiredMixin, ListView):
    model = Transaction
    template_name = 'finance/transaction_list.html'
    paginate_by = 10
    filterset_class = None  # Можно использовать django-filters
    
    def get_queryset(self):
        queryset = super().get_queryset().filter(user=self.request.user)
        # Фильтрация и сортировка
        return queryset.order_by('-date')
    
    def get_context_data(self, **kwargs):
        context = super().get_context_data(**kwargs)
        context['currency_rates'] = get_currency_rates()
        return context

class TransactionCreateView(LoginRequiredMixin, CreateView):
    model = Transaction
    form_class = TransactionForm
    template_name = 'finance/transaction_form.html'
    success_url = reverse_lazy('transaction-list')
    
    def form_valid(self, form):
        form.instance.user = self.request.user
        return super().form_valid(form)

# Аналогичные представления для других моделей