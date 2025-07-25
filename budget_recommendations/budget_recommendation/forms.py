from django import forms
from .models import Budget

class TransactionForm(forms.Form):
    amount = forms.DecimalField(min_value=0.01)  # Запрет отрицательных сумм
    date = forms.DateField(widget=forms.DateInput(attrs={'type': 'date'}))  # Удобный выбор даты

class BudgetForm(forms.ModelForm):
    class Meta:
        model = Budget
        fields = ['name', 'amount']