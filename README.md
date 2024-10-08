# Тестовое задание по Golang

## Дано

GIT репозиторий с заготовкой приложения для подсчёта статистики по содержимому бинарного файла.

Реализован HTTP сервер выдающий страницу, через которую на вход можно подать любой файл.

Сервер, получив файл, должен побайтово проанализировать его содержимое и выдать статистику в формате JSON.

При каждом запросе необходимо собирать статистику запросов, сохраняя общее количество запросов, имена входных файлов и их размер в байтах.

Для запуска приложения можно использовать `Makefile` (из консоли) или `launch.json` (для VSCode).

Для выполнения запроса нужно запустить приложение и зайти в браузере на [тестовую страницу](http://localhost:8080).

## Задание

### 1. От ветки `master` создайте отдельную ветку

### 2. В пакете bitcounter реализуйте функционал подсчёта статистики числа бит в каждом байте, полученного файла

На вход должна подаваться структура

```go
type Input struct {
 Data []byte
}
```

в которой Data - это содержимое файла.

Необходимо пройтись по каждому байту файла.
Для каждого байта в файле нужно посчитать количество бит, выставленных в 1 - назовём его raisedBitQuantity. Например, для байта 0x03 - raisedBitQuantity == 2, для байта 0xFE - raisedBitQuantity == 7.
Далее нужно собрать статистику по байтам, имеющим одинаковое raisedBitQuantity.
Нужно посчитать общее количество таких байт - назовём это byteCount.
Нужно перечислить значения этих байт - например байты со значениями 0x01, 0x02, 0x04 все имеют raisedBitQuantity == 1.
По возможным значениям нужно составить таблицу (map), где ключ - это значение байта, а значение - список смещений байтов с таким значением.

Пример выходного результата приведён в [файле](./etc/output_example.json).

```json
{
 "byteStatistic": [     - массив со статистикой
  {
   "raisedBitQuantity": 0,  - статистика по байтам, у которых 0 бит выставлено в 1
   "byteCount": 0,    - таких байт 0 штук
   "values": {}    - тут пусто, так как таких байтов нет
  },
  {
   "raisedBitQuantity": 1,  - статистика по байтам, у которых 1 бит выставлен в 1
   "byteCount": 83,   - таких байтов в файле 83
   "values": {     - их возможные значения
    "4": [     - значение 4 (0000 0100)
     6,     - массив с перечислением смещений байтов от начала файла, которые имеют значение 4
     27,
     43,
     74,
     122,
     155
    ],
    "16": [     - значение 16 (0000 1000)
     22,     - массив с перечислением смещений байтов от начала файла, которые имеют значение 16
     34,
     190
    ],
    
    ...

   }
  },
  {
   "raisedBitQuantity": 2,  - статистика по байтам, у которых 2 бита выставлено в 1
   "byteCount": 57,   - таких байтов в файле 57
   "values": {     - их возможные значения
    "3": [     - значение 3 (0000 0011)
     92,     - массив с перечислением смещений байтов от начала файла, которые имеют значение 3
     140,
     194
    ],
    "5": [     - значение 5 (0000 0101)
     200,    - массив с перечислением смещений байтов от начала файла, которые имеют значение 5
     300,
     400
    ],
    
    ...

   }
  },

  ...

 ]
}

```

Дополните структуру

```go
type Result struct {
}

```

таким образом, чтобы из неё можно было бы сформировать JSON, описанный в выше.

Для ускорения обработки можно реализовать подсчёт числа бит через таблицу.

### 3. В пакете httphandlers дополните метод Process, чтобы в ответ на запрос выдавался требуемый JSON

### 4. В пакете httphandlers реализуйте новый метод Stat, который будет выдавать статистику выполненных запросов

Пример результата статистики запросов приведён в [файле](./etc/stat_example.json).

```json
{
 "requestCount": 4,     - общее число запросов
 "files": [       - массив с информацией о файлах
  {
   "filename": "README.md", - имя файла
   "size": 3978    - размер файла в байтах
  },
  {
   "filename": "Makefile",
   "size": 34
  },
  {
   "filename": "go.sum",
   "size": 8715
  },
  {
   "filename": "go.mod",
   "size": 1477
  }
 ]
}
```

Для сохранения статистики используйте in-memory хранилище данных.

### 5. В пакете httphandlers добавьте обработку ошибок и их вывод в лог

### 6. Сохраните изменения в отдельной ветке, запакуйте архив и пришлите его в качестве ответа
