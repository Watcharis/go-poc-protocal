Service A:
-   ใช้ otelhttp.NewTransport สำหรับ HTTP Client เพื่อส่ง Trace Context ไปยัง Service B
    สร้าง Span ใหม่ (call-service-b) สำหรับการเรียก Service B


Service B:
-   ใช้ otelhttp.NewHandler สำหรับ HTTP Server เพื่อดึง Trace Context จาก Header
    สร้าง Span ใหม่ (process-request) ใน Trace เดียวกันกับ Service A