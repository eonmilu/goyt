# goyt [![Build status](https://travis-ci.org/eonmilu/goyt.svg?branch=master)](https://travis-ci.org/eonmilu/goyt/)

Library used by the [Your Time Extension](https://github.com/eonmilu/yourtime).

## Installation

Run:
`go get github.com/eonmilu/goyt`

## Requirements

You will need a database with the tables `users`, `reports` and `timemarks` as following:

```sql
                                                       Table "public.users"
   Column   |          Type           |                     Modifiers                      | Storage  | Stats target | Description
------------+-------------------------+----------------------------------------------------+----------+--------------+-------------
 id         | integer                 | not null default nextval('users_id_seq'::regclass) | plain    |              |
 token      | character varying       |                                                    | extended |              |
 identifier | character varying       | not null                                           | extended |              |
 username   | character varying(256)  |                                                    | extended |              |
 url        | character varying(2048) |                                                    | extended |              |
 upvotes    | integer[]               |                                                    | extended |              |
 downvotes  | integer[]               |                                                    | extended |              |
Indexes:
    "users_email_key" UNIQUE CONSTRAINT, btree (identifier)
```

```sql
                                                         Table "public.reports"
   Column   |            Type             |                      Modifiers                       | Storage  | Stats target | Description
------------+-----------------------------+------------------------------------------------------+----------+--------------+-------------
 id         | integer                     | not null default nextval('reports_id_seq'::regclass) | plain    |              |
 timemarkid | character varying(16)       |                                                      | extended |              |
 ip         | inet                        |                                                      | main     |              |
 content    | character varying(256)      |                                                      | extended |              |
 date       | timestamp without time zone |                                                      | plain    |              |
 reviewed   | boolean                     |                                                      | plain    |              |
```

```sql
                                                         Table "public.timemarks"
  Column   |            Type             |                       Modifiers                        | Storage  | Stats target | Description
-----------+-----------------------------+--------------------------------------------------------+----------+--------------+-------------
 id        | integer                     | not null default nextval('timemarks_id_seq'::regclass) | plain    |              |
 videoid   | character varying(64)       | not null                                               | extended |              |
 ip        | cidr                        | not null                                               | main     |              |
 timemark  | integer                     | not null                                               | plain    |              |
 content   | character varying(256)      | not null                                               | extended |              |
 votes     | integer                     | default 0                                              | plain    |              |
 author    | integer                     | default 1                                              | plain    |              |
 approved  | boolean                     | default false                                          | plain    |              |
 timestamp | timestamp without time zone | default now()                                          | plain    |              |
 reports   | integer[]                   |                                                        | extended |              |
Indexes:
    "timemarks_id_key" UNIQUE CONSTRAINT, btree (id)
```

## Use

Import the goyt package, then declare a YourTime variable:

```go
yourtime = goyt.YourTime{
    AuthTokenURL:   YOUR_GOOGLE_AUTH_TOKEN_URL,
    GoogleClientID: YOUR_GOOGLE_AUTH_CLIENT_ID,
    DB:             YOUR_DATABASE,
}
```

Then you can specify the path of the handler functions, which are wrapped by the createUsers middleware:

```go
    r.HandleFunc("/yourtime/search", yourtime.CreateUsers(yourtime.Search))
    r.HandleFunc("/yourtime/insert", yourtime.CreateUsers(yourtime.Insert))
    r.HandleFunc("/yourtime/votes", yourtime.CreateUsers(yourtime.Votes))
    r.HandleFunc("/yourtime/auth/validate", yourtime.CreateUsers(yourtime.ValidateAuth))
    r.HandleFunc("/yourtime/auth/remove", yourtime.CreateUsers(yourtime.RemoveAuth))
```
