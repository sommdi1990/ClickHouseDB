# Music school API with Go and ClickHouse

A Gin service records students and enrollments and exposes the course catalog over HTTP. The schema separates the school entities from enrollment data, and the seed script supplies instruments for the sample dataset.

Follow POST /students and POST /enrollments in the Go entry point, then compare the fields written by each handler with the SQL definitions. The students and courses read paths order their results by creation time. The schema and sample rows are maintained in separate files.

- [main.go](app/main.go)
- [01_schema.sql](clickhouse/init/01_schema.sql)
- [02_seed.sql](clickhouse/init/02_seed.sql)

## اجرا به صورت سریع

1) ساخت و اجرا:

```bash
docker compose up -d --build
```

2) بررسی سلامت سرویس‌ها:

- ClickHouse HTTP: `http://localhost:8123/`
- API: `http://localhost:8000/health`
- Tabix (رابط وب ClickHouse): `http://localhost:8080`

3) تست‌های سریع:

```bash
# لیست هنرجوها
curl http://localhost:8000/students

# ایجاد هنرجو
curl -X POST http://localhost:8000/students \
  -H 'Content-Type: application/json' \
  -d '{
        "first_name": "سینا",
        "last_name": "کاظمی",
        "phone": "09120009999",
        "email": "sina@example.com",
        "birth_date": "2004-07-21"
      }'
```

## رابط گرافیکی وب (Tabix)

- آدرس دسترسی: `http://localhost:8080`
- ساخت اتصال جدید در Tabix:
    - Host: `clickhouse` (از داخل شبکه Docker) یا `localhost`
    - Port: `8123`
    - User: `default`
    - Password: خالی در لوکال (یا طبق تنظیمات شما)
    - Database: `music_school`

نکته: اگر برای کاربر `default` رمز گذاشته‌اید، همان رمز را در Tabix وارد کنید. اگر از سرویس دیگری در همان شبکه Docker
وصل می‌شوید، نام سرویس `clickhouse` را به‌عنوان میزبان قرار دهید.

## مجوز

MIT
