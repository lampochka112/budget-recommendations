from django import forms

class TransactionForm(forms.Form):
    amount = forms.DecimalField(min_value=0.01)  # Запрет отрицательных сумм
    date = forms.DateField(widget=forms.DateInput(attrs={'type': 'date'}))  # Удобный выбор даты