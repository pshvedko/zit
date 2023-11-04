## Тестовое задание от компании [ЗащитаИнфоТранс](https://www.z-it.ru/) 

https://gist.github.com/zemlya25/585ab3fb3b0704880f920728c7598beb

Нужно сделать REST API на Go, который принимает два `user_id` и выдаёт ответ - являются ли они дублем или нет. Дублем считается пара `user_id`, у которых хотя бы два раза совпадает ip адрес в логе соединений. Для каждого пользователя может быть много соединений, причём нормально если много из них с одного ip адреса. Никаких ограничений на уникальность в логе соединений нет. Пара одинаковых `user_id` всегда является дублем.
Лог соединений можно нагенерить рандомом, в бд или файле - неважно. Структура такая:

```sql
create table conn_log ( user_id bigint, ip_addr varchar(15), ts timestamp)
```

IP в формате IPv4. Кол-во записей хотя бы миллион. Ответ сервиса должен быть быстрым - меньше 30мс. Бонусом будет, если данные будут вычитываться не один раз, а будет поддерживаться актуальность - то есть при вставке новых записей в `conn_log` они будут учитываться в новых запросах к сервису.
Писать следует так, как писался бы реальный боевой сервис. Оцениваться будет с расчётом на это.

Пример:

В conn_log такие записи:
```csv
1, 127.0.0.1, 17:51:59
2, 127.0.0.1, 17:52:59
1, 127.0.0.1, 17:52:59
1, 127.0.0.2, 17:53:59
2, 127.0.0.2, 17:54:59
2, 127.0.0.3, 17:55:59
3, 127.0.0.3, 17:55:59
3, 127.0.0.1, 17:56:59
4, 127.0.0.1, 17:57:59
```

Выполняем GET запрос: http://localhost:12345/1/2
Ответ:

```json
{ "dupes": true }
```

Выполняем GET запрос: http://localhost:12345/1/3
Ответ:

```json
{ "dupes": false }
```

Выполняем GET запрос: http://localhost:12345/2/1
Ответ:

```json
{ "dupes": true }
```

Выполняем GET запрос: http://localhost:12345/2/3
Ответ:
```json
{ "dupes": true }
```

Выполняем GET запрос: http://localhost:12345/3/2
Ответ:
```json
{ "dupes": true }
```

Выполняем GET запрос: http://localhost:12345/1/4
Ответ:

```json
{ "dupes": false }
```

Выполняем GET запрос: http://localhost:12345/3/1
Ответ:

```json
{ "dupes": false}
```

## Результаты

```sql
INSERT INTO log SELECT 7, '127.0.0.1'::inet + 2*g AS ip FROM generate_series(0,1000000) AS g;
INSERT INTO log SELECT 8, '127.0.0.2'::inet + 2*g AS ip FROM generate_series(0,1000000) AS g;
INSERT INTO log SELECT 9, '127.0.0.1'::inet + 2*g AS ip FROM generate_series(0,1000000) AS g;
```

```
$ ab -n 1000 -c 10 127.0.0.1:8080/7/8
This is ApacheBench, Version 2.3 <$Revision: 1879490 $>
Copyright 1996 Adam Twiss, Zeus Technology Ltd, http://www.zeustech.net/
Licensed to The Apache Software Foundation, http://www.apache.org/

Benchmarking 127.0.0.1 (be patient)
Completed 100 requests
Completed 200 requests
Completed 300 requests
Completed 400 requests
Completed 500 requests
Completed 600 requests
Completed 700 requests
Completed 800 requests
Completed 900 requests
Completed 1000 requests
Finished 1000 requests


Server Software:        
Server Hostname:        127.0.0.1
Server Port:            8080

Document Path:          /7/8
Document Length:        16 bytes

Concurrency Level:      10
Time taken for tests:   2.665 seconds
Complete requests:      1000
Failed requests:        0
Total transferred:      124000 bytes
HTML transferred:       16000 bytes
Requests per second:    375.23 [#/sec] (mean)
Time per request:       26.651 [ms] (mean)
Time per request:       2.665 [ms] (mean, across all concurrent requests)
Transfer rate:          45.44 [Kbytes/sec] received

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0    0   0.1      0       1
Processing:    19   26   8.1     25     107
Waiting:       19   26   8.1     25     106
Total:         19   26   8.2     25     107

Percentage of the requests served within a certain time (ms)
  50%     25
  66%     30
  75%     31
  80%     31
  90%     32
  95%     32
  98%     32
  99%     73
 100%    107 (longest request)
```
