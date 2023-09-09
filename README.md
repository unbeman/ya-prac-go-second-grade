# ya-prac-go-second-grade

## Менеджер паролей

todo: написать про запуск и билд

### Что использовалось

1. Протокол взаимодействия - gRPC + TLS + JWT
2. Хранилище - postgresql на сервере и sqlite на клиенте
3. Функции шифрования и хэширования из стандартной либы crypt
4. TOTP для двухфакторной авторизации (github.com/pquerna/otp/totp)
5. Терминальное приложение на базе promt. (хочу переписть на )

## Детали реализации

todo: про устройство клиента и сервера на го

### Вход

После старта работы клиента, пользователю необходимо пройти авторизацию с использованием его логина и мастер пароля.
Чтобы пройти онлайн авторизацию (офлайн вход на данный момент не реализован) + получить токен доступа (который будет необходим для последующих запросов на сервер) клиент вычисляет мастер ключ длиной 64 байт по алгоритму Scrypt из мастер пароля и логина пользователя. 
Перед отправкой на сервер этот ключ хэшируется по тому же алгоритму Scrypt. Полученный хэш отправляется вместе с логином пользователя на сервер.
(На сервере захешированный мастер ключ в момент регистрации был повторно захеширован алгоритмом Bcrypt и положен в базу.) Соответственно в момент логина на сервере сравниваются дважды захэшированные ключи. 

<img src="https://github.com/unbeman/ya-prac-go-second-grade/assets/16476703/b7dcf80e-da68-4063-967e-605e124087a8" height="500">

### Хранение и синхронизация

Чтобы безопасно хранить доверенные секреты пользователей, было принято решение их шифровать.
Все операции по шифрованию и дешифрованию происходят на клиенте. Шифрование происходит в момент создания секрета, а дешифрование перед тем как отобразить секрет пользователю.
Шифрование проиходит c помощью блочного алгоритма AES-256. Для него необходим ключ длиной 32 байта. Этот ключ получается путем "расширения" мастер ключа функцией HKDF (HMAC Key Derivation Function) с алгоритмом хеширования SHA-256.

<img src="https://github.com/unbeman/ya-prac-go-second-grade/assets/16476703/ea4e9b0b-2d11-454b-a83d-cc3d4cf60d56" height="500">

## Команды клиента

### `register {user_login} {master_password}`

Создает нового пользователя с заданным логином и паролем, если логин не занят.

Пример работы:

```
>>>register krolick malenkayaNoraVLesu
```

```
>>>register unbeman master
error occurred:  login already exists
```

---

### `login {user_login} {master_password}`

Авторизует пользователя по логину и паролю.
Если включена двухфакторная аутентификация, то попросит выполнить валидацию кода.
Если отключена, то сразу запускает фоновую синхронизацию данных.

Пример работы:

Если 2FA не включена

```
>>>login unbeman correctPassword
successfully logged in.
```

Если 2FA включена

```
>>>login unbeman master
please validate 2fa code
```

Если введен неправильный пароль

```
>>>login unbeman incorrectPassword
error occurred:  rpc error: code = InvalidArgument desc = invalid login or password
```

---

### `genTOTP`

Создает QR-код для подключения двухфакторной аутентификации.
Дополнительно отображает ссылку и ключ для ручного подключения.

```
>>>genTOTP
█████████████████████████████████████████████████
█████████████████████████████████████████████████
████ ▄▄▄▄▄ ████▀ █  █▀▄ █▀▄ ▀█▄ █▀█ ▀█ ▄▄▄▄▄ ████
████ █   █ █▄█ █▀▀▀ ███ ▄▀ ▀▀█▀▀▀█   █ █   █ ████
████ █▄▄▄█ ██ ▄▄█▀ ▀█ █▀▀▀▀▀▀   ██▄ ██ █▄▄▄█ ████
████▄▄▄▄▄▄▄█ █ █▄▀▄█▄▀ ▀▄▀ ▀ ▀ █ ▀ █ █▄▄▄▄▄▄▄████
████▄▄ █  ▄▄▄▀▀ ▄ █▄▄ ▀▄▀ ▀█  ██  ██▀▄ ██▀▄█▀████
████ ▀▀▀ ▀▄▄▀ ▄▀█ █▀█▀▀█▀ ▀██▀▀▀█▀▀▀██ ▄▀█▄▄ ████
████▄█ ▄▀▀▄ █   ▄█▀▀▀▀▀ ▀▀█  ▄██▄▄▀█▄  ▄████ ████
████▀█  ▄▄▄▄▀▄█▀▀▀ ▀▀▀█▀█▀▄▀ ▀ ▀▄  █▄ ▄ █▄▄ ▄████
████ ▀▀█▄▄▄▄█▄ █▄ ███ █▄ ▄██  ██  ██   ▀█▀█▄ ████
████▀▀ ▄█ ▄ ██▄█ ▀ ▀█ ▄█▀▀▄██  █▄ ▄▀▄█▀▄▀▄▄▄ ████
████▄  █▄▄▄ ▄  ▀▀  ▀ ▄ ▀ ▀█▀▀▄▀▄█ ██▀  ▄█▄██ ████
████▄▀ ▀▀▄▄▀█▄ ▀▀ ▀█▄▀▄▀▄▀▀█▀ ███▀▀▀▀▄▀ ▄▄▄▄ ████
████▀▄ ▄▀▄▄▄ █▀▄█ █   █▄ ▄▀▄▀▀███▄██▀▄ ▀█▄█▄▀████
██████▀█  ▄ ▄ █▀▀▀▀▀▄ █▀ ▀ ██ ▄█▄▀ █▀▄▄▀▀▄▄▄▄████
█████▄█ ▀▀▄ ██   ██ ▀▄  ▀▀▀▀  ███ █▄   ██ █▀▀████
████▄▀█▀██▄█ ▄▄▀█▀ ▀▄▀▀█  ▄▀▀▀▄▀▀▀▀▀▄ ▄▀▄▀█  ████
████▄███▄█▄▄▀█▄██ ▀    ▄   ▄▀ ▀▄▀▄▀▀ ▄▄▄ ▄█  ████
████ ▄▄▄▄▄ █ ▀▀██▀███ ▀▀█▀ ▀█▀▀▀▀ ▀  █▄█ █▄  ████
████ █   █ █▄▀█▀ ▀█▀ ██▀ █▀██ ▀██ █▀▄▄▄▄ █▀ ▄████
████ █▄▄▄█ █ ▀█▀▀▀██ ▀▀▀▄ ▀▀ ▀▄▀ ▀▄█ ▄█▀▄▄▄▀█████
████▄▄▄▄▄▄▄█▄█▄▄█▄██████▄▄██▄▄██▄▄█▄▄▄▄██▄█▄▄████
█████████████████████████████████████████████████
▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀
Secret key:  OTEPWHIQPNXIYA5R67FHZEZGFQ5DMAMU
URL:  otpauth://totp/passkeeper:unbeman?algorithm=SHA1&digits=6&issuer=passkeeper&period=30&secret=OTEPWHIQPNXIYA5R67FHZEZGFQ5DMAMU
please verify token to enable 2FA
```

---

### `verifyTOTP {code}`

Проверяет что подключение двухфакторной аутентификации выполнено успешно.

```
>>>verifyTOTP 881121
successful verified
```

---


### `validateTOTP {code}`

Проверяет код от аутентификатора. Необходимо вводить каждый раз после `login`, если включена 2FA.
Запускает фоновую синхронизацию данных, в случае успеха.

```
>>>validateTOTP 529468
successful validated
successful logged in
```

```
>>>validateTOTP 123456
error occurred:  rpc error: code = InvalidArgument desc = invalid otp token or user id
```

---

### `disableTOTP`

Отключает двухфакторку.
```
>>>disableTOTP
successful disabled 2FA
```
---

### `get-all`

Выводит список всех сохраненных учетных данных пользователя не раскрывая секрет.

```
>>>get-all
ID                                      Type    Info
e5a541b5-641c-484d-b4bf-5047244d5003    note    Title: MyNote
8934b85e-f878-4b91-8b2a-33e25efed44f    bank    Card Name: Tinkoff
08dc55c8-2755-4d53-ba01-0a1862d00366    login   Site: abba.com  Login: lucia
b6e8257d-6b39-4929-9d58-de2741ec9582    login   Site: game.org  Login: G@LG@D.T
a9b45a3c-544a-4c10-8050-ac66e6a8f315    login   Site: google.com  Login: vasya
9d8da783-399b-477d-a6f8-fe685b56ffbe    login   Site: haha.com  Login: ann
b53505e0-45fd-4f58-b9c0-81aa047ff621    login   Site: wow.ru    Login: orkadiy
Get row by id to show secret
```
---

### `search {search_string}`

Выводит список найденных секретов по одному ключевому слову.
Поиск происходит по метадате секрета. Для секрета типа `login` это сайт + логин, для `note` и `bank` его название.

```
>>>search tink
found 1 records 
ID                                      Type    Info
8934b85e-f878-4b91-8b2a-33e25efed44f    bank    Card Name: Tinkoff

```

---

### `get {uuid}`

Выводит расшифрованный секрет пользователя.

```
>>>get 9d8da783-399b-477d-a6f8-fe685b56ffbe
Site: haha.com
Login: ann
Password: secretpass
```

```
>>>get e5a541b5-641c-484d-b4bf-5047244d5003
Title: MyNote
Note:
Long insteresting text of my first note.
```

```
>>>get 17dd7499-881a-43ee-8054-7cd87d915362
Card Name: Tinkoff
Number: 1234123412341234
Exp: 05/32
CVC: 123
```

---

### `create {type} ...`

Создает секрет.

Для секрета типа `login` команда `create login {site} {email} {password}`

```
>>>create login ya.ru vasya.pupkin@ya.ru VasinSuperSecretParol
successfully created
```

Для секрета типа `bank` команда `create bank {title} {card_number} {expired_at} {cvc}`

```
>>>create bank VTB 1234123412341234 01/24 123
successfully created
```

Для секрета типа `note` команда `create note {title} ...`

```
>>>create note MainSecret Don't forget to call Mom.
successfully created
```

---

### `edit {uuid} ...`

Редактирует секрет.

```
>>>edit 1c09c08c-034c-4b74-aba7-9bf5acadc490 ya.ru vasya.pupkin@ya.ru UpdateVasinSuperSecretParol
successfully updated
```

---

### `exit`

Завершает работу клиента.
Останавливает синхронизацию.

```
>>>exit
goodbye!
```

---

## API

### Регистрация`/pass_keeper.AuthService/Register`

---

### Вход `/pass_keeper.AuthService/Login`

---

### Генерация ключа 2FA `/pass_keeper.OtpService/OTPGenerate`

---

### Подтверждение использования 2FA `/pass_keeper.OtpService/OTPVerify`

---

### Проверка одноразового кода `/pass_keeper.OtpService/OTPValidate`

---

### Отключение 2FA `/pass_keeper.OtpService/OTPDisable`

---

### Сохранение секретов `/pass_keeper.SyncService/Save`

---

### Получение секретов `/pass_keeper.SyncService/Load`

---

    


