# API для создания поздравительных открыток

1. Получение шаблона
    ```http request
    GET /api/preset?id={INT}
    ```
   Возвращает JSON в формате
   ```json
   {
    "Name": string, // название шаблона
    "PaperSize": string, // размер бумаги
    "Text": string, // текст
    "Greeting": string, // обращение
    "TextX": float64, // X левого угла блока с текстом
    "TextY": float64, // Y левого угла блока с текстом
    "GreetingY": float64, // X левого края строки с поздравлением
    "Image": string, // адрес картинки с фоном шаблона
   }
    ```