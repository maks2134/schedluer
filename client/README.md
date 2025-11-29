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

## Функциональность

### Просмотр расписания групп
- Поиск групп по номеру, специальности, факультету
- Отображение расписания по дням недели
- Показ экзаменов
- Информация о парах: время, аудитория, преподаватель, тип занятия

### Просмотр расписания преподавателей
- Поиск преподавателей по ФИО, должности
- Отображение расписания преподавателя
- Информация о группах и занятиях

### Дополнительные возможности
- Обновление данных из API БГУИРа
- Кэширование данных
- Адаптивный дизайн
- Темная тема (поддержка через Shadcn/ui)

## API Endpoints

Фронтенд использует следующие endpoints:

- `GET /api/v1/groups` - Список всех групп
- `GET /api/v1/groups/:groupNumber` - Информация о группе
- `GET /api/v1/schedule/group/:groupNumber` - Расписание группы
- `GET /api/v1/employees` - Список всех преподавателей
- `GET /api/v1/employees/:urlId` - Информация о преподавателе
- `GET /api/v1/schedule/employee/:urlId` - Расписание преподавателя
- `POST /api/v1/groups/refresh` - Обновить список групп
- `POST /api/v1/employees/refresh` - Обновить список преподавателей
- `POST /api/v1/schedule/group/:groupNumber/refresh` - Обновить расписание группы
- `POST /api/v1/schedule/employee/:urlId/refresh` - Обновить расписание преподавателя

## Настройка

### Переменные окружения

Создайте файл `.env.local`:

```env
NEXT_PUBLIC_API_URL=http://localhost:8080/api/v1
```

### Кастомизация

- **Цвета**: Настройте в `app/globals.css` через CSS переменные
- **Компоненты**: Используйте Shadcn/ui CLI для добавления новых компонентов
- **Стили**: Tailwind конфигурация в `tailwind.config.ts`

## Разработка

### Добавление новых компонентов Shadcn/ui

```bash
npx shadcn@latest add [component-name]
```

### Проверка типов

```bash
npm run type-check
```

### Линтинг

```bash
npm run lint
```

## Деплой

### Vercel (рекомендуется)

1. Подключите репозиторий к Vercel
2. Укажите переменную окружения `NEXT_PUBLIC_API_URL`
3. Деплой произойдет автоматически

### Docker

```bash
docker build -t schedluer-frontend .
docker run -p 3000:3000 -e NEXT_PUBLIC_API_URL=http://your-api-url schedluer-frontend
```

## Лицензия

Apache 2.0
