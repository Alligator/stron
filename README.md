# stron

show the structure of JSON.

stron reads JSON from stdin and prints all the paths through it, with example values if you'd like.

```
$ cat example.json
{
  "users": [
    {
      "name": "alligator",
      "repos": [
        {
          "name": "ely",
          "url": "https://github.com/alligator/ely"
        }
      ]
    },
    {
      "name": "sponge",
      "repos": [
        {
          "name": "slate2d",
          "url": "https://github.com/sponge/slate2d"
        }
      ]
    }
  ]
}

$ stron example.json
.users[].name
.users[].repos[].name
.users[].repos[].url

$ stron -v example.json
.users[].name = "alligator"
.users[].repos[].name = "ely"
.users[].repos[].url = "https://github.com/alligator/ely"
```
