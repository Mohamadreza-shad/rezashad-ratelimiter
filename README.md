*********************************Architecture*********************************
Project Overview: Rate Limiter Implementation
In this project, I implemented a rate-limiter that controls the number of requests a user can make within a specified time window. Its core responsibility is to ensure fair usage of the API by restricting users from exceeding the allowed rate of requests.

Key Features:
Middleware-Based User Rate Limiting:
Each user's request limit is checked via a middleware that intercepts API requests. The middleware verifies whether the user has exceeded the rate limit within the defined period using Redis as the data storage.

Redis for Data Storage:
All rate-limiting data is stored in Redis. The INCR command is used to track the number of requests per user, and hashes are utilized to store and manage user-specific rate limits. Redis automatically resets the request count after a predefined time window using the EXPIRE command.

Manual Rate-Limit Configuration:
I implemented an endpoint (/api/v1/user/rate-limit) that allows setting custom, manual rate limits for specific users. This feature provides flexibility for managing special rate limits on a per-user basis.

Service Layer Logic:
All the core business logic for rate limiting, request counting, and Redis interactions is encapsulated in the service layer, ensuring clean separation of concerns.

Testing:
End-to-End (E2E) Testing:
I have written 3 tests for this project:
1. "Test_SetUserManualRateLimit_Successful" 
-- The first test (Test_SetUserManualRateLimit_Successful) validates the ability to set a manual rate limit for a specific user by posting configuration data to an endpoint and ensuring that the rate limit is correctly stored in Redis.

2. "TestRateLimiter_1LimitForManualAnd3ForLimit_OneReqShouldGet429"
-- checks if a user, with both manual and default rate limits, gets blocked after exceeding the limit, verifying that one out of five concurrent requests receives an HTTP 429 (Too Many Requests) response.

3. "Test_5ConcurrentReques_ThenWaitForWindowToClose_ThenCallAgain_WeShouldGet200InsteadOf429"
-- confirms that after hitting the rate limit, the user is blocked within the time window but can successfully make a request (HTTP 200) after the rate-limiting window resets.

This setup ensures a robust and scalable rate-limiting solution for controlling API usage across multiple users.


--NOTICE: Please ensure that the Docker Compose containers are running in the background before executing the tests. Use the following command:
***"docker compose -f docker-compose.yml up"***

*********************************Benchmark*********************************
The BenchmarkRateLimit function measures the performance of the RateLimit service. This benchmark simulates repeated requests for a specific userID to test how efficiently the rate-limiting mechanism operates in a Redis-backed environment. The benchmark is particularly useful for evaluating how well the rate limiter scales with increasing traffic and how Redis handles repeated access.

*********************************Scaling*********************************
Scaling in the Rate Limiting Solution:
1. Redis as a Central Data Store: Redis is used for managing rate limits, and it's inherently scalable. Redis can be deployed in a clustered configuration, allowing data to be sharded across multiple nodes. This ensures that even with increasing traffic and data volume, the rate-limiting solution can scale horizontally by adding more Redis nodes.

2. Multiple Application Instances: Since Redis is a centralized data store, the rate limiter can be deployed across multiple instances of the application without any risk of inconsistent rate limit states. Each instance can perform atomic operations (using INCR and EXPIRE) without conflict, ensuring thread safety and consistency across instances.

*********************************Handling Bottlenecks*********************************
1. Redis Bottleneck: If Redis becomes a bottleneck due to high traffic, you can implement Redis clustering or use Redis Sentinel for high availability and automatic failover. This would help distribute the load and provide redundancy.

2. Sharding Requests: Rate limits for different users can be distributed across multiple Redis nodes using sharding based on user IDs. This helps spread the load evenly across Redis instances and prevents overloading a single Redis node.