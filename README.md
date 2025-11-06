# پروژه نمونه آموزشگاه موسیقی با ClickHouse

این مخزن شامل یک نمونه‌ی ساده از سامانه آموزشگاه موسیقی است که با ClickHouse به عنوان دیتابیس و یک سرویس API مبتنی بر
Go (Gin) پیاده‌سازی شده و به کمک Docker اجرا می‌شود.

## اجرا به صورت سریع

1) ساخت و اجرا:

```bash
docker compose up -d --build
```

2) بررسی سلامت سرویس‌ها:

- ClickHouse HTTP: `http://localhost:8123/`
- API: `http://localhost:8000/health`

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

## ساختار پروژه

- `docker-compose.yml`: تعریف سرویس‌ها
- `clickhouse/init/`: اسکریپت‌های ساخت اسکیما و داده‌ی نمونه
- `app/`: سرویس Go/Gin
    - `Dockerfile`
    - `go.mod`, `go.sum`
    - `main.go`: مسیرهای API و اتصال ClickHouse

## ویکی و مستندات

- مستندات کامل به زبان فارسی در پوشه‌ی `ClickHouseDB.wiki` موجود است.
- از صفحه‌ی [Home](ClickHouseDB.wiki/Home.md) شروع کنید.

## توقف و پاک‌سازی

```bash
docker compose down
# حذف کامل داده‌های لوکال (اختیاری)
docker volume rm clickhousedb_clickhouse_data
```

## مجوز

MIT


