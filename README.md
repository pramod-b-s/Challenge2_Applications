# Challenge_Eluvio

## Problem Statement

Imagine you have a program that needs to look up information about items using their item ID, often in large batches. <br\>

Unfortunately, the only API available for returning this data takes one item at a time, which means you will have to perform one query per item. Additionally, the API is limited to five simultaneous requests. Any additional requests will be served with HTTP 429 (too many requests).<br\>

Write a client utility for your program to use that will retrieve the information for all given IDs as quickly as possible without triggering the simultaneous requests limit, and without performing unnecessary queries for item IDs that have already been seen.<br\>

API Usage:<br\>

GET https://challenges.qluv.io/items/:id<br\>

Required headers:<br\>

Authorization: Base64(:id)<br\>

Example:<br\>

curl https://challenges.qluv.io/items/cRF2dvDZQsmu37WGgK6MTcL7XjH -H "Authorization: Y1JGMmR2RFpRc211MzdXR2dLNk1UY0w3WGpI"<br\>
