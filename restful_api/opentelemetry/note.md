# OpenTelemetry (OTel) การใช้ Log + Distributed Tracing
```เป็นเฟรมเวิร์กที่ใช้สำหรับการสังเกตการณ์ระบบซอฟต์แวร์ (observability) และช่วยรวบรวมข้อมูลที่สำคัญในระบบกระจาย (distributed systems) ซึ่งข้อมูลดังกล่าวประกอบด้วย Logs, Distributed Tracing และ Metrics โดยแต่ละอย่างมีบทบาทเฉพาะตัวและสัมพันธ์กันในมุมของการช่วยระบุปัญหาและปรับปรุงประสิทธิภาพของระบบ```

## Logs
    - คืออะไร: Logs เป็นข้อความที่ระบบหรือแอปพลิเคชันสร้างขึ้น เพื่อบันทึกเหตุการณ์ที่เกิดขึ้น เช่น ข้อผิดพลาด (errors), คำเตือน (warnings), หรือข้อมูลดีบ(debugging information)
    - บทบาท: Logs ช่วยบอกว่าเกิดอะไรขึ้นในแอปพลิเคชันในเวลาที่แน่นอน โดยมักจะละเอียดและเฉพาะเจาะจง
    - ตัวอย่างการใช้งาน: เมื่อเซิร์ฟเวอร์ล้มเหลว Logs จะบันทึกรายละเอียด เช่น "Database connection timeout at 10:05 PM"

#  Distributed Tracing
    - คืออะไร: Distributed Tracing ใช้ติดตามคำขอ (requests) หรือกระบวนการ (transactions) ที่วิ่งผ่านหลายๆ บริการ (services) ในระบบกระจาย เช่น Microservices
    - บทบาท: ช่วยให้เข้าใจว่าเส้นทางของคำขอนั้นผ่านบริการไหนบ้าง และใช้เวลาเท่าไรในแต่ละจุด
    - ตัวอย่างการใช้งาน: หากระบบช้าลง Distributed Tracing ช่วยระบุว่าความล่าช้าเกิดที่บริการใด เช่น API Gateway หรือ Database

# Metrics
    - คืออะไร: Metrics เป็นข้อมูลเชิงตัวเลขที่เก็บรวบรวมเพื่อบอกสถานะและประสิทธิภาพของระบบ เช่น CPU Usage, Memory Usage, Request Rate, Error Rate
    - บทบาท: ใช้ตรวจวัดและตั้งค่าการแจ้งเตือน (alerts) หากค่าของ Metrics ผิดปกติ
    - ตัวอย่างการใช้งาน: ถ้าจำนวน Error Rate เพิ่มขึ้น Metrics จะช่วยเตือนว่ามีปัญหา


### ความสัมพันธ์ระหว่าง Logs, Tracing และ Metrics
ทั้งสามทำงานร่วมกันเพื่อให้มุมมองที่ครอบคลุมเกี่ยวกับระบบ:
- Logs + Metrics: หาก Metrics บ่งชี้ถึงปัญหา เช่น Error Rate เพิ่มขึ้น Logs จะช่วยให้เจาะลึกได้ว่า Error เกิดจากอะไร
- Tracing + Metrics: Distributed Tracing สามารถบอกได้ว่าปัญหาเกิดขึ้นที่บริการใด ขณะที่ Metrics ช่วยติดตามแนวโน้มของปัญหานั้น
- Tracing + Logs: Tracing ช่วยแสดงเส้นทางการทำงานและบริบท (context) ของคำขอ ขณะที่ Logs บันทึกรายละเอียดเฉพาะของแต่ละเหตุการณ์ในคำขอนั้น

Service A:
-   ใช้ otelhttp.NewTransport สำหรับ HTTP Client เพื่อส่ง Trace Context ไปยัง Service B
    สร้าง Span ใหม่ (call-service-b) สำหรับการเรียก Service B


Service B:
-   ใช้ otelhttp.NewHandler สำหรับ HTTP Server เพื่อดึง Trace Context จาก Header
    สร้าง Span ใหม่ (process-request) ใน Trace เดียวกันกับ Service A