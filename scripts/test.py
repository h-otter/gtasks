import json

data = [
  {
    "title": "changed",
    "state": "Done",
    "labels": {
      "test-label": "hoge",
    },
    "description": "foobar",
    "dueDate": "2019-05-15",
    "dependsOn": [],
  },
  {
    "title": "ToDo",
    "state": "ToDo",
    "labels": {
      "test-label": "hoge",
    },
    "description": "foobar",
    "dueDate": "2019-05-15",
    "dependsOn": [
      0,
    ],
  },
  {
    "title": "testing",
    "state": "Doing",
    "labels": {
      "test-label": "hoge",
    },
    "description": "foobar",
    "dueDate": "2019-05-15",
    "dependsOn": [
      1,
    ],
  },
]

for d in data:
    print("curl http://localhost:12345/api/task/ -X POST -H 'Content-Type: application/json' -d '{}'".format(json.dumps(d)))
