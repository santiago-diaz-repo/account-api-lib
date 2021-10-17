# accountapi-lib-form3

This implementation was built by Santiago Diaz.

Here are my considerations while developing the assessment.

## Requirements
These requirements are based on the ones listed on the assessment repository:
1. Library should be a client library suitable for use in another software project.
2. Library have to implement the Create, Fetch, and Delete operations on the accounts resource.
3. Library should be well tested.
4. Library should not implement client-side validation
5. Library should not use any external library, it should be `go get` friendly

## Considerations
1. As it is a library, it should be decoupled from external resources, that is why I decided 
to allow components that use this library to configure backend host, port and version of 
   the account API. I implemented a configuration with default values, however, they can be easily modified.
   
2. The default scheme is `http` as compose version of the account API does not support `https`. Nonetheless,
it is possible to extend capabilities easily due to library implementation.
   
3. In order to improve maintainability and flexibility when configuring this library, I implemented a `builder pattern`
that allows us to custom:
   
   a. `http.Client`: I defined a default http.Client with 4 seconds of *timeout* (as per my experience this is a good time) to
orientate our services to be more resilient. A component that uses this library can use a particular http.Client by configuring it
   through the builder.
   
   b. `Host`: Default implementation points out to `localhost` because everything was tested on a local laptop, but the host can
be modified by using the `WithHost()` method.
   
   c. `Port`: if the backend API has a particular port, it is possible to configure it by invoking the `WithPort()` method.
   
   d. `API version`: I noticed that `v1` is the first version of the account API, however, it is possible to configure a different version by invoking the `WithAPIVersion()` method.

4. Debugging is important, that is why I defined a mechanism to print information about request and response, however, it is important to mention that
Enabling logging verbose by invoking the `Verbose()`method  reduces performance up to 90%. I implemented a benchmark to show this impact. It can be found in the *benchmark* folder.
   
5. As I am using a default `http.Client.Transport`, it implies that I am using a pool of connections, it helps in improving performance, however,
I want to mention that a component that uses this library may configure timeout of *alive connections*, *connections per host*
   and other interesting things that can help in improving performance even more in a particular scenario.
   
6. To test this library and according to the restriction to use external libraries,
   I implemented `fakes` in order to emulate the behaviour of some components. Nonetheless, it is possible to use
   other tools such as `gomock` to implement mocks and improve tests configuration and execution without writing too much code.
   
7. I added `// +build integration` to integration tests to allow us to decouple the execution of tests.
As integration tests depend on external resources, those resources may be failing when
   executing tests, that is why integration tests can be disabled by removing the flag `-tags` from the execution command `go test -tags integration ./...`.
   By doing so we can prevent development delays or even failures in our *CI/CD* pipeline.
   
8. Integration tests are configured with a fixed `UUID` as they may be executed several times, and it can lead to increased 
storage in the account-API database. Nonetheless, it is possible to use an external library such as *google/uuid*.
   
9. I defined a custom error struct, when there is an error either internally or when the account API is invoked, this library
returns the custom error struct. Some status code can be found in the *Specification of errors* section.

## Instruction to use library
According to requirements, this library can be implemented in any project, to do that you can follow the following steps:
1. Create a configuration as follows:
```
config := configuration.NewDefaultConfigBuilder().
		WithHost("account-api-host").
		WithPort("8080").
		WithAPIVersion("v2").
		Build()
```

2. Create an AccountService by sending the `config`object of step 1:
```
accountService := NewAccountService(&config)
```

3. Create a request object, it depends on what you want to execute, `Create`, `Delete` or `Fetch`. For simplicity, 
let's create a `DeleteRequest` object:
   
   ```
    req := models.DeleteRequest{ 
       AccountId: "12ab1977-6894-4d82-9968-4044df675fd9",
       Version:   0}
   ```
   
4. Invoke method that you need to execute against the account API. Following the step 3, let's implement `Delete` operation.
   You should validate that there was no error, if so, you can get the statusCode and message. You can get the response when there was no error:
```
res, err := subject.DeleteAccount(&input)
if err != nil {
	acctErr := err.(*error_handling.AccountError)
	statusCode := acctErr.GetCode()
	errMsg := acctErr.GetMessage()
	fmt.Printf("%d - %s",statusCode, errMsg)
	return
}

fmt.Printf("%d",res.StatusCode)

```

## Considerations about docker-compose
1. I created a Dockerfile in which I am executing unit and integrations tests with verbose output.
2. I took your docker-compose, add my image, which is built and executed when executing `docker-compose up`,  and versioned it in my repository.
3. To run your docker-compose I added the Database creation script to my repository. 
4. To avoid errors related to network when executing `docker-compose up`, I am defining a network called *acct-lib-santiago*.


## Specification of errors

| Code | Description |
|------|-------------|
|1| failed marshalling request|
|2| failed creating request|
|3| failed invoking backend|
|4| failed reading response body|
|5| failed decoding error response|
|6| failed decoding response|
|404| Resource does not exist|
|400| You sent something wrong to the account API|
|409| There was a conflict when trying to create resource, it may already exist|
