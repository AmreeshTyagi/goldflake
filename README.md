# Goldflake
Generator to generate Bigint ID or unique key for 174 years (Approx 200k Ids/sec/machine on 8k distributed machines at once) inspired by Twitter's Snowflake and Sonyflake


[![GoDoc](https://godoc.org/github.com/AmreeshTyagi/goldflake?status.svg)](http://godoc.org/github.com/AmreeshTyagi/goldflake)
[![Build Status](https://travis-ci.org/AmreeshTyagi/goldflake.svg?branch=master)](https://travis-ci.org/AmreeshTyagi/goldflake)
[![Coverage Status](https://coveralls.io/repos/AmreeshTyagi/goldflake/badge.svg?branch=master&service=github)](https://coveralls.io/github/AmreeshTyagi/goldflake?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/AmreeshTyagi/goldflake)](https://goreportcard.com/report/github.com/AmreeshTyagi/goldflake)

Goldflake is a distributed unique ID generator inspired by [Twitter's Snowflake](https://blog.twitter.com/2010/announcing-snowflake) & [Sonyflake](https://github.com/sony/sonyflake)

Goal of Goldflake is to provide a unique ID generator which can be used on database level as well as application level to generate bigint unique IDs for lifetime without compromising performance.

So it has a different bit assignment from Snowflake & Sonyflake.
A Goldflake ID is composed of

    39 bits for time in units of 10 msec                      // Similar to Sonyflake
    13 bits for a machine id or app instance id or shard id   // Balance between Snowflake and Sonyflake (If you can live with 8192 machines)
    11 bits for a sequence number                             // To generate 2^11=2048 unique ids per 10 millisecond or 204800 keys/sec on one node
    
## Comparison chart
|                                                  	| Twitter Snowflake       	| Sonyflake                                                   	| Goldflake                                                                                                                                                                                                                                       	|
|--------------------------------------------------	|-------------------------	|-------------------------------------------------------------	|-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------	|
| Lifetime to generate unique id on single machine 	| 69                      	| 174                                                         	| 174                                                                                                                                                                                                                                             	|
| No. of distributed machines                      	| 1024                    	| 65536                                                       	| 8192                                                                                                                                                                                                                                            	|
| No. of unique IDs per 10 milliseconds            	| 40960                   	| 256                                                         	| 2048                                                                                                                                                                                                                                            	|
| Max No. of unique IDs per second on single node  	| 4096000 = 4096k         	| 25600 = 25.6k                                               	| 204800 = 204.8k                                                                                                                                                                                                                                 	|
| Max No. of unique IDs per second on all nodes    	| 4194304000              	| 1677721600                                                  	| 1677721600                                                                                                                                                                                                                                      	|
| Database oriented                                	| Yes                     	| No                                                          	| Yes                                                                                                                                                                                                                                             	|
| Application oriented                             	| Yes                     	| Yes                                                         	| Yes                                                                                                                                                                                                                                             	|
| Comment                                          	| Less number of machines 	| DB will return duplicate soon in  case of heavy concurrency 	| 50% application generated IDs & 50% database generated IDs makes it a best fit to use it on app server and db server.  I believe 4096 database shards are enough to handle potential large load. Though it depends on solution design as well.  	|
