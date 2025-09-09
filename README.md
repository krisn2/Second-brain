# Second-brain
a note taking app that hold the imp or read later articles or notes

## End pointes 

### Sign up

Post /api/v1/signup

```
{
    "username":"krisn",
    "password":"123123"
}

```

**Constrainst** -

1. username should be 3-10 lettes.

2. password should e 8 to 20 letter, should have atlest one uppercase, one lowercase. on special character, one number .

**Responses**
1. Status 200 - Sign up
2. Status 422 - Error in inputs
3. Status 403 - User already exists with this username
4. Status 500 - Server error


### Sign in

POST  /api/v1/signin

```
{
    "username":"krisn",
    "password:"1232112"
}

```

> Returns
> **200**

```
{
    "token":"jwt_token"
}
```

> **403 - Wrong email password**
> **500 - Internal server error**


### Add new content
> **Note:** only auth person add 

POST /api/v1/content

```
{
    "type":"document" | "tweet" | "youtube" | "link",
    "link":"url",
    "title":"Title of dc/video",
    "tags":["productivity", "politics",...]
}
```

> Return 
> **200**


### Fetching all existing documents (no pagination)

GET /api/v1/content

Returns 
```
{
	"content": [
		{
			"id": 1,
			"type": "document" | "tweet" | "youtube" | "link",
			"link": "url",
			"title": "Title of doc/video",
			"tags": ["productivity", "politics", ...]
		}
	
	]
}
```

### search Brain 

GET /api/v1/content/:query

return 
```
{
	"ai_response":"string",
	"original_content:" "content": [
		{
			"id": 1,
			"type": "document" | "tweet" | "youtube" | "link",
			"link": "url",
			"title": "Title of doc/video",
			"tags": ["productivity", "politics", ...]
		}
	
	]"
}

### Delete a document

DELETE /api/v1/content

```
{
	"contentId": "1"
}
```

>Returns

1. 200 - Delete suceeded
2. 403 - Trying to delete a doc you don't own