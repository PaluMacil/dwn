Waiting long time after starting oauth will cause response of

Code exchange failed with 'oauth2: cannot fetch token: 400 Bad Request
Response: {
  "error" : "invalid_grant",
  "error_description" : "Bad Request"
}'


Badger repro:
Start db, close
start db, panic
start db, truncate, panic
now badger should complain that pointer is to beginning of database, which means whole thing is corrupt