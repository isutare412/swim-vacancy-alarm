# Swim Vacancy Alarm

판교스포츠센터 빈자리 텔레그램으로 알림받자!

## Configuration

```yaml
search:
  swim-course:
    every: 10s
    course-names:
      - 07시_연수
      - 09시_연수
register:
  seongnam-sdc-url: https://spo.isdc.co.kr/courseRegist.do
telegram:
  bot-token: <bot_token>
  chat-id: <chat_id>
```

## Screenshots

![App logs](docs/asset/app_log.png)

![Telegram alarm](docs/asset/telegram_alarm.png)
