# Challenge_Applications_Eluvio

## Problem Statement

Imagine you have a program that needs to look up information about items using their item ID, often in large batches. 

Unfortunately, the only API available for returning this data takes one item at a time, which means you will have to perform one query per item. Additionally, the API is limited to five simultaneous requests. Any additional requests will be served with HTTP 429 (too many requests).

Write a client utility for your program to use that will retrieve the information for all given IDs as quickly as possible without triggering the simultaneous requests limit, and without performing unnecessary queries for item IDs that have already been seen.

API Usage:

GET https://challenges.qluv.io/items/:id

Required headers:

Authorization: Base64(:id)

Example:

curl https://challenges.qluv.io/items/cRF2dvDZQsmu37WGgK6MTcL7XjH -H "Authorization: Y1JGMmR2RFpRc211MzdXR2dLNk1UY0w3WGpI"

```
curl https://challenges.qluv.io/items/cRF2dvDZQsmu37WGgK6MTcL7XjH -H "Authorization: Y1JGMmR2RFpRc211MzdXR2dLNk1UY0w3WGpI"
```

## Solution

5 GoRoutines are used to independently schedule the requests from the same queue. A channel with buffer size of 5 is used to implement this queue and send the HTTP GET request for each user ID.


### Backoff Timer
The backoff handler is called by the function `requestInfo()` if a particular itemID's query is not successful.
```
func backoffHandler(item string, iter int, channel chan []interface{}) {}
```
This function runs as another GoRoutine as other items' scheduling should not be blocked on this.


### GoRoutine to request info for URL
Five concurrent instances of this routine will keep calling the `getResponse()` function to send the GET requests and get the response status, if there is a failure, the backoff mechanism is invoked.
```
func requestInfo(base_url string, channel chan []interface{}) {}
```
If the maximum number of attempts is reached for a particular itemID, it is dropped and marked as 'not visited' in a dictionary to be considered again.


### GET request
This function will be invoked by each of the 5 GoRoutines to send a HTTP GET request and get the response.
```
func getResponse(base_url string, _id string) []interface{} {
```
The Authorization header is generated using the item ID by converting it to base64 and the GET request is sent to the URL generated with this header. The response of this GET request is returned to the caller.


### Main
The main function will invoke the 5 GoRoutines explained earlier. A channel of capacity 5 is used to continuously queue itemIDs to be queried.
A dictionary is used to ensure that the itemIDs that have already been queried are not re-queried.
This will be updated by the main function to mark the itemID before sending it to the channel and also by the GoRoutine to mark an itemID if the query fails. A mutex is hence used to concurrently access this dictionary at both places. A Waitgroup ensures that the queries are atomic before the channel is closed. A number of random IDs are generated to verify the functionality of the client API.


### Test Server
A server is simulated to check the behavior indicated in the problem statement. If more than 5 requests are sent concurrently, Error 429 is thrown and some random 404 errors are generated so the server behaves realistically. The implementation has been verified.
