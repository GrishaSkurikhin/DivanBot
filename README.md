# DivanBot
### [Телеграм-бот](https://t.me/divan_movienights_bot) для киноклуба МАИ "Диван"

Проект полностью написан на языке Golang, с использованием пакета [go-telegram](github.com/go-telegram/bot).

**Возможности, предоставляемые пользователям:**
- Просмотр фильмов, которые планируется показывать
- Просмотр уже показанных фильмов
- Запись и отмена записи на фильм в один клик (перед этим необходимо зарегестрироваться в системе)
- Просмотр информации о предстоящем показе фильма: место проведения, количество мест, открыта ли регистрация
- Просмотр своих записей на фильмы
- Отправка отзыва или предложения

Для удобства использования бота, в проекте используются ui компоненты, предоставляемые API telegram:
- **ReplyKeyboard** - клавиатура под строкой ввода сообщения

![image](https://github.com/GrishaSkurikhin/DivanBot/assets/71190776/6c724aea-9b51-4218-91cd-85d198dbe4a7)
- **InlineKeyboard** - клавиатура под сообщением

![image](https://github.com/GrishaSkurikhin/DivanBot/assets/71190776/7886372d-aad1-49dc-abdc-d469ca3be48a)
- **Slider** - компонент, основанный на InlineKeyboard, для просмотра нескольких слайдов в одном сообщении

![image](https://github.com/GrishaSkurikhin/DivanBot/assets/71190776/92253b41-ebc8-4cc7-b295-c229091ed852)

## Схема деплоя
Проект залит на Yandex Cloud с применением serverless-технологий, что существенно снизило затраты на обслуживание бота.

**Serverless** - это модель предоставления серверных услуг без аренды или покупки оборудования. 
При таком подходе управлением ресурсами инфраструктуры, её настройкой и обслуживанием занимается провайдер.
Серверное оборудование задействуется только тогда, когда срабатывает заранее запрограммированное событие, или триггер.
В случае бота, на него приходит сообщение от пользователя.

![Drawing 2023-11-23 13 51 53 excalidraw](https://github.com/GrishaSkurikhin/DivanBot/assets/71190776/69cc0cf0-4efe-45da-ae93-ab8e5a7e8005)

## База данных
Для хранения информации была выбрана база данных Yandex Database, потому что она также предоставляет возможность serverless-вычислений.

Схема БД:

![бот - схема бд](https://github.com/GrishaSkurikhin/DivanBot/assets/71190776/65cec956-c4ff-4944-8008-322b11ea334b)

## Доработки пакета go-telegram

- Была разработана система диалогов между пользователем и ботом. Диалоги используются при регистрации, изменении своих данных и отпправке отзывов. В основе диалогов лежит компонент, отвечающий за хранение состояния диалога (для хранения состояний также используется yandex DB).
- Были переработаны ui компоненты под serverless (в оригинальной библиотеке, кнопки переставали работать после перезапуска бота).
