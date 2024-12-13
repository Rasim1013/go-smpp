# go-smpp

## Как использовать данный сервис

Сервис предоставляет маршруты для отправки сообщений через SMPP.

---

### Примеры использования

**1. Отправка одиночного сообщения**  
**2. Отправка массовых сообщений**

**Маршрут:** `POST /smpp-api/v1/sendone` и `POST /smpp-api/v1/sendbulk`

**Пример запроса:**

```json
{
  "server": "1.1.1.1:1013",
  "username": "user",
  "password": "pass",
  "sender": "MRPRO",
  "msisdn": "992123456789",                // Для одиночного сообщения
  "message": "Hello FROM MRPRO!"
}


curl -X POST http://your-domain.com/smpp-api/v1/sendone \
-H "Content-Type: application/json" \
-d '{
  "server": "1.1.1.1:1013",
  "username": "user",
  "password": "pass",
  "sender": "MRPRO",
  "msisdn": "992123456789",
  "message": "Hello FROM MRPRO!"
}'

curl -X POST http://your-domain.com/smpp-api/v1/sendbulk \
-H "Content-Type: application/json" \
-d '{
  "server": "1.1.1.1:1013",
  "username": "user",
  "password": "pass",
  "sender": "MRPRO",
  "msisdn": ["992123456789", "992987654321"],
  "message": "Hello FROM MRPRO!"
}'
