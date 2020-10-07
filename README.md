# Goldflake
Bigint ID or unique key generator for more than 150 years inspired by Twitter's Snowflake and Sonyflake


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
    11 bits for a sequence number
    
