*********************************Architecture*******************************************
Project Overview: Rate Limiter Implementation
In this project, I implemented a rate-limiter that controls the number of requests a user can make within a specified time window. Its core responsibility is to ensure fair usage of the API by restricting users from exceeding the allowed rate of requests.

Key Features:
Middleware-Based User Rate Limiting:
Each user's request limit is checked via a middleware that intercepts API requests. The middleware verifies whether the user has exceeded the rate limit within the defined period using Redis as the backend.

Redis for Data Storage:
All rate-limiting data is stored in Redis. Redis's sorted sets and hashes are utilized to track user requests and manage dynamic rate limits.

Manual Rate-Limit Configuration:
I implemented an endpoint (/api/v1/user/rate-limit) that allows setting custom, manual rate limits for specific users. This feature provides flexibility for managing special rate limits on a per-user basis.

Service Layer Logic:
All the core business logic for rate limiting, request counting, and Redis interactions is encapsulated in the service layer, ensuring clean separation of concerns.

Testing:
End-to-End (E2E) Testing:
I wrote comprehensive E2E tests for the /user/rate-limit endpoint to ensure correct functionality under various conditions.

Unit Tests and Benchmarking:
Additional unit tests and benchmarks were written for the service layer to ensure the rate-limiting logic performs efficiently and handles edge cases effectively.

This setup ensures a robust and scalable rate-limiting solution for controlling API usage across multiple users.

*********************************RateLimit Function*****************************************
#About RateLimit Function

The RateLimit method enforces a rate-limiting policy for a specific user, identified by a unique userID. It tracks the number of requests made by the user within a sliding window of time, using Redis sorted sets. If the user exceeds their allowed request limit within that window, the method denies further requests until the window resets.

The method also supports fetching and applying additional user-specific limits stored in Redis.

Parameters
userID (string):
A unique identifier for the user whose rate is being limited.

limit (int):
The default rate limit (maximum number of requests) that the user is allowed within the sliding window. This is a baseline limit and may be adjusted dynamically based on user-specific settings.

Return Value
bool:
Returns true if the request is allowed (i.e., the user has not exceeded the rate limit).
Returns false if the request is denied (i.e., the user has exceeded the rate limit).

Redis Data Structures
Sorted Set (ZSET):
Used to store the timestamps of requests for each user, where the score is the timestamp (in milliseconds). The set stores requests within the sliding time window.

Key Format: "user:<userID>:limit"
Hash (HSET):
Stores the user-specific rate limit configuration in Redis.

Key Format: "userID:<userID>"
Field: REDIS_RATE_LIMIT_FIELD
Contains the user-specific rate limit.
Logic Breakdown
Key Definition:

The Redis key used to track the requests for the user is defined as "user:<userID>:limit".
Removing Old Requests:

Any requests that fall outside the current sliding window are removed from the sorted set using ZRemRangeByScore, where the scores represent timestamps outside of the current window.

Counting Requests:
The method counts the number of requests within the sliding window using ZCount, which checks the sorted set between the start of the window (windowStart) and the current time.

Fetching User-Specific Rate Limit:
The method attempts to fetch a user-specific rate limit from Redis using the HGet command. If the user's rate limit is not found, a default rate limit is set using HSet.

Rate Limiting Check:
The method compares the number of requests (count) within the sliding window to the combined default and user-specific rate limits (limit + userLimit). If the user has exceeded the limit, the method returns false, indicating the request is denied.

Allowing the Request:
If the user has not exceeded the limit, the method adds the current timestamp as a new request to the sorted set and ensures the sorted set has an expiration time (TTL) set.
