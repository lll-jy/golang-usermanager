# Design

## Table of Content
1. [Architecture/Structure](#architecturestructure)
1. [Strategies & Algorithms](#strategies-&-algorithms)
1. [API Design](#api-design)
1. [Data Schema Design](#data-schema-design)
1. [Data Flow](#data-flow)
1. [Frameworks Used](#frameworks-used)
1. [References](#references)

-----

## Architecture/Structure

The app is a simple web application with the web app server connected to a database (see Figure below). HTTP requests 
are sent by the clients to the server and handled as queries to database. The web app server gets data from the database 
and send the information as HTTP response back to the client.

![architecture](diagrams/architecture.png)

## Strategies & Algorithms

1. A server socket port is started waiting for incoming connection requests.
1. A client socket is created to connect to the server socket as connection socket.
1. An HTTP request is sent from the client socket and read by the connection socket. Possible requests are seen in
   [API Design](#api-design) section. Socket programming as well as HTTP requests handling are written in GoLang.
1. A database query checking authentication is first sent to the database to ask whether the user has the right to 
   perform the desired action, which is in essence user login check. The database used is MySQL and 
   github.com/go-sql-driver/mysql is used to integrate with GoLang.
1. Successful login would allow the server to enter a session. The session is handled using cookies with the go modules
   github.com/gorilla/mux and github.com/gorilla/securecookie.
1. The server continue to handle the database query made based on the requests and data is fetched from the database 
   (to be elaborated separately below). Database queries are made based on the requests and data is fetched from the 
   database (to be elaborated separately below). 
    * In particular, the photo is fetched by the photo URL stored in the database, which is encrypted using the password 
      of the user. When a user logs in, a temporary photo file is created after decryption to the local device and 
      accessed by the web application server.
1. The connection socket generates corresponding HTTP responses based on the results from database and send to the 
   client socket.
1. The client socket read the response from connection socket and report to the user.
1. Repeat 3 to 8 to handle more requests from the user.
1. Close the current client socket.
1. Repeat 2 to 10 to handle requests from more users, possibly simultaneously.
1. Close the connection socket.

## API Design

| Method | Path             | Description                               |
|--------|------------------|-------------------------------------------|
| GET    | /                | Home page (log in)                        |
| GET    | /signup          | Sign up page                              |
| POST   | /signup          | Add a new user to the database            |
| POST   | /login           | Authenticate users with correct password  |
| POST   | /logout          | Log out                                   |
| GET    | /view            | User page (with username and profile)     |
| POST   | /delete          | Delete the user from database             |
| GET    | /edit            | Edit profile page                         |
| POST   | /edit            | Upgrade user profile in the database      |
| GET    | /reset           | Reset user info page                      |
| POST   | /reset           | Upgrade user info in the database         |
| POST   | /discard         | Cancel on the edit page                   |
| POST   | /remove          | Remove profile photo                      |
* Note: user info refers to username and password, and user profile refers to the user profile photo and nickname.

### GET `/`
Display the homepage to the user from the HTML template that allows the user to choose to provide the username and 
password to log in.

### GET `/signup`
Display the sign up page to the user from the HTML template that allows the user to provide the username and password 
that the user wants to sign up for.

### POST `/signup`
Insert a new user to the database with the username and password specified. 

Successful sign up would redirect to `GET /edit` to ask the user to provide initial profile.

### POST `/login`
Checks validity of the user's identity.

Successful login would result to `GET /view`.

Unrecognized but valid username log in would redirect to `GET /signup`. The user can also choose to sign up directly.

### POST `/logout`
User logs out and go back to the home page. Remove local photo file.

### GET `/view`
Display the profile page to the user from the HTML template that allows the user to view the username, nickname, and 
profile photo, and choose to edit profile, reset password, or delete user.

### POST `/delete`
Delete the user from the database (including the photo in the file base) and redirect to the homepage.

### GET `/edit`
Display the edit profile page to the user from the HTML template that allows the user to (re-)upload profile photo 
and/or change a nickname. Remove unused photos when necessary.

### POST `/edit`
Save the photo uploaded by the user and get its address. Then save the changes to the user profile to the database and 
redirect to `GET /view`.

### GET `/reset`
Display the reset user info page to the user from the HTML template that allows the user to reset the username and/or 
password.

### POST `/reset`
Save the new user info and update the database accordingly. Then redirect to `GET /view`.

### POST `/discard`
Choosing 'Cancel' on the edit page, to discard the photo uploaded and use the original photo instead.

### POST `/remove`
Remove the profile photo for the user (and the file in the actual file base is to be deleted when the changes
are saved, i.e. when the 'Save' button on the `/edit` page is clicked).

## Data Schema Design

Only on table is needed, so ER diagram is not needed. The table creation query in SQL is as follows.

```
CREATE TABLE users (
    username    VARCHAR(20) PRIMARY KEY,
    password    VARCHAR(100) NOT NULL,
    photo       VARCHAR(50),
    nickanme    VARCHAR(30) COLLATE utf8mb4_unicode_ci
);
```

To integrate UTF8 encoding, before creating the table, the following query is run.

```
ALTER DATABASE EntryTask CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci;
```

### Username
The username is the primary key that uniquely identifies each user. The username should be a 4-20 character string 
containing numbers, letters, hyphens ("-"), and underscores ("_") only, and the first character must be a letter. 
Duplicate username is not allowed.

### Password
The password should not be null, and the password that user provides should be a 4-20 character string. The original 
password is hashed before it is saved to the database.

### Photo
The photos are stored as files at some particular place with a path known, and this path is the content of the photo 
attribute in the database. This method of photo handling is preferred over storing the photo directly as `VARBINARY` in 
the database because of the following reasons:
1. Photos are usually big files, which may use too much space in the database.
1. If the photo is directly stored in the database, anyone with access to the database can access the photo, whereas 
   storing the photo files separately makes separate authentication and access to the photos possible, so it potentially 
   protects the user's information in a sense.

### Nickname
The nickname can be any UTF-8 string of no more than 30 characters.

## Data Flow

The dataflow is seen as follows.

![data flow](diagrams/data_flow.png)

## Frameworks Used

| Framework | Usage |
|-----------|-------|
| `crypto/aes`, `crypto/cipher`, `crypto/rand` | file encryption and decryption framework |
| `database/sql` | database helper framework |
| `errors`, `fmt`, `io`, `io/ioutil`, `log`, `os`, `reflect`, `regexp`, `strings`, `sync`, `testing` | fundamental Go in-built functions |
| `net/http` | handle HTTP clients, requests, etc. |
| `github.com/go-sql-driver/mysql` | Go MySQL driver |
| `github.com/gorilla/mux` | router framework |
| `github.com/gorilla/securecookie` | cookie encoding and decoding framework |
| `golang.org/x/crypto/bcrypt`  | password hashing and checking |
| `google.golang.org/protobuf/proto`, `google.golang.org/protobuf/reflect/protoreflect`, `google.golang.org/protobuf/runtime/protoimpl`| protobuf framework |
| `text/template` | HTML template parsing framework |

## References

1. GoLang API: https://golang.org (given reference)
1. GoLang web app: 
   1. Tutorial: https://golang.org/doc/articles/wiki/ (given reference)
   1. Structure and assets file directory: https://www.sohamkamani.com/golang/how-to-build-a-web-application/
1. Go MySQL driver: 
   1. Source repository: https://github.com/go-sql-driver/mysql (given reference)
   1. Examples: https://github.com/go-sql-driver/mysql/wiki/Examples
1. Protobuf: https://developers.google.com/protocol-buffers (given reference)
1. Cookie encoding and decoding for sessions: https://gist.github.com/mschoebel/9398202
1. File encryption and decryption: https://www.melvinvivas.com/how-to-encrypt-and-decrypt-data-using-aes/
1. Go Profiling:
   1. Introduction: https://blog.golang.org/pprof (given reference)
   1. Tutorial: https://code.tutsplus.com/tutorials/make-your-go-programs-lightning-fast-with-profiling--cms-29809
