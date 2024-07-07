# TaskManager

Task manager is a test problem for "Effective Mobile". It is supposed that the service is used by authorized persons. It doesn't have any authorization mechanism and isn't supposed to have one.

The primary task of the service is to make a simple way of starting and stopping a countdown for any task and then getting the tasks the user was working on over a period of time.

## Installation

1. `git clone https://github.com/noname0443/TaskManager.git`
2. Change .env variables
3. Install Docker
4. `make run` in the project root

## API

You can check more API details by using swagger at "/swagger/index.html".

### /api/v1/users \[GET\]

Get users' data.

Request parameters:
- offset: Pagination offset
- limit: Pagination limit
- filters: filter in format: "filter1=value1,filter2=value2,..."

Response:
- User model
- Error

### /api/v1/users \[POST\]

Create user. It sends a request to EXTERNAL_WEBSERIVCE_URL to get detailed information about the user.

Request parameters:
- passportNumber

Response:
- int: User ID
- Error

### /api/v1/users/{userId} \[GET\]

Get the user's tasks.

Request parameters:
- to: (optional) strint in format YYYY-MM-DDThh:mm:ss.SSSZ
- from: (optional) strint in format YYYY-MM-DDThh:mm:ss.SSSZ
- offset: Pagination offset
- limit: Pagination limit
- userId: User ID

Response:
- Task model
- Error

### /api/v1/users/{userId} \[PUT\]

Update user's data.

Request parameters:
- userData

Response:
- "ok"
- Error

### /api/v1/users/{userId} \[DELETE\]

Request parameters:
- userId: User ID

Response:
- "ok"
- Error

### /api/v1/tasks/{taskId} \[PUT\]

Request parameters:
- taskId: Task ID
- status: New status of the task

Response:
- "ok"
- Error

### /api/v1/tasks/ \[POST\]

Request parameters:
- userId: User ID
- description

Response:
- int: taskId
- Error