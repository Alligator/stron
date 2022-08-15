# stron

show the structure of JSON in a compact, gron-esque format.

stron reads JSON from stdin or a file, and prints the distinct paths through it.
it can print example values too, with `-v`.

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
