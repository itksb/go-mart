# go-mart

Проект HTTP API "Накопительная система лояльности"

# Обновление шаблона

Чтобы иметь возможность получать обновления автотестов и других частей шаблона, выполните команду:

```
git remote add -m master template https://github.com/yandex-praktikum/go-musthave-diploma-tpl.git
```

Для обновления кода автотестов выполните команду:

```
git fetch template && git checkout template/master .github
```

Затем добавьте полученные изменения в свой репозиторий.
.github/workflows/gophermart.yml

# Миграции

Создать:

```shell
goose postgres "host=localhost user=user password=password dbname=postgres sslmode=disable" status
```

# CURL

curl -i -X POST -d '{"login": "legat", "password": "legat"}' http://localhost:8000/api/user/register
curl -i -X POST -d '{"login": "legat", "password": "legat"}' http://localhost:8000/api/user/login

curl -i -X POST -d 18 -H "Authorization: Bearer YOUR TOKEN HERE" http://localhost:8000/api/user/orders
curl -i  -H "Authorization: Bearer TOKEN" http://localhost:8000/api/user/withdrawals
