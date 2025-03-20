# pip install docxtpl
import json
import locale
import os
import sys
from datetime import datetime as dt
from docxtpl import DocxTemplate

def filling_doc(input_file):
    # Проверяем существование файла
    if not os.path.exists(input_file):
        print(f"Файл {input_file} не найден")
        return

    # Устанавливаем локаль для корректного отображения дат
    locale.setlocale(locale.LC_ALL, '')
    
    # Загружаем шаблон документа
    doc = DocxTemplate("eptemplate.docx")
    
    # Загружаем все данные из JSON файла
    try:
        with open(input_file, encoding='utf-8') as f:
            data = json.load(f)
    except json.JSONDecodeError:
        print("Ошибка при чтении JSON файла")
        return
    
    # Создаем словарь для подстановки, используя ключи из JSON
    context = {}
    for item in data:
        for key, value in item.items():
            context["#"+key] = value
    
    # Заполняем шаблон и сохраняем документ
    try:
        doc.render(context)
        doc.save(os.path.join(os.getcwd(), 'result.docx'))
    except Exception as e:
        print(f"Ошибка при создании документа: {str(e)}")

def main():
    # Проверяем наличие аргумента
    if len(sys.argv) < 2:
        print("Необходимо указать имя JSON файла в качестве параметра")
        return
    
    input_file = sys.argv[1]
    filling_doc(input_file)

if __name__ == "__main__":
    main()