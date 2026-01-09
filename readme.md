malanianchal last chat - https://chatgpt.com/share/6910f720-ff60-800c-93ec-d9fd68ac2ec2

continue from there 
Models implemented jwt 


Completed 
1. JWT Verification middleware is implemented but not yet used
2. Model and model methods are done
3. Login method impl -- user is logged in password verified and JWT token is generated, and token is returned
5. Auth abstraction for auth based APIs (for correct access)
6. Register and Login routes are working fine. 
7. When login is clicked you are getting a JWT token, now from here you need to impl how and where this is to be used. 



Pick From :
1. you need to make sure that token is mapped with correct user and only that user's info is shown to him, understand how this is done - this is frontend, on the backend you just need to check if there is an Auth token or not, and use JWT claims to check for roles.
2. Now I think the one part left if writing APIs with appropriate authorization checks and deploying and then frontend.
3. but before that let us deploy.


Frontend Side of Things :
1. JWT storage decision -b whether in a cookie from backend OR in local storage.
