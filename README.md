# PingME

Get notifications easily when someone needs you.

Create your personal room, share your room link and receive notifications, when people are trying to reach you (open link). Remember to leave your tab open (Could be changed soon with Service Workers) and allow notifications.

Connection is held within a single tab. Everytime you close it, your room is removed. Max room living time is around 30 days (session expiration time).

## Future plans

* Implement queue, so you receive only a summary of all requests (no need to send dozens of notifications)
* Review backend
* <s>Add referer to url, so you can see where users came from</s>
* Play with Service Workers
* Clean up (when app is closed remove old entries from the database)

**NOTE:** Plans might change depending on free time and interest.

## Q&A

**WHY?** The only reason for this repo is my will to learn GO for writing server-side applications.

**Why do you send ping every 20 seconds?** This keeps connection open, otherwise heroku will close it in ~30 seconds. 


**LICENSE** - MIT
