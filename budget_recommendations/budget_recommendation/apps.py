from flask import Flask, render_template, request, jsonify
import pickle
import pandas as pd
import numpy as np

app = Flask(__name__)

# Загрузка модели и данных
try:
    with open('model.pkl', 'rb') as f:
        model = pickle.load(f)
except FileNotFoundError:
    model = None
    print("Предупреждение: Модель model.pkl не найдена")

try:
    data = pd.read_csv('data.csv')
except FileNotFoundError:
    data = pd.DataFrame()
    print("Предупреждение: Файл data.csv не найден")

@app.route('/')
def index():
    return render_template('index.html')

@app.route('/recommend', methods=['POST'])
def recommend():
    try:
        if model is None:
            return jsonify({'error': 'Модель не загружена'})
        
        # Получаем данные из формы
        age = float(request.form['age'])
        income = float(request.form['income'])
        expenses = float(request.form['expenses'])
        savings = float(request.form['savings'])
        
        # Создаем DataFrame для предсказания
        input_data = pd.DataFrame({
            'age': [age],
            'income': [income],
            'expenses': [expenses],
            'savings': [savings]
        })
        
        # Делаем предсказание
        prediction = model.predict(input_data)[0]
        
        # Формируем рекомендации
        recommendations = generate_recommendations(age, income, expenses, savings, prediction)
        
        return jsonify({
            'prediction': float(prediction),
            'recommendations': recommendations
        })
        
    except Exception as e:
        return jsonify({'error': str(e)})

def generate_recommendations(age, income, expenses, savings, prediction):
    recommendations = []
    
    # Анализ соотношения доходов и расходов
    if expenses > income * 0.7:
        recommendations.append("Сократите расходы: ваши расходы превышают 70% доходов")
    elif expenses > income * 0.5:
        recommendations.append("Оптимизируйте расходы: рассмотрите возможность сокращения необязательных трат")
    
    # Анализ сбережений
    if savings < income * 3:
        recommendations.append("Увеличьте сбережения: рекомендуется иметь минимум 3 месячных дохода в качестве резерва")
    
    # Возрастные рекомендации
    if age < 30:
        recommendations.append("Рассмотрите долгосрочные инвестиции: у вас есть время для роста капитала")
    elif age > 50:
        recommendations.append("Увеличьте консервативные инвестиции: снизьте риски в инвестиционном портфеле")
    
    # Общие рекомендации
    recommendations.extend([
        "Ведите бюджет: отслеживайте все доходы и расходы",
        "Создайте финансовую подушку безопасности",
        "Диверсифицируйте источники дохода",
        "Регулярно пересматривайте финансовый план"
    ])
    
    return recommendations

if __name__ == '__main__':
    app.run(debug=True, host='0.0.0.0', port=5000)