# finance/tasks.py
from celery import shared_task
from django.core.mail import send_mail
from django.conf import settings

@shared_task
def send_weekly_report(user_id):
    from users.models import User
    from finance.models import Transaction
    try:
        user = User.objects.get(pk=user_id)
        transactions = Transaction.objects.filter(user=user).order_by('-date')[:10]
        
        message = f"Ваши последние транзакции:\n\n"
        for t in transactions:
            message += f"{t.date}: {t.amount} ({t.get_transaction_type_display()})\n"
        
        send_mail(
            'Ваш еженедельный финансовый отчет',
            message,
            settings.DEFAULT_FROM_EMAIL,
            [user.email],
            fail_silently=False,
        )
        return f"Report sent to {user.email}"
    except Exception as e:
        return str(e)