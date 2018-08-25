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

Game library
  scoring tools
  list of games
  mark want to play (unmarked when played)
  mark favorites (stays marked)
  mark as played
  save wishlist
  respond to wishlist

Alerts
  Expiration
  ExpirationActions (e.g. delete or archive alert, lock account)
alerttypes
  one time (some system alerts should never be sent an additional time)
  user dismissable (usually true, but some only by event, such as verifying an email)
  chat
  tip
  comment reply
  post reply
  mention
  tenant comment
  tenant action
  site admin (flagged post, server restart)
  
photo album

file manager
  file access: owner only, authenticated (or maybe this is "user" membership), group only
  tags (maybe use bleve to return files by tag), bleve should also index file ext, upload user disp name and email, date range
  feeds photo albums, profile pictures, images for blogging, downloads, and amazon summary project

forum

blog
  stores pre-calculated summary snippets of last X number of posts
post

auth
  site registration (non-oauth)
  email verification
  pw retry lockout (per ip? both ip and overall?)
  two factor (mandatory for admins)
  session timeout
  heartbeat via websockets

pinned content (pin a vault, project, page, etc)

chat

projects: small projects that don't have their own domain but that I want to put behind the auth and data storage of DWN
  amazon storage: summary data of one's Amazon history that parses Amazon order CSV, removes duplicates, and shows history (table and charts)
  journal: micro blogging only visible to you
  personal crm: business card and contact info storage
  friend schedule tracker: show if a friend with a weird shift schedule is working
  