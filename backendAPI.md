# Register

| parameter       | description               |
| --------------- | ------------------------- |
| Name            | User                      |
| Email           | User EmailID              |
| Password        | Password for user account |
| Profile Picture | User picture              |

# Login

| parameter | description               |
| --------- | ------------------------- |
| Email     | User EmailID              |
| Password  | Password for user account |

# Update

| parameter       | description               |
| --------------- | ------------------------- |
| Name            | User                      |
| Email           | User EmailID              |
| Password        | Password for user account |
| Profile Picture | User picture              |

# Delete

| parameter | description                                 |
| --------- | ------------------------------------------- |
| JWT Token | JWT token of the user details to be deleted |

# Message

| parameter | description           |
| --------- | --------------------- |
| ID        | Message ID            |
| Send to   | The recipient details |
| From      | Sender Details        |

## loginAPI

Users can login to the application  
**Method**: POST  
**URL**: http://localhost:8000/login <br />
**Auth required**: Need JWT token  
If login error

```
Status Code: 401
Response:
{
  "message": "Incorrect Credentials",
}
```

If login success

```
 Status Code: 200
 Response:
{
 "message": "User Login Successful",
 "token": <JWT Token>,
 "user": {
   "email": "testuser@gmail.com",
   "id": 1,
   "name": "Test User",
   "password": <Encrypted Password>,
   "profilepic": <Image as string>,
 }
}
```

## RegisterAPI

Users can register to the application  
**Method**:POST  
**URL**:http://localhost:8000/register  
**Auth required**: Need token  
If register error

```
{
  "message": "Invalid user details",
  "code":    400,
}
```

If add register success

```
{
 "message": "All is fine",
 "code":    201,
}
```

## UpdateAPI

Users can update their information in the application  
**Method**:POST  
**URL**:http://localhost:8000/user/update  
**Auth required**: Need token  
If update error

```
{
  "message": "Invalid user details",
  "code":    400,
}
```

If update success

```
{
 "message": "All is fine",
 "code":    201,
}
```

## DeleteAPI

Users can delete their details in the application  
**Method**:POST  
**URL**:http://localhost:8000/user/delete  
**Auth required**: Need token  
If delete error

```
{
  "message": "Invalid user details",
  "code":    400,
}
```

If delete success

```
{
 "message": "All is fine",
 "code":    201,
}
```

## MessageAPI

Users can Message another user  
**Method**:POST  
**URL**:http://localhost:8000/message  
**Auth required**: Need token  
If message error

```
{
  "message": "Invalid user details",
  "code":    400,
}
```

If message success

```
{
 "message": "All is fine",
 "code":    200,
}
```
