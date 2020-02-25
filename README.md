
Learning each step towards building a Golang coded app template
not using any 3rd party packages, e.g. routers, this so that all basics can be learned & worked through
TO DATE:

Added Home, About, Register & Signin pages
Added Member area
Added password encryption - encrypted password saved to PostgresDb hosted by Heroku
Added dbconnect file with global variable & init for opening Postgres & Mongo databases
Added browser-side javascript registration input verification
Added Server-side (Golang) registration input verification - including check for duplicate e.mails & usernames
  -  all essential info' requested has been input, all requested elements of the password included - numbers, length etc
Added e.mailed link as part of registration process
Added authorisation
That website pages look good on all types of device - making available on Heroku for testing
Added session management cookies
Added signin time limited
Encrypt website HTTPS, SSL (Heroku does this)
Added gzip site compression
Added suggestions page, with voting

TODO:

Feedback option for non-members to comment on site
browser caching issues - as per vary.com/pagespeed test
Render blocking CSS/JS needs sorting - as per vary.com/pagespeed test
e.mailed link for registration process time limited
Log out
Idle timeout
Login e.mail reminder for forgotten password/username at signin
Member suggestion page with ratings
Upgrade to paid Heroku service to include SSL on email messaging
Outlook365 tls cert
Write tests
Dockerize.

Current state can be seen on: gocloudcoding.com redirected to  https://golangwebapp.herokuapp.com
