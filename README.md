# example-dating-apps
### How to Run ? 
- Build the docker image using ```docker build ./ -t dating-apps```
- Run the docker compose using ```docker-compose up -d```
- Now the application has start on port ```8843```

# Database
- I always using uuid_key on every table that used for the unique id of the table that will received by user. I learn it based on OWASP top 10 ```Broken Access Control``` that told that we need to use the unique UUID that used for get detailing data cause it so hard to brute forced.

# Docker
- I'm using docker and docker compose to dockerize my app. it help me to move to other environment and the application still working fine. 
- My docker-compose contain postgresql, redis, and the application.

# Code Descriptions:
## token
- I generate Token with JWT Format. 

## authentication
- This Package will contain all function that used on custom_endpoint to validate the incoming request. On this case, every incoming message has been wrapped on the custom_endpoint and every Authentication Function must have structure like this:
    ```go
    func(
        ctx context.Context,
        header map[string]string,
    ) error
    ``` 
- I've create example implementation of user_authentication that can be used to validate the incoming token that will checked to Redis and the user data will parsed to the custom_context and it will be sent to the service layer.

## bundles & i18n
- This package is used to read the structure of i18n messaging. i18n is used for dictionary for multiple language response.
- i18n is used to save every known message, such as error message and constant field name. So the response api can converting the ERROR_CODE into the message. 

## config
- Based on Twelve Factor Apps, we need to separate all of secret and everything that possible to configured by ops on different environment. For the example, i want key of jwt token is different key on the production and staging, by changing the environment i can solve that problem.

## custom_contex
- This is the context that will store user data based on authentication layer on the controller. So on the service, i can make sure that user has autenticated and authorized to access and i can get the user data.

## custom_endpoint
- I creating the custom endpoint that wrap some repetitive action. for the example marshaling incoming request body, parsing to DTO, validate the incoming request, authentication and authorization, Standardization Output response on error and success so it help us on read the response.
- This is how i wrapping the router. it used to collecting all data before register it on the gorilla Mux. I use it for generating swagger api documentation
    ```
    NewHandleFuncParam(
        path string,
        f func(http.ResponseWriter, *http.Request),
        method ...string,
    )
    ```
- This is how i wrap the default golang routing func(http.ResponseWriter, *http.Request). I want all of request is clean and validated before processed by the service layer. I need to parse the Service cause i need to read the DTO that will be used, the FunctionServe is the function that will be served and the ServerAccessValitor is the function that will used on the authentication layer
    ```
    NewWarpServiceParam(
        service Services,
        serve FunctionServe,
        cv ServerAccessValidator,
    ) 
    ```

## custom_error
- This used for generating standardization error response that i can create the error message that contain the Status code, Error code that will converted on i18n, adding the caused by for logging etc.

## password 
- This repository contain function that used to encrypt/hash the password. 

## repository
- This is the implementation of table on database that implemented on golang structure.

## server_attribute
- Declare all of the needs, such as third party app, controller, service, and dao.

## sql_migrations
- List of files that will used executed to creating/updating database that suit to the ERD
