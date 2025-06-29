sheep_farm_backend_go/
├── cmd/                          # نقطه ورود اصلی برنامه
│   └── api/                      # برای وب API
│       └── main.go               # فایل اصلی اجرای سرور
├── internal/                     # کدهای داخلی برنامه (منطق اصلی)
│   ├── domain/                   # لایه دامنه: موجودیت‌ها، مدل‌های داده‌ای
│   │   ├── sheep.go
│   │   ├── vaccine.go
│   │   ├── treatment.go
│   │   ├── reminder.go
│   │   └── errors.go             # خطاهای سفارشی برنامه
│   ├── application/              # لایه کاربرد: Use Cases (خدمات برنامه)
│   │   ├── ports/                # پورت‌ها (اینترفیس‌ها) برای ارتباط با لایه‌های بیرونی
│   │   │   ├── repository.go     # اینترفیس برای پایگاه داده
│   │   │   ├── reminder_notifier.go # اینترفیس برای ارسال یادآورها
│   │   ├── services/             # پیاده‌سازی Use Cases (منطق درخواست/پاسخ برنامه)
│   │   │   ├── sheep_service.go
│   │   │   └── vaccine_service.go
│   │   │   └── reminder_service.go
│   ├── infrastructure/           # لایه زیرساخت: پیاده‌سازی پورت‌ها (آداپتورها)
│   │   ├── persistence/          # پیاده‌سازی‌های مربوط به پایگاه داده
│   │   │   └── firestore_repository.go
│   │   ├── http/                 # پیاده‌سازی وب سرور و هندلرها
│   │   │   ├── router.go         # تعریف مسیرها
│   │   │   ├── handlers/         # توابع هندلر برای APIها
│   │   │   │   ├── sheep_handler.go
│   │   │   │   └── vaccine_handler.go
│   │   │   │   └── auth_handler.go # (اگر احراز هویت اضافه شود)
│   │   │   └── dto/              # Data Transfer Objects (DTOs) برای ورودی/خروجی API
│   │   ├── external/             # پیاده‌سازی سرویس‌های خارجی (مثلاً ارسال پیامک)
│   │   │   └── console_notifier.go # مثال: یادآور به کنسول
│   │   ├── scheduler/            # سیستم زمان‌بندی برای یادآورها
│   │   │   └── cron_scheduler.go # پیاده‌سازی زمان‌بندی با کتابخانه cron
├── go.mod                        # فایل ماژول Go (وابستگی‌ها)
├── go.sum                        # فایل چک‌سام Go
└── .env                          # متغیرهای محیطی (مانند کلیدهای API)