# accountapi-lib-form3

Hello, my name is Santiago Diaz and I am applying for a Software Engineer position. 
I have worked with Golang for about nine months.

I am going to list some considerations that I had to develop this assessment.

## Requirements
These requirements are based on the ones listed on the assessment repository:
1. Library should be a client library suitable for use in another software project.
2. Library have to implement the Create, Fetch, and Delete operations on the accounts resource.
3. Library should be well tested.

## Considerations
1. As it is a library, it should be decoupled from external resources, that is why I decided 
to allow components that use this library to configure backend host, port and version of 
   the account API.
   
2. The default scheme is `http` as compose version of the account API does not support `https`. Nonetheless,
it is possible to extend capabilities easily due to implementation.
   
3. In order to improve maintainability and flexibility when configuring this library, I implemented a `builder pattern
that allows us to custom:
   
   a. `http.Client`: I defined a default http.Client with 4 seconds of *timeout* (as per my experience this is good time) to
orientate our services to be more resilient. A component that uses this library can use a particular http.Client by configuring it
   through the builder.
   
   b. Port: if backend API has a particular port, it is possible to configure it.
   
   c. API version: I noticed that `v1` is the first version of the account API, however, it is possible to configure a different version.

4. Debugging is important, that is why I defined a mechanism to print information about request and response, however, it is important to mention that
it reduces performance a lot as per I discovered by executing a benchmark that you can find in benchmark folder. Enabling verbose log is only to debug.
   
5. As I am using a default `http.Client.Transport`, it implies that I am using a pool of connections, it improves performance, however,
I want to mention that a component that uses this library may configure timeout of alive connections, connections per host
   and other interesting things that can help in improving performance even more to a particular scenario.
   
6. To test this library and according to the restriction to use external libraries,
   I implemented `fakes` in order to emulate behaviour of some components. Nonetheless, it is possible to use
   other tools such as `gomock` to implement mocks and improve tests.
   
7. I added `// +build integration` to integration test files to allow us to decouple execution of tests,
as integration tests depend on external resources, and it is possible that those resources are failing when
   executing tests, it is possible to avoid using the flag `-tags` in a pipeline for example, to disable integration tests execution.