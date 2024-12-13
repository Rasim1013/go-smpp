# go-smpp
Как использовать данный сервис:
имеется 2 роута.
1) 1013:/smpp-api/v1/sendone
{
  "server": "1.1.1.1:1013",
  "username": "user",
  "password": "pass",
  "sender": "MRPRO",
  "msisdn": "992123456789",
  "message": "Hello FROM MRPRO!"
}   
2) 1013^/smpp-api/v1/sendbulk
{
  "server": "1.1.1.1:1013",
  "username": "user",
  "password": "pass",
  "sender": "MRPRO",
  "msisdn": ["992123456789","992987654321"],
  "message": "Hello FROM MRPRO!"
}   

