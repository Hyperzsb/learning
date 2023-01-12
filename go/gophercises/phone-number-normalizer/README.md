# phone-number-normalizer

## PostgreSQL

I use following command to set up a PostgreSQL environment using Docker for this project:

```bash
$ docker run --rm -d --name postgres -p 5432:5432 -v "$PWD/.data/:/var/lib/postgresql/data" -e POSTGRES_PASSWORD=postgres -u postgres  postgres
```
