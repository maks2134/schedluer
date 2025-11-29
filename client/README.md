# Schedluer Frontend

Современный фронтенд для просмотра расписания БГУИРа, построенный на Next.js 15, TypeScript, Tailwind CSS и Shadcn/ui.

## Технологии

- **Next.js 15** - React фреймворк с App Router
- **TypeScript** - Типизированный JavaScript
- **Tailwind CSS** - Utility-first CSS фреймворк
- **Shadcn/ui** - Высококачественные компоненты UI
- **Lucide React** - Иконки

## Установка

```bash
# Установить зависимости
npm install

# Создать файл .env.local
cp .env.local.example .env.local

# Отредактировать .env.local и указать URL бэкенда
# NEXT_PUBLIC_API_URL=http://localhost:8080/api/v1
```

## Запуск

### Режим разработки

```bash
npm run dev
```

Приложение будет доступно по адресу [http://localhost:3000](http://localhost:3000)

### Production сборка

```bash
npm run build
npm start
```

## Структура проекта

```
client/
├── app/                    # Next.js App Router
│   ├── layout.tsx         # Корневой layout
│   ├── page.tsx           # Главная страница
│   └── globals.css        # Глобальные стили
├── components/            # React компоненты
│   ├── ui/                # Shadcn/ui компоненты
│   ├── schedule/          # Компоненты расписания
│   ├── groups/            # Компоненты групп
│   └── employees/         # Компоненты преподавателей
├── lib/                   # Утилиты и хелперы
│   ├── api.ts             # API клиент
│   └── hooks/             # React хуки
├── types/                 # TypeScript типы
│   └── api.ts             # Типы для API
└── public/                # Статические файлы
```
