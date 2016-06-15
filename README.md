Learning each step towards building a Golang coded app template  
Not using any 3rd party packages, e.g. routers, this so that all basics ca be learned     
TO DATE:  
	Added Home, About, Register & Signin pages  
	Added Member area & a display result page  
	Added password encryption - encrypted password saved to PostgresDb hosted by Heroku  
	Added dbconnect file with global variable & init for opening Postgres & Mongo databases (concurrency)  
    Added browser-side javascript registration input verification  
    Added Server-side (Golang) registration input verification - including check for duplicate e.mails & usernames  
    e.g. all essential info' requested has been input, all requested elements of the password included - numbers, length etc  
	Added e.mailed link as part of registration process  
	Added some middleware  
	That website pages look good on all types of device - making available on Heroku for testing   
	Added session management cookies  
	Added idle timeout after member signin  

TODO:  

e.mailed link for registration process time limited  
log out
login e.mail reminder for forgotten password/username at signin  
short number code sent to mobile phone for extra security    
Encrypt website HTTPS, SSL  (Heroku does this)  
Write tests  
Dockerize  

Current state can be seen on: http://blanktemplateapp.herokuapp.com
