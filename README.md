# send-query-result

[![Stable Version](https://img.shields.io/github/v/tag/anothrNick/send-query-result)](https://img.shields.io/github/v/tag/anothrNick/send-query-result)

Sends the text output of a Postgres query to a Slack channel via Incoming Web Hook (or any web hook). The query is run on a configurable interval.

For example, if you wanted to know the count of records for a specific table in Postgres (and you don't have metrics), you could run this container with the following query:

```sql
select count(id) as user_count from app_users;
```

The message will appear as standard SQL output:

```sh
 user_count 
------------
          5 
(1 row)
```

### Environment Variables

A configuration file would be better, but env vars will do for this project

|Key|Description|
|---|-----------|
|POSTGRES_USER|Username of the Postgres user|
|POSTGRES_PW|Password of the Postgres user|
|POSTGRES_HOST|Postgres hostname or ip address|
|POSTGRES_DB|Postgres database name|
|POSTGRES_SSL|Postgres [SSL Mode](https://www.postgresql.org/docs/9.1/libpq-ssl.html#LIBPQ-SSL-SSLMODE-STATEMENTS)|
|STAT_QUERY|The SQL query to be executed. Can result in single or multiple rows|
|STAT_INTERVAL|The interval, in minutes, to run the query and send the result to the web hook url|
|STAT_URL|The web hook URL to send the result (does not have to be Slack)|

### POST Body Format

This script is meant to send to a Slack Incoming Web Hook, so the request body will be formatted like so:

```json
{
    "text": "```<sql_output>```"
}
```

_Howevever, the sql output can be sent to any web hook_

### Dependencies

[github.com/shomali11/xsql](github.com/shomali11/xsql)
